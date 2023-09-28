package tfengine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/eddycharly/tf-kyverno/pkg/apis/v1alpha1"
	"github.com/eddycharly/tf-kyverno/pkg/engine"
	"github.com/eddycharly/tf-kyverno/pkg/engine/blocks/loop"
	"github.com/eddycharly/tf-kyverno/pkg/engine/builder"
	"github.com/eddycharly/tf-kyverno/pkg/match"
	"github.com/eddycharly/tf-kyverno/pkg/plan"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
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

var (
	variables = regexp.MustCompile(`\{\{(.*?)\}\}`)
	parser    = parsing.NewParser()
	caller    = interpreter.NewFunctionCaller(functions.GetDefaultFunctions()...)
)

func jp(statement string, policy *v1alpha1.Policy, rule *v1alpha1.Rule, resource interface{}) (interface{}, error) {
	bindings := binding.NewBindings()
	for _, entry := range rule.Context {
		bindings = bindings.Register("$"+entry.Name, entry.Variable.Value)
	}
	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"policy":   policy,
		"rule":     rule,
		"resource": resource,
	}
	interpreter := interpreter.NewInterpreter(data, caller, bindings)
	return interpreter.Execute(compiled, data)
}

func message(message string, policy *v1alpha1.Policy, rule *v1alpha1.Rule, resource interface{}) string {
	groups := variables.FindAllStringSubmatch(message, -1)
	for _, group := range groups {
		statement := strings.TrimSpace(group[1])
		result, err := jp(statement, policy, rule, resource)
		if err != nil {
			message = strings.ReplaceAll(message, group[0], fmt.Sprintf("ERR (%s - %s)", statement, err))
		} else if result == nil {
			message = strings.ReplaceAll(message, group[0], fmt.Sprintf("ERR (%s not found)", statement))
		} else if result, ok := result.(string); !ok {
			message = strings.ReplaceAll(message, group[0], fmt.Sprintf("ERR (%s not a string)", statement))
		} else {
			message = strings.ReplaceAll(message, group[0], result)
		}
	}
	return message
}

func New() engine.Engine[TfEngineRequest, TfEngineResponse] {
	type request struct {
		Policy   *v1alpha1.Policy
		Rule     *v1alpha1.Rule
		Resource interface{}
	}
	looper := func(r TfEngineRequest) []request {
		var requests []request
		for _, resource := range r.Plan.Resources {
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
		Function(func(r request) TfEngineResponse {
			response := TfEngineResponse{
				Policy:   r.Policy,
				Rule:     r.Rule,
				Resource: r.Resource,
			}
			if !match.Match(r.Rule.Validation.Pattern, r.Resource, match.WithWildcard()) {
				response.Error = errors.New(message(r.Rule.Validation.Message, r.Policy, r.Rule, r.Resource))
			}
			return response
		}).
		Predicate(func(r request) bool {
			return !match.MatchResources(r.Rule.ExcludeResources, r.Resource, match.WithWildcard())
		}).
		Predicate(func(r request) bool {
			return match.MatchResources(r.Rule.MatchResources, r.Resource, match.WithWildcard())
		})
	// TODO: we can't use the builder package for loops :(
	return loop.New(inner, looper)
}
