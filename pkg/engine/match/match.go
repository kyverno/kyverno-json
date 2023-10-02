package match

import (
	"fmt"
	"reflect"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

func MatchResources(match *v1alpha1.MatchResources, actual interface{}) (bool, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return false, nil
	}
	if len(match.Any) != 0 {
		if match, err := MatchAny(match.Any, actual); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	if len(match.All) != 0 {
		if match, err := MatchAll(match.All, actual); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	return true, nil
}

func MatchAny(filters v1alpha1.ResourceFilters, actual interface{}) (bool, error) {
	for _, filter := range filters {
		if match, err := Match(filter.Resource, actual); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}

func MatchAll(filters v1alpha1.ResourceFilters, actual interface{}) (bool, error) {
	for _, filter := range filters {
		if match, err := Match(filter.Resource, actual); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	return true, nil
}

func Match(expected, actual interface{}) (bool, error) {
	return match(expected, actual)
}

func match(expected, actual interface{}) (bool, error) {
	if expected != nil {
		switch reflectutils.GetKind(expected) {
		case reflect.Slice:
			if reflectutils.GetKind(actual) != reflect.Slice {
				return false, fmt.Errorf("invalid actual value, must be a slice, found %s", reflectutils.GetKind(actual))
			}
			if reflect.ValueOf(expected).Len() != reflect.ValueOf(actual).Len() {
				return false, nil
			}
			for i := 0; i < reflect.ValueOf(expected).Len(); i++ {
				if inner, err := match(reflect.ValueOf(expected).Index(i).Interface(), reflect.ValueOf(actual).Index(i).Interface()); err != nil {
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
				if inner, err := match(iter.Value().Interface(), actualValue.Interface()); err != nil {
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
