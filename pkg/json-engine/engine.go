package jsonengine

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
)

func New() engine.Engine[Request, Response] {
	type ruleRequest struct {
		policy   v1alpha1.ValidatingPolicy
		rule     v1alpha1.ValidatingRule
		resource any
		bindings binding.Bindings
	}
	type policyRequest struct {
		policy   v1alpha1.ValidatingPolicy
		resource any
		bindings binding.Bindings
	}
	compilers := compilers.DefaultCompilers.WithDefaultCompiler(compilers.CompilerJP)
	ruleCompiler := compiler{}
	ruleEngine := builder.
		Function(func(ctx context.Context, r ruleRequest) []RuleResponse {
			compilers := compilers
			if r.policy.Spec.Compiler != nil {
				compilers = compilers.WithDefaultCompiler(string(*r.policy.Spec.Compiler))
			}
			compiled, err := ruleCompiler.compileRule(nil, compilers, r.rule)
			if err != nil {
				return []RuleResponse{{
					Rule:      r.rule,
					Timestamp: time.Now(),
					Error:     err,
				}}
			}
			bindings := r.bindings.Register("$rule", binding.NewBinding(r.rule))
			return compiled(r.resource, bindings)
		})
	policyEngine := builder.
		Function(func(ctx context.Context, r policyRequest) PolicyResponse {
			response := PolicyResponse{
				Policy: r.policy,
			}
			bindings := r.bindings.Register("$policy", binding.NewBinding(r.policy))
			for _, rule := range r.policy.Spec.Rules {
				response.Rules = append(response.Rules, ruleEngine.Run(ctx, ruleRequest{
					policy:   r.policy,
					rule:     rule,
					resource: r.resource,
					bindings: bindings.Register("$rule", binding.NewBinding(rule)),
				})...)
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
