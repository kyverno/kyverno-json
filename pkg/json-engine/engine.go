package jsonengine

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func New() engine.Engine[Request, Response] {
	type ruleRequest struct {
		compilers compilers.Compilers
		path      *field.Path
		policy    v1alpha1.ValidatingPolicy
		rule      v1alpha1.ValidatingRule
		resource  any
		bindings  binding.Bindings
	}
	type policyRequest struct {
		policy   v1alpha1.ValidatingPolicy
		resource any
		bindings binding.Bindings
	}
	ruleEngine := builder.
		Function(func(ctx context.Context, r ruleRequest) *RuleResponse {
			ruleCompiler := compiler{}
			compiled, err := ruleCompiler.compileRule(r.path, r.compilers, r.rule)
			if err != nil {
				return &RuleResponse{
					Rule:      r.rule,
					Timestamp: time.Now(),
					Error:     err,
				}
			}
			bindings := r.bindings.Register("$rule", binding.NewBinding(r.rule))
			return compiled(r.resource, bindings)
		})
	policyEngine := builder.
		Function(func(ctx context.Context, r policyRequest) PolicyResponse {
			response := PolicyResponse{
				Policy: r.policy,
			}
			path := field.NewPath("spec", "rules")
			compilers := compilers.DefaultCompilers.WithDefaultCompiler(compilers.CompilerJP)
			if r.policy.Spec.Compiler != nil {
				compilers = compilers.WithDefaultCompiler(string(*r.policy.Spec.Compiler))
			}
			bindings := r.bindings.Register("$policy", binding.NewBinding(r.policy))
			for i, rule := range r.policy.Spec.Rules {
				if ruleResponse := ruleEngine.Run(ctx, ruleRequest{
					compilers: compilers,
					path:      path.Index(i),
					policy:    r.policy,
					rule:      rule,
					resource:  r.resource,
					bindings:  bindings.Register("$rule", binding.NewBinding(rule)),
				}); ruleResponse != nil {
					response.Rules = append(response.Rules, *ruleResponse)
				}
			}
			return response
		})
	resourceEngine := builder.
		Function(func(ctx context.Context, r Request) Response {
			response := Response{
				Resource: r.Resource,
			}
			bindings := binding.NewBindings()
			for k, v := range r.Bindings {
				bindings = bindings.Register("$"+k, binding.NewBinding(v))
			}
			bindings = bindings.Register("$payload", binding.NewBinding(r.Resource))
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
