package jsonengine

import (
	"context"
	"fmt"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
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
	Result     PolicyResult
	Message    string
}

// PolicyResult specifies state of a policy result
type PolicyResult string

const (
	StatusPass  PolicyResult = "pass"
	StatusFail  PolicyResult = "fail"
	StatusWarn  PolicyResult = "warn"
	StatusError PolicyResult = "error"
	StatusSkip  PolicyResult = "skip"
)

func New() engine.Engine[Request, Response] {
	type ruleRequest struct {
		rule     v1alpha1.ValidatingRule
		value    any
		bindings jpbinding.Bindings
	}
	type policyRequest struct {
		policy   *v1alpha1.ValidatingPolicy
		value    any
		bindings jpbinding.Bindings
	}
	ruleEngine := builder.
		Function(func(ctx context.Context, r ruleRequest) []RuleResponse {
			if r.rule.Match == nil {
				errs, err := matching.Match(ctx, nil, r.rule.Match, r.value, r.bindings)
				if err != nil {
					// TODO return error
				}
				// didn't match
				if len(errs) != 0 {
					return nil
				}
			}
			if r.rule.Exclude != nil {
				errs, err := matching.Match(ctx, nil, r.rule.Exclude, r.value, r.bindings)
				if err != nil {
					// TODO return error
				}
				// matched
				if len(errs) == 0 {
					return nil
				}
			}
			errs, err := matching.MatchAssert(ctx, nil, r.rule.Assert, r.value, r.bindings)
			if err != nil {
				// TODO return error
			}
			if len(errs) == 0 {
				return []RuleResponse{{
					Rule:    r.rule,
					Result:  StatusPass,
					Message: "",
				}}
			}
			identifier := ""
			if r.rule.Identifier != "" {
				result, subjectErr := template.Execute(context.Background(), r.rule.Identifier, r.value, nil)
				if subjectErr != nil {
					identifier = fmt.Sprintf("(error: %s)", subjectErr)
				} else {
					identifier = fmt.Sprint(result)
				}
			}
			var failures []RuleResponse
			for _, err := range errs {
				failures = append(failures, RuleResponse{
					Rule:       r.rule,
					Identifier: identifier,
					Result:     StatusFail,
					Message:    err.Error(),
				})
			}
			return failures
		})
	policyEngine := builder.
		Function(func(ctx context.Context, r policyRequest) PolicyResponse {
			response := PolicyResponse{
				Policy: r.policy,
			}
			for _, rule := range r.policy.Spec.Rules {
				response.Rules = append(response.Rules, ruleEngine.Run(ctx, ruleRequest{
					rule:  rule,
					value: r.value,
				})...)
			}
			return response
		})
	e := builder.
		Function(func(ctx context.Context, r Request) Response {
			response := Response{
				Resource: r.Resource,
			}
			for _, policy := range r.Policies {
				response.Policies = append(response.Policies, policyEngine.Run(ctx, policyRequest{
					policy: policy,
					value:  r.Resource,
				}))
			}
			return response
		})
	return e
	// looper := func(r Request) []request {
	// 	var requests []request
	// 	bindings := jpbinding.NewBindings()
	// 	bindings = bindings.Register("$payload", jpbinding.NewBinding(r.Resource))
	// 	for _, policy := range r.Policies {
	// 		bindings = bindings.Register("$policy", jpbinding.NewBinding(policy))
	// 		for _, rule := range policy.Spec.Rules {
	// 			bindings = bindings.Register("$rule", jpbinding.NewBinding(rule))
	// 			bindings = binding.NewContextBindings(bindings, r.Resource, rule.Context...)
	// 			requests = append(requests, request{
	// 				policy:   policy,
	// 				rule:     rule,
	// 				value:    r.Resource,
	// 				bindings: bindings,
	// 			})
	// 		}
	// 	}
	// 	return requests
	// }
	// inner := builder.
	// 	Function(func(ctx context.Context, r request) RuleResponse {
	// 		errs, err := matching.MatchAssert(ctx, nil, r.rule.Assert, r.value, r.bindings)
	// 		response := buildResponse(r, errs, err)
	// 		return response
	// 	}).
	// 	Predicate(func(ctx context.Context, r request) bool {
	// 		if r.rule.Exclude == nil {
	// 			return true
	// 		}
	// 		errs, err := matching.Match(ctx, nil, r.rule.Exclude, r.value, r.bindings)
	// 		// TODO: handle error and skip
	// 		return err == nil && len(errs) != 0
	// 	}).
	// 	Predicate(func(ctx context.Context, r request) bool {
	// 		if r.rule.Match == nil {
	// 			return true
	// 		}
	// 		errs, err := matching.Match(ctx, nil, r.rule.Match, r.value, r.bindings)
	// 		// TODO: handle error and skip
	// 		return err == nil && len(errs) == 0
	// 	})
	// // TODO: we can't use the builder package for loops :(
	// return loop.New(inner, looper)
}

// func buildResponse(req request, fails []error, ruleErr error) RuleResponse {
// 	response := RuleResponse{
// 		PolicyName: req.policy.Name,
// 		RuleName:   req.rule.Name,
// 	}

// 	response.Identifier = ""
// 	if req.rule.Identifier != "" {
// 		result, subjectErr := template.Execute(context.Background(), req.rule.Identifier, req.value, nil)
// 		if subjectErr != nil {
// 			response.Identifier = fmt.Sprintf("(error: %s)", subjectErr)
// 		} else {
// 			response.Identifier = fmt.Sprint(result)
// 		}
// 	}

// 	if ruleErr != nil {
// 		response.Result = StatusError
// 		response.Message = ruleErr.Error()
// 	} else if err := multierr.Combine(fails...); err != nil {
// 		response.Result = StatusFail
// 		response.Message = err.Error()
// 	} else {
// 		// TODO: handle skip result
// 		response.Result = StatusPass
// 	}

// 	return response
// }
