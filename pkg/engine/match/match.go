package match

import (
	"fmt"
	"reflect"

	"github.com/eddycharly/json-kyverno/pkg/apis/v1alpha1"
	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	"github.com/kyverno/kyverno/pkg/utils/wildcard"
)

func MatchResources(match *v1alpha1.MatchResources, actual interface{}, options ...option) (bool, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return false, nil
	}
	if len(match.Any) != 0 {
		if match, err := MatchAny(match.Any, actual, options...); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	if len(match.All) != 0 {
		if match, err := MatchAll(match.All, actual, options...); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	return true, nil
}

func MatchAny(filters v1alpha1.ResourceFilters, actual interface{}, options ...option) (bool, error) {
	for _, filter := range filters {
		if match, err := Match(filter.Resource, actual, options...); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}

func MatchAll(filters v1alpha1.ResourceFilters, actual interface{}, options ...option) (bool, error) {
	for _, filter := range filters {
		if match, err := Match(filter.Resource, actual, options...); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	return true, nil
}

func Match(expected, actual interface{}, options ...option) (bool, error) {
	return match(expected, actual, newMatchOptions(options...))
}

func match(expected, actual interface{}, options matchOptions) (bool, error) {
	if expected != nil {
		switch reflectutils.GetKind(expected) {
		case reflect.String:
			if options.wildcard {
				if reflectutils.GetKind(actual) != reflect.String {
					return false, fmt.Errorf("invalid actual value, must be a string, found %s", reflectutils.GetKind(actual))
				}
				return wildcard.Match(reflect.ValueOf(expected).String(), reflect.ValueOf(actual).String()), nil
			}
		case reflect.Slice:
			if reflectutils.GetKind(actual) != reflect.Slice {
				return false, fmt.Errorf("invalid actual value, must be a slice, found %s", reflectutils.GetKind(actual))
			}
			if reflect.ValueOf(expected).Len() != reflect.ValueOf(actual).Len() {
				return false, nil
			}
			for i := 0; i < reflect.ValueOf(expected).Len(); i++ {
				if inner, err := match(reflect.ValueOf(expected).Index(i).Interface(), reflect.ValueOf(actual).Index(i).Interface(), options); err != nil {
					return false, err
				} else if !inner {
					return false, nil
				}
			}
			return true, nil
		case reflect.Map:
			if reflectutils.GetKind(actual) != reflect.Map {
				return false, fmt.Errorf("invalid actual value, must be a map, found %s", reflectutils.GetKind(actual))
			}
			iter := reflect.ValueOf(expected).MapRange()
			for iter.Next() {
				actualValue := reflect.ValueOf(actual).MapIndex(iter.Key())
				if !actualValue.IsValid() {
					return false, nil
				}
				if inner, err := match(iter.Value().Interface(), actualValue.Interface(), options); err != nil {
					return false, err
				} else if !inner {
					return false, nil
				}
			}
			return true, nil
		}
	}
	return reflectutils.MatchScalar(expected, actual)
}
