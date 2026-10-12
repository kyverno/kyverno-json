package jsonengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestServiceCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"foo": "bar"}`))
			return
		}
		if r.URL.Path == "/post" && r.Method == "POST" {
			var m map[string]any
			json.NewDecoder(r.Body).Decode(&m)
			if m["key"] == "value" && m["key2"] == "value2" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success": true}`))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		apiCall v1alpha1.APICall
		want    string
		wantErr bool
	}{
		{
			name: "basic get",
			apiCall: v1alpha1.APICall{
				Service: &v1alpha1.ServiceCall{
					URL: v1alpha1.NewAny(server.URL + "/test"),
				},
			},
			want: `{"foo":"bar"}`,
		},
		{
			name: "post with data",
			apiCall: v1alpha1.APICall{
				Method: "POST",
				Service: &v1alpha1.ServiceCall{
					URL: v1alpha1.NewAny(server.URL + "/post"),
				},
				Data: []v1alpha1.RequestData{
					{
						Key:   "key",
						Value: v1alpha1.NewAny("value"),
					},
					{
						Key:   "key2",
						Value: v1alpha1.NewAny("value2"),
					},
				},
			},
			want: `{"success":true}`,
		},
		{
			name: "jmespath filter",
			apiCall: v1alpha1.APICall{
				Service: &v1alpha1.ServiceCall{
					URL: v1alpha1.NewAny(server.URL + "/test"),
				},
				JMESPath: "foo",
			},
			want: `"bar"`,
		},
		{
			name: "default fallback",
			apiCall: v1alpha1.APICall{
				Service: &v1alpha1.ServiceCall{
					URL: v1alpha1.NewAny(server.URL + "/404"),
				},
				Default: func() *v1alpha1.Any { a := v1alpha1.NewAny("fallback"); return &a }(),
			},
			want: `"fallback"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := v1alpha1.ContextEntry{
				Name:    "testCall",
				APICall: &tt.apiCall,
			}
			path := field.NewPath("context")
			fn, err := CompileContextEntry(path, compilers.DefaultCompilers, nil, in)
			if err != nil {
				t.Fatalf("CompileContextEntry failed: %v", err)
			}
			b := binding.NewBindings()
			b = fn(context.Background(), nil, b)
			bind, _ := b.Get("$testCall")
			v, errVal := bind.Value()
			if tt.wantErr {
				if errVal == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if errVal != nil {
				t.Fatalf("unexpected error: %v", errVal)
			}
			gotBytes, _ := json.Marshal(v)
			if string(gotBytes) != tt.want {
				t.Errorf("got %s, want %s", string(gotBytes), tt.want)
			}
		})
	}
}

type fakeK8sClient struct {
	data           []byte
	err            error
	expectedMethod string
	expectedPath   string
	expectedData   string
	t              *testing.T
}

func (c *fakeK8sClient) DoRaw(ctx context.Context, method, path string, data []byte) ([]byte, error) {
	if method != c.expectedMethod {
		c.t.Errorf("expected method %s, got %s", c.expectedMethod, method)
	}
	if path != c.expectedPath {
		c.t.Errorf("expected path %s, got %s", c.expectedPath, path)
	}
	if string(data) != c.expectedData {
		c.t.Errorf("expected data %s, got %s", c.expectedData, string(data))
	}
	return c.data, c.err
}

func TestURLPath(t *testing.T) {
	client := &fakeK8sClient{
		data:           []byte(`{"kind": "Pod"}`),
		expectedMethod: "GET", // Default method
		expectedPath:   "/api/v1/pods",
		expectedData:   "",
		t:              t,
	}

	in := v1alpha1.ContextEntry{
		Name: "testK8s",
		APICall: &v1alpha1.APICall{
			URLPath: v1alpha1.NewAny("/api/v1/pods"),
		},
	}
	path := field.NewPath("context")
	fn, err := CompileContextEntry(path, compilers.DefaultCompilers, client, in)
	if err != nil {
		t.Fatalf("CompileContextEntry failed: %v", err)
	}

	b := binding.NewBindings()
	b = fn(context.Background(), nil, b)
	bind, _ := b.Get("$testK8s")
	v, errVal := bind.Value()
	if errVal != nil {
		t.Fatalf("unexpected error: %v", errVal)
	}

	gotBytes, _ := json.Marshal(v)
	if string(gotBytes) != `{"kind":"Pod"}` {
		t.Errorf("got %s, want %s", string(gotBytes), `{"kind":"Pod"}`)
	}
}
