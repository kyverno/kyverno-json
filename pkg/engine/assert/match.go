package assert

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// TODO: this is the only file we reference the apis package.
// We should remove this dependency.
// Either move the Match struct in this package or move this file in a more specific package.

func MatchAssert(ctx context.Context, path *field.Path, match *v1alpha1.Assert, actual interface{}, bindings binding.Bindings) ([]error, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return nil, field.Invalid(path, match, "an empty assert is not valid")
	} else {
		if len(match.Any) != 0 {
			var fails []error
			path := path.Child("any")
			for i, assertion := range match.Any {
				checkFails, err := assert(ctx, path.Index(i).Child("check"), Parse(ctx, assertion.Check.Value), actual, bindings)
				if err != nil {
					return fails, err
				}
				if len(checkFails) == 0 {
					fails = nil
					break
				}
				if assertion.Message != "" {
					msg := template.String(ctx, assertion.Message, actual, bindings)
					msg += ": " + checkFails.ToAggregate().Error()
					fails = append(fails, errors.New(msg))
				} else {
					fails = append(fails, checkFails.ToAggregate())
				}
			}
			if fails != nil {
				return fails, nil
			}
		}
		if len(match.All) != 0 {
			var fails []error
			path := path.Child("all")
			for i, assertion := range match.All {
				checkFails, err := assert(ctx, path.Index(i).Child("check"), Parse(ctx, assertion.Check.Value), actual, bindings)
				if err != nil {
					return fails, err
				}
				if len(checkFails) > 0 {
					if assertion.Message != "" {
						msg := template.String(ctx, assertion.Message, actual, bindings)
						msg += ": " + checkFails.ToAggregate().Error()
						fails = append(fails, errors.New(msg))
					} else {
						fails = append(fails, checkFails.ToAggregate())
					}
				}
			}
			return fails, nil
		}
		return nil, nil
	}
}

func Match(ctx context.Context, path *field.Path, match *v1alpha1.Match, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return nil, field.Invalid(path, match, "an empty match is not valid")
	} else {
		var errs field.ErrorList
		if len(match.Any) != 0 {
			_errs, err := MatchAny(ctx, path.Child("any"), match.Any, actual, bindings)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		if len(match.All) != 0 {
			_errs, err := MatchAll(ctx, path.Child("all"), match.All, actual, bindings)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		return errs, nil
	}
}

func MatchAny(ctx context.Context, path *field.Path, assertions []v1alpha1.Any, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := assert(ctx, path.Index(i), Parse(ctx, assertion.Value), actual, bindings)
		if err != nil {
			return errs, err
		}
		if len(_errs) == 0 {
			return nil, nil
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}

func MatchAll(ctx context.Context, path *field.Path, assertions []v1alpha1.Any, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := assert(ctx, path.Index(i), Parse(ctx, assertion.Value), actual, bindings)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}
