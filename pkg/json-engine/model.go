package jsonengine

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
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
	Policy v1alpha1.ValidatingPolicy
	Rules  []RuleResponse
}

type RuleResponse struct {
	Rule       v1alpha1.ValidatingRule
	Timestamp  time.Time
	Identifier string
	Feedback   map[string]Feedback
	Error      error
	Violations Results
}

//nolint:errname
type Result struct {
	field.ErrorList
	Message string
}

func (r Result) Error(prefix string) string {
	var lines []string
	if r.Message != "" {
		lines = append(lines, prefix+"-> "+r.Message)
	}
	for _, err := range r.ErrorList {
		lines = append(lines, prefix+fmt.Sprintf(" -> %s (PATH=%s)", err.ErrorBody(), err.Field))
	}
	return strings.Join(lines, "\n")
}

//nolint:errname
type Results []Result

func (r Results) Error(prefix string) string {
	var lines []string
	for _, err := range r {
		lines = append(lines, err.Error(prefix))
	}
	return strings.Join(lines, "\n")
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
