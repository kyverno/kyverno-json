package match

import (
	"reflect"

	"github.com/eddycharly/tf-kyverno/pkg/apis/v1alpha1"
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
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}
	switch reflect.TypeOf(expected).Kind() {
	case reflect.String:
		if options.wildcard {
			return wildcard.Match(reflect.ValueOf(expected).String(), reflect.ValueOf(actual).String())
		}
	case reflect.Slice:
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
	return reflect.DeepEqual(expected, actual)
}
