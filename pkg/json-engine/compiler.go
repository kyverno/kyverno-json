package jsonengine

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func CompileContextEntry(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.ContextEntry,
) (func(any, binding.Bindings) binding.Bindings, *field.Error) {
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	handler, err := in.Variable.Compile(path.Child("variable"), compilers)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) binding.Bindings {
		binding := binding.NewDelegate(
			sync.OnceValues(
				func() (any, error) {
					projected, err := handler(resource, bindings)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				},
			),
		)
		return bindings.Register("$"+in.Name, binding)
	}, nil
}

func CompileContextEntries(
	path *field.Path,
	compilers compilers.Compilers,
	in ...v1alpha1.ContextEntry,
) (func(any, binding.Bindings) binding.Bindings, *field.Error) {
	var out []func(any, binding.Bindings) binding.Bindings
	for i, entry := range in {
		entry, err := CompileContextEntry(path.Index(i), compilers, entry)
		if err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return func(resource any, bindings binding.Bindings) binding.Bindings {
		for _, entry := range out {
			bindings = entry(resource, bindings)
		}
		return bindings
	}, nil
}

func CompileMatch(
	path *field.Path,
	compilers compilers.Compilers,
	in *v1alpha1.Match,
) (func(any, binding.Bindings) (field.ErrorList, error), *field.Error) {
	if in == nil {
		return nil, nil
	}
	if len(in.Any) == 0 && len(in.All) == 0 {
		return nil, field.Invalid(path, in, "an empty match is not valid")
	}
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	_any, err := CompileAssertionTrees(path.Child("any"), compilers, in.Any...)
	if err != nil {
		return nil, err
	}
	_all, err := CompileAssertionTrees(path.Child("all"), compilers, in.All...)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) (field.ErrorList, error) {
		var errs field.ErrorList
		for _, assertion := range _any {
			_errs, err := assertion(resource, bindings)
			if err != nil {
				return errs, err
			}
			if len(_errs) == 0 {
				return nil, nil
			}
			errs = append(errs, _errs...)
		}
		for _, assertion := range _all {
			_errs, err := assertion(resource, bindings)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		return errs, nil
	}, nil
}

func CompileAssert(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.Assert,
) (func(any, binding.Bindings) (Results, error), *field.Error) {
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	if len(in.Any) == 0 && len(in.All) == 0 {
		return nil, field.Invalid(path, in, "an empty assert is not valid")
	}
	_any, err := CompileAssertions(path.Child("any"), compilers, in.Any...)
	if err != nil {
		return nil, err
	}
	_all, err := CompileAssertions(path.Child("all"), compilers, in.All...)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) (Results, error) {
		if len(_any) != 0 {
			var fails Results
			for _, assertion := range _any {
				result, err := assertion(resource, bindings)
				if err != nil {
					return fails, err
				}
				if len(result.ErrorList) == 0 {
					fails = nil
					break
				}
				fails = append(fails, result)
			}
			if fails != nil {
				return fails, nil
			}
		}
		if len(_all) != 0 {
			var fails Results
			for _, assertion := range _all {
				result, err := assertion(resource, bindings)
				if err != nil {
					return fails, err
				}
				if len(result.ErrorList) > 0 {
					fails = append(fails, result)
				}
			}
			return fails, nil
		}
		return nil, nil
	}, nil
}

func CompileAssertions(
	path *field.Path,
	compilers compilers.Compilers,
	in ...v1alpha1.Assertion,
) ([]func(any, binding.Bindings) (Result, error), *field.Error) {
	var out []func(any, binding.Bindings) (Result, error)
	for i, in := range in {
		if in, err := CompileAssertion(path.Index(i), compilers, in); err != nil {
			return nil, err
		} else {
			out = append(out, in)
		}
	}
	return out, nil
}

func CompileAssertion(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.Assertion,
) (func(any, binding.Bindings) (Result, error), *field.Error) {
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	check, err := CompileAssertionTree(path.Child("check"), compilers, in.Check)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) (Result, error) {
		var result Result
		errs, err := check(resource, bindings)
		if len(errs) != 0 {
			result.ErrorList = errs
			message := fmt.Sprintf("(CHECK=%s)", path.String())
			if in.Message != nil {
				message = fmt.Sprintf("%s %s", in.Message.Format(resource, bindings, compilers.Jp.Options()...), message)
			}
			result.Message = message
		}
		return result, err
	}, nil
}

func CompileAssertionTrees(
	path *field.Path,
	compilers compilers.Compilers,
	in ...v1alpha1.AssertionTree,
) ([]func(any, binding.Bindings) (field.ErrorList, error), *field.Error) {
	var out []func(any, binding.Bindings) (field.ErrorList, error)
	for i, in := range in {
		if in, err := CompileAssertionTree(path.Index(i), compilers, in); err != nil {
			return nil, err
		} else {
			out = append(out, in)
		}
	}
	return out, nil
}

