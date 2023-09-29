package jsonengine

import (
	"errors"

	"github.com/eddycharly/json-kyverno/pkg/apis/v1alpha1"
	"github.com/eddycharly/json-kyverno/pkg/engine"
	"github.com/eddycharly/json-kyverno/pkg/engine/blocks/loop"
	"github.com/eddycharly/json-kyverno/pkg/engine/builder"
	"github.com/eddycharly/json-kyverno/pkg/engine/match"
	"github.com/eddycharly/json-kyverno/pkg/engine/template"
)

type JsonEngineRequest struct {
	Resources []interface{}
	Policies  []*v1alpha1.Policy
}

type JsonEngineResponse struct {
	Policy   *v1alpha1.Policy
	Rule     *v1alpha1.Rule
	Resource interface{}
	Error    error
}

func New() engine.Engine[JsonEngineRequest, JsonEngineResponse] {
	type request struct {
		Policy   *v1alpha1.Policy
		Rule     *v1alpha1.Rule
		Resource interface{}
	}
	looper := func(r JsonEngineRequest) []request {
		var requests []request
		for _, resource := range r.Resources {
			for _, policy := range r.Policies {
				for _, rule := range policy.Spec.Rules {
					requests = append(requests, request{
						Policy:   policy,
						Rule:     &rule,
						Resource: resource,
					})
				}
			}
		}
		return requests
	}
	inner := builder.
		Function(func(r request) JsonEngineResponse {
			response := JsonEngineResponse{
				Policy:   r.Policy,
				Rule:     r.Rule,
				Resource: r.Resource,
			}
			template := template.New(r, r.Rule.Context...)
			if !match.Match(r.Rule.Validation.Pattern, r.Resource, match.WithWildcard(), match.WithTemplate(template)) {
				message := r.Rule.Validation.Message
				if message != "" {
					message = template.String(message)
				} else {
					message = "failed to match pattern"
				}
				response.Error = errors.New(message)
			}
			return response
		}).
		Predicate(func(r request) bool {
			return !match.MatchResources(r.Rule.ExcludeResources, r.Resource, match.WithWildcard())
		}).
		Predicate(func(r request) bool {
			return r.Rule.MatchResources == nil || match.MatchResources(r.Rule.MatchResources, r.Resource, match.WithWildcard())
		})
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}
