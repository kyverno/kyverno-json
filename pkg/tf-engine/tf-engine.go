package tfengine

import (
	"errors"

	"github.com/eddycharly/tf-kyverno/pkg/apis/v1alpha1"
	"github.com/eddycharly/tf-kyverno/pkg/engine"
	"github.com/eddycharly/tf-kyverno/pkg/engine/blocks/loop"
	"github.com/eddycharly/tf-kyverno/pkg/engine/builder"
	"github.com/eddycharly/tf-kyverno/pkg/match"
	"github.com/eddycharly/tf-kyverno/pkg/plan"
)

type TfEngineRequest struct {
	Plan     *plan.Plan
	Policies []*v1alpha1.Policy
}

type TfEngineResponse struct {
	Policy   *v1alpha1.Policy
	Rule     *v1alpha1.Rule
	Resource interface{}
	Error    error
}

func New() engine.Engine[TfEngineRequest, TfEngineResponse] {
	type request struct {
		policy   *v1alpha1.Policy
		rule     *v1alpha1.Rule
		resource interface{}
	}
	looper := func(r TfEngineRequest) []request {
		var requests []request
		for _, resource := range r.Plan.Resources {
			for _, policy := range r.Policies {
				for _, rule := range policy.Spec.Rules {
					requests = append(requests, request{
						policy:   policy,
						rule:     &rule,
						resource: resource,
					})
				}
			}
		}
		return requests
	}
	inner := builder.
		Function(func(r request) TfEngineResponse {
			response := TfEngineResponse{
				Policy:   r.policy,
				Rule:     r.rule,
				Resource: r.resource,
			}
			if !match.Match(r.rule.Validation.Pattern, r.resource, match.WithWildcard(true)) {
				response.Error = errors.New(r.rule.Validation.Message)
			}
			return response
		}).
		Predicate(func(r request) bool { return !match.MatchResources(r.rule.ExcludeResources, r.resource) }).
		Predicate(func(r request) bool { return match.MatchResources(r.rule.MatchResources, r.resource) })
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}