func CompileAssertionTree(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.AssertionTree,
) (func(any, binding.Bindings) (field.ErrorList, error), *field.Error) {
	check, err := in.Compile(path, compilers)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) (field.ErrorList, error) {
		return check.Assert(nil, resource, bindings)
	}, nil
}

func CompileIdentifier(
	path *field.Path,
	compilers compilers.Compilers,
	in string,
) (func(any, binding.Bindings) string, *field.Error) {
	if in == "" {
		return func(resource any, bindings binding.Bindings) string {
			return ""
		}, nil
	}
	program, err := compilers.Jp.Compile(in)
	if err != nil {
		return nil, field.InternalError(path, err)
	}
	return func(resource any, bindings binding.Bindings) string {
		result, err := program(resource, bindings)
		if err != nil {
			return fmt.Sprintf("(error: %s)", err)
		} else {
			return fmt.Sprint(result)
		}
	}, nil
}

func CompileFeedbacks(
	path *field.Path,
	compilers compilers.Compilers,
	in ...v1alpha1.Feedback,
) (func(any, binding.Bindings) map[string]Feedback, *field.Error) {
	if len(in) == 0 {
		return func(any, binding.Bindings) map[string]Feedback {
			return nil
		}, nil
	}
	feedback := map[string]func(any, binding.Bindings) Feedback{}
	for i, in := range in {
		f, err := CompileFeedback(path.Index(i), compilers, in)
		if err != nil {
			return nil, err
		}
		feedback[in.Name] = f
	}
	return func(resource any, bindings binding.Bindings) map[string]Feedback {
		out := map[string]Feedback{}
		for name, f := range feedback {
			out[name] = f(resource, bindings)
		}
		return out
	}, nil
}

func CompileFeedback(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.Feedback,
) (func(any, binding.Bindings) Feedback, *field.Error) {
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	handler, err := in.Value.Compile(path.Child("value"), compilers)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) Feedback {
		var out Feedback
		if projected, err := handler(resource, bindings); err != nil {
			out.Error = err
		} else {
			out.Value = projected
		}
		return out
	}, nil
}

func CompileRule(
	path *field.Path,
	compilers compilers.Compilers,
	in v1alpha1.ValidatingRule,
) (func(any, binding.Bindings) *RuleResponse, *field.Error) {
	if in.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*in.Compiler))
	}
	context, err := CompileContextEntries(path.Child("context"), compilers, in.Context...)
	if err != nil {
		return nil, err
	}
	identifier, err := CompileIdentifier(path.Child("identifier"), compilers, in.Identifier)
	if err != nil {
		return nil, err
	}
	match, err := CompileMatch(path.Child("match"), compilers, in.Match)
	if err != nil {
		return nil, err
	}
	exclude, err := CompileMatch(path.Child("exclude"), compilers, in.Exclude)
	if err != nil {
		return nil, err
	}
	feedback, err := CompileFeedbacks(path.Child("feedback"), compilers, in.Feedback...)
	if err != nil {
		return nil, err
	}
	assert, err := CompileAssert(path.Child("assert"), compilers, in.Assert)
	if err != nil {
		return nil, err
	}
	return func(resource any, bindings binding.Bindings) *RuleResponse {
		// register context bindings
		bindings = context(resource, bindings)
		// process match clause
		if match != nil {
			if errs, err := match(resource, bindings); err != nil {
				return &RuleResponse{
					Rule:       in,
					Timestamp:  time.Now(),
					Identifier: identifier(resource, bindings),
					Feedback:   feedback(resource, bindings),
					Error:      err,
				}
			} else if len(errs) != 0 {
				// didn't match
				return nil
			}
		}
		// process exclude clause
		if exclude != nil {
			if errs, err := exclude(resource, bindings); err != nil {
				return &RuleResponse{
					Rule:       in,
					Timestamp:  time.Now(),
					Identifier: identifier(resource, bindings),
					Feedback:   feedback(resource, bindings),
					Error:      err,
				}
			} else if len(errs) == 0 {
				// matched
				return nil
			}
		}
		// evaluate assertions
		violations, err := assert(resource, bindings)
		if err != nil {
			return &RuleResponse{
				Rule:       in,
				Timestamp:  time.Now(),
				Identifier: identifier(resource, bindings),
				Feedback:   feedback(resource, bindings),
				Error:      err,
			}
		}
		return &RuleResponse{
			Rule:       in,
			Timestamp:  time.Now(),
			Identifier: identifier(resource, bindings),
			Feedback:   feedback(resource, bindings),
			Violations: violations,
		}
	}, nil
}
