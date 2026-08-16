package jsonengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/projection"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var defaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func compileAPICall(
	path *field.Path,
	comps compilers.Compilers,
	k8sClient K8sClient,
	in v1alpha1.ContextEntry,
) (func(context.Context, any, binding.Bindings) binding.Bindings, *field.Error) {
	apiCall := in.APICall
	if apiCall == nil {
		return nil, field.Invalid(path, in, "apiCall cannot be nil")
	}

	var urlHandler projection.ScalarHandler
	var err *field.Error
	isK8s := false

	if apiCall.Service != nil {
		urlHandler, err = apiCall.Service.URL.Compile(path.Child("service", "url"), comps)
	} else if !apiCall.URLPath.IsNil() {
		isK8s = true
		urlHandler, err = apiCall.URLPath.Compile(path.Child("urlPath"), comps)
	} else {
		return nil, field.Invalid(path, in, "either service or urlPath must be provided")
	}

	if err != nil {
		return nil, err
	}

	method := apiCall.Method
	if method == "" {
		method = http.MethodGet
	}

	var dataHandlers []func(any, binding.Bindings) (string, any, error)
	for i, d := range apiCall.Data {
		h, err := d.Value.Compile(path.Child("data").Index(i).Child("value"), comps)
		if err != nil {
			return nil, err
		}
		key := d.Key
		dataHandlers = append(dataHandlers, func(res any, b binding.Bindings) (string, any, error) {
			val, err := h(res, b)
			return key, val, err
		})
	}

	var jpHandler projection.ScalarHandler
	if apiCall.JMESPath != "" {
		c := comps.Compiler(compilers.CompilerJP)
		if c == nil {
			return nil, field.Invalid(path.Child("jmesPath"), apiCall.JMESPath, "jmespath compiler not found")
		}
		program, err2 := c.Compile(apiCall.JMESPath)
		if err2 != nil {
			return nil, field.Invalid(path.Child("jmesPath"), apiCall.JMESPath, err2.Error())
		}
		jpHandler = func(value any, bindings binding.Bindings) (any, error) {
			return program(value, bindings)
		}
	}

	var defaultHandler projection.ScalarHandler
	if apiCall.Default != nil {
		defaultHandler, err = apiCall.Default.Compile(path.Child("default"), comps)
		if err != nil {
			return nil, err
		}
	}

	return func(ctx context.Context, resource any, bindings binding.Bindings) binding.Bindings {
		b := binding.NewDelegate(
			sync.OnceValues(
				func() (any, error) {
					urlVal, err2 := urlHandler(resource, bindings)
					if err2 != nil {
						if defaultHandler != nil {
							return defaultHandler(resource, bindings)
						}
						return nil, err2
					}
					urlStr, ok := urlVal.(string)
					if !ok {
						err2 = fmt.Errorf("url must evaluate to string, got %T", urlVal)
						if defaultHandler != nil {
							return defaultHandler(resource, bindings)
						}
						return nil, err2
					}

					var reqData []byte
					if len(dataHandlers) > 0 {
						dataMap := make(map[string]any)
						for _, dh := range dataHandlers {
							k, v, err2 := dh(resource, bindings)
							if err2 != nil {
								if defaultHandler != nil {
									return defaultHandler(resource, bindings)
								}
								return nil, err2
							}
							dataMap[k] = v
						}
						reqData, _ = json.Marshal(dataMap)
					}

					var respData []byte

					if isK8s {
						if k8sClient == nil {
							err2 = fmt.Errorf("Kubernetes API calls are not supported in offline CLI mode")
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, err2
						}
						respData, err2 = k8sClient.DoRaw(ctx, method, urlStr, reqData)
					} else {
						client := defaultHTTPClient

						var bodyReader io.Reader
						if len(reqData) > 0 {
							bodyReader = bytes.NewReader(reqData)
						}
						req, err3 := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
						if err3 != nil {
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, err3
						}
						if bodyReader != nil {
							req.Header.Set("Content-Type", "application/json")
						}

						resp, err3 := client.Do(req)
						if err3 != nil {
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, err3
						}
						defer resp.Body.Close()

						if resp.StatusCode < 200 || resp.StatusCode >= 300 {
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
						}

						resp.Body = http.MaxBytesReader(nil, resp.Body, 2*1024*1024)
						respData, err3 = io.ReadAll(resp.Body)
						if err3 != nil {
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, err3
						}
					}

					if err2 != nil {
						if defaultHandler != nil {
							return defaultHandler(resource, bindings)
						}
						return nil, err2
					}

					var parsed any
					if err3 := json.Unmarshal(respData, &parsed); err3 != nil {
						if defaultHandler != nil {
							return defaultHandler(resource, bindings)
						}
						return nil, err3
					}

					if jpHandler != nil {
						parsed, err2 = jpHandler(parsed, bindings)
						if err2 != nil {
							if defaultHandler != nil {
								return defaultHandler(resource, bindings)
							}
							return nil, err2
						}
					}

					return parsed, nil
				},
			),
		)
		return bindings.Register("$"+in.Name, b)
	}, nil
}
