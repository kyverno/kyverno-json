package jsonengine

import (
	"errors"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/loop"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"github.com/kyverno/kyverno-json/pkg/engine/match"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

type JsonEngineRequest struct {
	Resources []interface{}
	Policies  []*v1alpha1.Policy
}

type JsonEngineResponse struct {
	Policy   *v1alpha1.Policy
	Rule     *v1alpha1.Rule
	Resource interface{}
	Failure  error
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
			bindings := jpbinding.NewBindings()
			bindings = bindings.Register("$resource", jpbinding.NewBinding(r.Resource))
			bindings = bindings.Register("$rule", jpbinding.NewBinding(r.Rule))
			bindings = bindings.Register("$policy", jpbinding.NewBinding(r.Policy))
			bindings = assert.NewContextBindings(bindings, r.Resource, r.Rule.Context...)
			errs, err := assert.Assert(r.Rule.Validation.Pattern, r.Resource, bindings)
			if err != nil {
				response.Failure = err
			} else if err := errs.ToAggregate(); err != nil {
				if r.Rule.Validation.Message != "" {
					response.Error = errors.New(template.String(r.Rule.Validation.Message, r.Resource, bindings))
				} else {
					response.Error = err
				}
			}
			return response
		}).
		Predicate(func(r request) bool {
			match, err := match.MatchResources(r.Rule.ExcludeResources, r.Resource, match.WithWildcard())
			return err == nil && !match
		}).
		Predicate(func(r request) bool {
			if r.Rule.MatchResources == nil {
				return true
			}
			match, err := match.MatchResources(r.Rule.MatchResources, r.Resource, match.WithWildcard())
			return err == nil && match
		})
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}
