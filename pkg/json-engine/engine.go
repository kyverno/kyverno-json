package jsonengine

import (
	"context"
	"fmt"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/loop"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/matching"
	"go.uber.org/multierr"
)

type Request struct {
	Resources []interface{}
	Policies  []*v1alpha1.ValidatingPolicy
}

type Response struct {
	Results []RuleResponse `json:"results"`
}

type RuleResponse struct {
	PolicyName string       `json:"policy"`
	RuleName   string       `json:"rule"`
	Identifier string       `json:"identifier,omitempty"`
	Result     PolicyResult `json:"result"`
	Message    string       `json:"message"`
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

type request struct {
	policy   *v1alpha1.ValidatingPolicy
	rule     v1alpha1.ValidatingRule
	value    interface{}
	bindings jpbinding.Bindings
}

func New() engine.Engine[Request, RuleResponse] {
	looper := func(r Request) []request {
		var requests []request
		bindings := jpbinding.NewBindings()
		for _, resource := range r.Resources {
			bindings = bindings.Register("$payload", jpbinding.NewBinding(resource))
			for _, policy := range r.Policies {
				bindings = bindings.Register("$policy", jpbinding.NewBinding(policy))
				for _, rule := range policy.Spec.Rules {
					bindings = bindings.Register("$rule", jpbinding.NewBinding(rule))
					bindings = binding.NewContextBindings(bindings, resource, rule.Context...)
					requests = append(requests, request{
						policy:   policy,
						rule:     rule,
						value:    resource,
						bindings: bindings,
					})
				}
			}
		}
		return requests
	}
	inner := builder.
		Function(func(ctx context.Context, r request) RuleResponse {
			errs, err := matching.MatchAssert(ctx, nil, r.rule.Assert, r.value, r.bindings)
			response := buildResponse(r, errs, err)
			return response
		}).
		Predicate(func(ctx context.Context, r request) bool {
			if r.rule.Exclude == nil {
				return true
			}
			errs, err := matching.Match(ctx, nil, r.rule.Exclude, r.value, r.bindings)
			// TODO: handle error and skip
			return err == nil && len(errs) != 0
		}).
		Predicate(func(ctx context.Context, r request) bool {
			if r.rule.Match == nil {
				return true
			}
			errs, err := matching.Match(ctx, nil, r.rule.Match, r.value, r.bindings)
			// TODO: handle error and skip
			return err == nil && len(errs) == 0
		})
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}

func buildResponse(req request, fails []error, ruleErr error) RuleResponse {
	response := RuleResponse{
		PolicyName: req.policy.Name,
		RuleName:   req.rule.Name,
	}

	response.Identifier = ""
	if req.rule.Identifier != "" {
		result, subjectErr := template.Execute(context.Background(), req.rule.Identifier, req.value, nil)
		if subjectErr != nil {
			response.Identifier = fmt.Sprintf("(error: %s)", subjectErr)
		} else {
			response.Identifier = fmt.Sprint(result)
		}
	}

	if ruleErr != nil {
		response.Result = StatusError
		response.Message = ruleErr.Error()
	} else if err := multierr.Combine(fails...); err != nil {
		response.Result = StatusFail
		response.Message = err.Error()
	} else {
		// TODO: handle skip result
		response.Result = StatusPass
	}

	return response
}
