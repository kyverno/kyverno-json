package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// TODO: this is the only file we reference the apis package.
// We should remove this dependency.
// Either move the Match struct in this package or move this file in a more specific package.

func Match(ctx context.Context, path *field.Path, match *v1alpha1.Match, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		errs = append(errs, field.Invalid(path, match, "an empty match is not valid"))
	} else {
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
	}
	return errs, nil
}

func MatchAny(ctx context.Context, path *field.Path, assertions v1alpha1.Assertions, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := validate(ctx, path.Index(i), assertion.Value, actual, bindings)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}

func MatchAll(ctx context.Context, path *field.Path, assertions v1alpha1.Assertions, actual interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := validate(ctx, path.Index(i), assertion.Value, actual, bindings)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}
