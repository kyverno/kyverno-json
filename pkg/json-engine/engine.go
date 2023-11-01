package jsonengine

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/loop"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"go.uber.org/multierr"
)

type JsonEngineRequest struct {
	Resources []interface{}
	Policies  []*v1alpha1.ValidatingPolicy
}

type JsonEngineResponse struct {
	Policy   *v1alpha1.ValidatingPolicy
	Rule     v1alpha1.ValidatingRule
	Resource interface{}
	Failure  error
	Error    error
}

func New() engine.Engine[JsonEngineRequest, JsonEngineResponse] {
	type request struct {
		policy   *v1alpha1.ValidatingPolicy
		rule     v1alpha1.ValidatingRule
		value    interface{}
		bindings binding.Bindings
	}
	looper := func(r JsonEngineRequest) []request {
		var requests []request
		bindings := jpbinding.NewBindings()
		for _, resource := range r.Resources {
			bindings = bindings.Register("$payload", jpbinding.NewBinding(resource))
			for _, policy := range r.Policies {
				bindings = bindings.Register("$policy", jpbinding.NewBinding(policy))
				for _, rule := range policy.Spec.Rules {
					bindings = bindings.Register("$rule", jpbinding.NewBinding(rule))
					bindings = assert.NewContextBindings(bindings, resource, rule.Context...)
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
		Function(func(ctx context.Context, r request) JsonEngineResponse {
			response := JsonEngineResponse{
				Policy:   r.policy,
				Rule:     r.rule,
				Resource: r.value,
			}
			errs, err := assert.MatchAssert(ctx, nil, r.rule.Assert, r.value, r.bindings)
			if err != nil {
				response.Error = err
			} else if err := multierr.Combine(errs...); err != nil {
				response.Failure = err
			}
			return response
		}).
		Predicate(func(ctx context.Context, r request) bool {
			if r.rule.Exclude == nil {
				return true
			}
			errs, err := assert.Match(ctx, nil, r.rule.Exclude, r.value, r.bindings)
			return err == nil && len(errs) != 0
		}).
		Predicate(func(ctx context.Context, r request) bool {
			if r.rule.Match == nil {
				return true
			}
			errs, err := assert.Match(ctx, nil, r.rule.Match, r.value, r.bindings)
			return err == nil && len(errs) == 0
		})
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}
