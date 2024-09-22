package matching

import (
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

//nolint:errname
type Result struct {
	field.ErrorList
	Message string
}

func (r Result) Error() string {
	var lines []string
	if r.Message != "" {
		lines = append(lines, "-> "+r.Message)
	}
	for _, err := range r.ErrorList {
		lines = append(lines, " -> "+err.Error())
	}
	return strings.Join(lines, "\n")
}

//nolint:errname
type Results []Result

func (r Results) Error() string {
	var lines []string
	for _, err := range r {
		lines = append(lines, err.Error())
	}
	return strings.Join(lines, "\n")
}

func Assert(path *field.Path, in v1alpha1.Assert, actual any, bindings binding.Bindings, compiler Compiler, defaultCompiler string) ([]Result, error) {
	if in.Compiler != nil {
		defaultCompiler = string(*in.Compiler)
	}
	if len(in.Any) == 0 && len(in.All) == 0 {
		return nil, field.Invalid(path, in, "an empty assert is not valid")
	} else {
		if len(in.Any) != 0 {
			var fails []Result
			path := path.Child("any")
			for i, assertion := range in.Any {
				defaultCompiler := defaultCompiler
				if assertion.Compiler != nil {
					defaultCompiler = string(*assertion.Compiler)
				}
				checkFails, err := assert(path.Index(i).Child("check"), assertion.Check, actual, bindings, compiler, defaultCompiler)
				if err != nil {
					return fails, err
				}
				if len(checkFails) == 0 {
					fails = nil
					break
				}
				fail := Result{
					ErrorList: checkFails,
				}
				if assertion.Message != nil {
					fail.Message = assertion.Message.Format(actual, bindings, compiler.Jp.Options()...)
				}
				fails = append(fails, fail)
			}
			if fails != nil {
				return fails, nil
			}
		}
		if len(in.All) != 0 {
			var fails []Result
			path := path.Child("all")
			for i, assertion := range in.All {
				defaultCompiler := defaultCompiler
				if assertion.Compiler != nil {
					defaultCompiler = string(*assertion.Compiler)
				}
				checkFails, err := assert(path.Index(i).Child("check"), assertion.Check, actual, bindings, compiler, defaultCompiler)
				if err != nil {
					return fails, err
				}
				if len(checkFails) > 0 {
					fail := Result{
						ErrorList: checkFails,
					}
					if assertion.Message != nil {
						fail.Message = assertion.Message.Format(actual, bindings, compiler.Jp.Options()...)
					}
					fails = append(fails, fail)
				}
			}
			return fails, nil
		}
		return nil, nil
	}
}

func Match(path *field.Path, in *v1alpha1.Match, actual any, bindings binding.Bindings, compiler Compiler, defaultCompiler string) (field.ErrorList, error) {
	if in.Compiler != nil {
		defaultCompiler = string(*in.Compiler)
	}
	if in == nil || (len(in.Any) == 0 && len(in.All) == 0) {
		return nil, field.Invalid(path, in, "an empty match is not valid")
	} else {
		var errs field.ErrorList
		if len(in.Any) != 0 {
			_errs, err := matchAny(path.Child("any"), in.Any, actual, bindings, compiler, defaultCompiler)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		if len(in.All) != 0 {
			_errs, err := matchAll(path.Child("all"), in.All, actual, bindings, compiler, defaultCompiler)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		return errs, nil
	}
}

func matchAny(path *field.Path, in []v1alpha1.AssertionTree, value any, bindings binding.Bindings, compiler Compiler, defaultCompiler string) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range in {
		_errs, err := assert(path.Index(i), assertion, value, bindings, compiler, defaultCompiler)
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

func matchAll(path *field.Path, in []v1alpha1.AssertionTree, value any, bindings binding.Bindings, compiler Compiler, defaultCompiler string) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range in {
		_errs, err := assert(path.Index(i), assertion, value, bindings, compiler, defaultCompiler)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}

func assert(path *field.Path, assertion v1alpha1.AssertionTree, value any, bindings binding.Bindings, compiler Compiler, defaultCompiler string) (field.ErrorList, error) {
	check, err := assertion.Compile(compiler.CompileAssertion, defaultCompiler)
	if err != nil {
		return nil, err
	}
	return check.Assert(path, value, bindings)
}
