package jsonengine

import (
	"context"
	"fmt"
	"time"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/builder"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/matching"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Request struct {
	Resource any
	Policies []*v1alpha1.ValidatingPolicy
	Bindings map[string]any
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
	Timestamp  time.Time
	Identifier string
	Feedback   map[string]Feedback
	Error      error
	Violations matching.Results
}

type Feedback struct {
	Error error
	Value any
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
			// TODO: this doesn't seem to be the right path
			var path *field.Path
			path = path.Child("context")
			for i, entry := range r.rule.Context {
				bindings = bindings.Register("$"+entry.Name, template.NewContextBinding(path.Index(i), bindings, r.resource, entry.Variable.Value()))
			}
			identifier := ""
			if r.rule.Identifier != "" {
				result, err := template.ExecuteJP(context.Background(), r.rule.Identifier, r.resource, bindings)
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
						Timestamp:  time.Now(),
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
						Timestamp:  time.Now(),
						Identifier: identifier,
						Error:      err,
					}}
				}
				// matched
				if len(errs) == 0 {
					return nil
				}
			}
			var feedback map[string]Feedback
			for _, f := range r.rule.Feedback {
				result, err := template.ExecuteJP(context.Background(), f.Value, r.resource, bindings)
				if feedback == nil {
					feedback = map[string]Feedback{}
				}
				if err != nil {
					feedback[f.Name] = Feedback{
						Error: err,
					}
				} else {
					feedback[f.Name] = Feedback{
						Value: result,
					}
				}
			}
			violations, err := matching.MatchAssert(ctx, nil, r.rule.Assert, r.resource, bindings)
			if err != nil {
				return []RuleResponse{{
					Rule:       r.rule,
					Timestamp:  time.Now(),
					Identifier: identifier,
					Feedback:   feedback,
					Error:      err,
				}}
			}
			return []RuleResponse{{
				Rule:       r.rule,
				Timestamp:  time.Now(),
				Identifier: identifier,
				Feedback:   feedback,
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
			bindings := jpbinding.NewBindings()
			for k, v := range r.Bindings {
				bindings = bindings.Register("$"+k, jpbinding.NewBinding(v))
			}
			bindings = bindings.Register("$payload", jpbinding.NewBinding(r.Resource))
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
