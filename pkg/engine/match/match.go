package match

import (
	"reflect"

	"github.com/eddycharly/json-kyverno/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno/pkg/utils/wildcard"
)

func MatchResources(match *v1alpha1.MatchResources, actual interface{}, options ...option) bool {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return false
	}
	if len(match.Any) != 0 {
		if !MatchAny(match.Any, actual, options...) {
			return false
		}
	}
	if len(match.All) != 0 {
		if !MatchAll(match.All, actual, options...) {
			return false
		}
	}
	return true
}

func MatchAny(filters v1alpha1.ResourceFilters, actual interface{}, options ...option) bool {
	for _, filter := range filters {
		if Match(filter.Resource, actual, options...) {
			return true
		}
	}
	return false
}

func MatchAll(filters v1alpha1.ResourceFilters, actual interface{}, options ...option) bool {
	for _, filter := range filters {
		if !Match(filter.Resource, actual, options...) {
			return false
		}
	}
	return true
}

func Match(expected, actual interface{}, options ...option) bool {
	return match(expected, actual, newMatchOptions(options...))
}

func match(expected, actual interface{}, options matchOptions) bool {
	if options.template != nil && reflect.TypeOf(expected).Kind() == reflect.String {
		expected = options.template.Interface(reflect.ValueOf(expected).String())
	}
	// if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
	// 	return false
	// }
	switch reflect.TypeOf(expected).Kind() {
	case reflect.String:
		if options.wildcard {
			if reflect.TypeOf(actual).Kind() != reflect.String {
				return false
			}
			return wildcard.Match(reflect.ValueOf(expected).String(), reflect.ValueOf(actual).String())
		}
	case reflect.Slice:
		if reflect.TypeOf(actual).Kind() != reflect.Slice {
			return false
		}
		if reflect.ValueOf(expected).Len() != reflect.ValueOf(actual).Len() {
			return false
		}
		for i := 0; i < reflect.ValueOf(expected).Len(); i++ {
			if !match(reflect.ValueOf(expected).Index(i).Interface(), reflect.ValueOf(actual).Index(i).Interface(), options) {
				return false
			}
		}
		return true
	case reflect.Map:
		if reflect.TypeOf(actual).Kind() != reflect.Map {
			return false
		}
		iter := reflect.ValueOf(expected).MapRange()
		for iter.Next() {
			actualValue := reflect.ValueOf(actual).MapIndex(iter.Key())
			if !actualValue.IsValid() {
				return false
			}
			if !match(iter.Value().Interface(), actualValue.Interface(), options) {
				return false
			}
		}
		return true
	}
	return matchScalar(expected, actual)
}

func matchScalar(expected, actual interface{}) bool {
	// if they are the same type we can use reflect.DeepEqual
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) {
		return reflect.DeepEqual(expected, actual)
	}
	e := reflect.ValueOf(expected)
	a := reflect.ValueOf(actual)
	if !a.IsValid() && !e.IsValid() {
		return true
	}
	if a.IsZero() && e.IsZero() {
		return true
	}
	if a.CanComplex() && e.CanComplex() {
		return a.Complex() == e.Complex()
	}
	if a.CanFloat() && e.CanFloat() {
		return a.Float() == e.Float()
	}
	if a.CanInt() && e.CanInt() {
		return a.Int() == e.Int()
	}
	if a.CanUint() && e.CanUint() {
		return a.Uint() == e.Uint()
	}
	return false
}
