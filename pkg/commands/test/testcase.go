package test

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
)

type TestCase struct {
	Path  string
	Tests *v1alpha1.Tests
	Err   error
}

type TestCases []TestCase

func (tc TestCases) Errors() []TestCase {
	var errors []TestCase
	for _, test := range tc {
		if test.Err != nil {
			errors = append(errors, test)
		}
	}
	return errors
}
