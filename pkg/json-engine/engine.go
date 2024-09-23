package jsonengine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"github.com/kyverno/kyverno-json/pkg/matching"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Request struct {
	Resource any
	Policies []*v1alpha1.ValidatingPolicy
	Bindings map[string]any
}

type Response struct {
	Resource any
	Policies []PolicyResponse
}

type PolicyResponse struct {
	Policy v1alpha1.ValidatingPolicy
	Rules  []RuleResponse
}

type RuleResponse struct {
	Rule       v1alpha1.ValidatingRule
	Timestamp  time.Time
	Identifier string
	Feedback   map[string]Feedback
	Error      error
	Violations matching.Results
}

type Feedback struct {
	Error error
	Value any
}

// PolicyResult specifies state of a policy result
type PolicyResult string

const (
	StatusPass PolicyResult = "pass"
	StatusFail PolicyResult = "fail"
	// StatusWarn  PolicyResult = "warn"
	StatusError PolicyResult = "error"
	// StatusSkip  PolicyResult = "skip"
)

func New() engine.Engine[Request, Response] {
	type ruleRequest struct {
		policy   v1alpha1.ValidatingPolicy
		rule     v1alpha1.ValidatingRule
		resource any
		bindings jpbinding.Bindings
	}
	type policyRequest struct {
		policy   v1alpha1.ValidatingPolicy
		resource any
		bindings jpbinding.Bindings
	}
	compiler := matching.NewCompiler(compilers.DefaultCompilers, 256)
	ruleEngine := builder.
		Function(func(ctx context.Context, r ruleRequest) []RuleResponse {
			bindings := r.bindings.Register("$rule", jpbinding.NewBinding(r.rule))
			defaultCompiler := expression.CompilerJP
			if r.policy.Spec.Compiler != nil {
				defaultCompiler = string(*r.policy.Spec.Compiler)
			}
			if r.rule.Compiler != nil {
				defaultCompiler = string(*r.rule.Compiler)
			}
			// TODO: this doesn't seem to be the right path
			var path *field.Path
			path = path.Child("context")
			for _, entry := range r.rule.Context {
				defaultCompiler := defaultCompiler
				if entry.Compiler != nil {
					defaultCompiler = string(*entry.Compiler)
				}
				bindings = func(variable v1alpha1.Any, bindings jpbinding.Bindings) jpbinding.Bindings {
					return bindings.Register(
						"$"+entry.Name,
						binding.NewDelegate(
							sync.OnceValues(
								func() (any, error) {
									handler, err := variable.Compile(compiler.CompileProjection, defaultCompiler)
									if err != nil {
										return nil, field.InternalError(path.Child("variable"), err)
									}
									projected, err := handler(r.resource, bindings)
									if err != nil {
										return nil, field.InternalError(path.Child("variable"), err)
									}
									return projected, nil
								},
							),
						),
					)
				}(entry.Variable, bindings)
			}
			identifier := ""
			if r.rule.Identifier != "" {
				result, err := compilers.Execute(r.rule.Identifier, r.resource, bindings, compiler.Jp)
				if err != nil {
					identifier = fmt.Sprintf("(error: %s)", err)
				} else {
					identifier = fmt.Sprint(result)
				}
			}
			if r.rule.Match != nil {
				defaultCompiler := defaultCompiler
				if r.rule.Match.Compiler != nil {
					defaultCompiler = string(*r.rule.Match.Compiler)
				}
				errs, err := matching.Match(nil, r.rule.Match, r.resource, bindings, compiler, defaultCompiler)
				if err != nil {
					return []RuleResponse{{
						Rule:       r.rule,
						Timestamp:  time.Now(),
						Identifier: identifier,
						Error:      err,
					}}
				}
				// didn't match
				if len(errs) != 0 {
					return nil
				}
			}
			if r.rule.Exclude != nil {
				defaultCompiler := defaultCompiler
				if r.rule.Exclude.Compiler != nil {
					defaultCompiler = string(*r.rule.Exclude.Compiler)
				}
				errs, err := matching.Match(nil, r.rule.Exclude, r.resource, bindings, compiler, defaultCompiler)
				if err != nil {
					return []RuleResponse{{
						Rule:       r.rule,
						Timestamp:  time.Now(),
						Identifier: identifier,
						Error:      err,
					}}
				}
				// matched
				if len(errs) == 0 {
					return nil
				}
			}
			var feedback map[string]Feedback
			for _, f := range r.rule.Feedback {
				entry := Feedback{}
				if f.Value != nil {
					defaultCompiler := defaultCompiler
					if f.Compiler != nil {
						defaultCompiler = string(*f.Compiler)
					}
					if handler, err := f.Value.Compile(compiler.CompileProjection, defaultCompiler); err != nil {
						entry.Error = err
					} else if projected, err := handler(r.resource, bindings); err != nil {
						entry.Error = err
					} else {
						entry.Value = projected
					}
				}
				if feedback == nil {
					feedback = map[string]Feedback{}
				}
				feedback[f.Name] = entry
			}
			violations, err := matching.Assert(nil, r.rule.Assert, r.resource, bindings, compiler, defaultCompiler)
			if err != nil {
				return []RuleResponse{{
					Rule:       r.rule,
					Timestamp:  time.Now(),
					Identifier: identifier,
					Feedback:   feedback,
					Error:      err,
				}}
			}
			return []RuleResponse{{
				Rule:       r.rule,
				Timestamp:  time.Now(),
				Identifier: identifier,
				Feedback:   feedback,
				Violations: violations,
			}}
		})
	policyEngine := builder.
		Function(func(ctx context.Context, r policyRequest) PolicyResponse {
			response := PolicyResponse{
				Policy: r.policy,
			}
			bindings := r.bindings.Register("$policy", jpbinding.NewBinding(r.policy))
			for _, rule := range r.policy.Spec.Rules {
				response.Rules = append(response.Rules, ruleEngine.Run(ctx, ruleRequest{
					rule:     rule,
					resource: r.resource,
					bindings: bindings.Register("$rule", jpbinding.NewBinding(rule)),
				})...)
			}
			return response
		})
	resourceEngine := builder.
		Function(func(ctx context.Context, r Request) Response {
			response := Response{
				Resource: r.Resource,
			}
			bindings := jpbinding.NewBindings()
			for k, v := range r.Bindings {
				bindings = bindings.Register("$"+k, jpbinding.NewBinding(v))
			}
			bindings = bindings.Register("$payload", jpbinding.NewBinding(r.Resource))
			for _, policy := range r.Policies {
				response.Policies = append(response.Policies, policyEngine.Run(ctx, policyRequest{
					policy:   *policy,
					resource: r.Resource,
					bindings: bindings,
				}))
			}
			return response
		})
	return resourceEngine
}
