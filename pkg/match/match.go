package match

import (
	"reflect"

	"github.com/eddycharly/tf-kyverno/pkg/apis/v1alpha1"
)

func MatchResources(match *v1alpha1.MatchResources, actual interface{}) bool {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return false
	}
	if len(match.Any) != 0 {
		if !MatchAny(match.Any, actual) {
			return false
		}
	}
	if len(match.All) != 0 {
		if !MatchAll(match.All, actual) {
			return false
		}
	}
	return true
}

func MatchAny(filters v1alpha1.ResourceFilters, actual interface{}) bool {
	for _, filter := range filters {
		if Match(filter.Resource, actual) {
			return true
		}
	}
	return false
}

func MatchAll(filters v1alpha1.ResourceFilters, actual interface{}) bool {
	for _, filter := range filters {
		if !Match(filter.Resource, actual) {
			return false
		}
	}
	return true
}

func Match(expected, actual interface{}) bool {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	switch reflect.TypeOf(expected).Kind() {
	case reflect.Slice:
		if reflect.ValueOf(expected).Len() != reflect.ValueOf(actual).Len() {
			return false
		}
		for i := 0; i < reflect.ValueOf(expected).Len(); i++ {
			if !Match(reflect.ValueOf(expected).Index(i).Interface(), reflect.ValueOf(actual).Index(i).Interface()) {
				return false
			}
		}
	case reflect.Map:
		iter := reflect.ValueOf(expected).MapRange()

		for iter.Next() {
			actualValue := reflect.ValueOf(actual).MapIndex(iter.Key())
			if !actualValue.IsValid() {
				return false
			}
			if !Match(iter.Value().Interface(), actualValue.Interface()) {
				return false
			}
		}
	default:
		return false
	}
	return true
}
