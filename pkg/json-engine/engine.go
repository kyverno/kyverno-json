package jsonengine

import (
	"context"
	"fmt"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/matching"
)

type Request struct {
	Resource any
	Policies []*v1alpha1.ValidatingPolicy
}

type Response struct {
	Resource any
	Policies []PolicyResponse
}

type PolicyResponse struct {
	Policy *v1alpha1.ValidatingPolicy
	Rules  []RuleResponse
}

type RuleResponse struct {
	Rule       v1alpha1.ValidatingRule
	Identifier string
	Error      error
	Violations []error
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
		rule     v1alpha1.ValidatingRule
		resource any
		bindings jpbinding.Bindings
	}
	type policyRequest struct {
		policy   *v1alpha1.ValidatingPolicy
		resource any
		bindings jpbinding.Bindings
	}
	ruleEngine := builder.
		Function(func(ctx context.Context, r ruleRequest) []RuleResponse {
			bindings := r.bindings.Register("$rule", jpbinding.NewBinding(r.rule))
			bindings = binding.NewContextBindings(bindings, r.resource, r.rule.Context...)
			identifier := ""
			if r.rule.Identifier != "" {
				result, err := template.Execute(context.Background(), r.rule.Identifier, r.resource, bindings)
				if err != nil {
					identifier = fmt.Sprintf("(error: %s)", err)
				} else {
					identifier = fmt.Sprint(result)
				}
			}
			if r.rule.Match != nil {
				errs, err := matching.Match(ctx, nil, r.rule.Match, r.resource, bindings)
				if err != nil {
					return []RuleResponse{{
						Rule:       r.rule,
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
				errs, err := matching.Match(ctx, nil, r.rule.Exclude, r.resource, bindings)
				if err != nil {
					return []RuleResponse{{
						Rule:       r.rule,
						Identifier: identifier,
						Error:      err,
					}}
				}
				// matched
				if len(errs) == 0 {
					return nil
				}
			}
			violations, err := matching.MatchAssert(ctx, nil, r.rule.Assert, r.resource, bindings)
			if err != nil {
				return []RuleResponse{{
					Rule:       r.rule,
					Identifier: identifier,
					Error:      err,
				}}
			}
			return []RuleResponse{{
				Rule:       r.rule,
				Identifier: identifier,
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
			bindings := jpbinding.NewBindings().Register("$payload", jpbinding.NewBinding(r.Resource))
			for _, policy := range r.Policies {
				response.Policies = append(response.Policies, policyEngine.Run(ctx, policyRequest{
					policy:   policy,
					resource: r.Resource,
					bindings: bindings,
				}))
			}
			return response
		})
	return resourceEngine
}
