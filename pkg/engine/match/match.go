package match

import (
	"fmt"
	"reflect"

	"github.com/eddycharly/json-kyverno/pkg/apis/v1alpha1"
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
	if options.template != nil && reflect.TypeOf(expected).Kind() == reflect.String {
		result, err := options.template.Interface(reflect.ValueOf(expected).String(), actual)
		if err != nil {
			return false, err
		}
		expected = result
	}
	if expected != nil {
		switch reflect.TypeOf(expected).Kind() {
		case reflect.String:
			if options.wildcard {
				if reflect.TypeOf(actual).Kind() != reflect.String {
					return false, fmt.Errorf("invalid actual value, must be a string, found %s", reflect.TypeOf(actual).Kind())
				}
				return wildcard.Match(reflect.ValueOf(expected).String(), reflect.ValueOf(actual).String()), nil
			}
		case reflect.Slice:
			if reflect.TypeOf(actual).Kind() != reflect.Slice {
				return false, fmt.Errorf("invalid actual value, must be a slice, found %s", reflect.TypeOf(actual).Kind())
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
			keyType := reflect.TypeOf(expected).Key()
			if options.template != nil && keyType.Kind() == reflect.String {
				iter := reflect.ValueOf(expected).MapRange()
				for iter.Next() {
					key := iter.Key().String()
					result, err := options.template.Interface(key, actual)
					if err != nil {
						return false, err
					}
					switch actualKey := result.(type) {
					case string:
						if actualKey == key {
							if reflect.TypeOf(actual).Kind() != reflect.Map {
								return false, fmt.Errorf("invalid actual value, must be a map, found %s", reflect.TypeOf(actual).Kind())
							}
							actualValue := reflect.ValueOf(actual).MapIndex(iter.Key())
							if !actualValue.IsValid() {
								return false, nil
							}
							if inner, err := match(iter.Value().Interface(), actualValue.Interface(), options); err != nil {
								return false, err
							} else if !inner {
								return false, nil
							}
						} else {
							if inner, err := match(iter.Value().Interface(), actualKey, options); err != nil {
								return false, err
							} else if !inner {
								return false, nil
							}
						}
					default:
						if inner, err := match(iter.Value().Interface(), actualKey, options); err != nil {
							return false, err
						} else if !inner {
							return false, nil
						}
					}
				}
			} else {
				if reflect.TypeOf(actual).Kind() != reflect.Map {
					return false, fmt.Errorf("invalid actual value, must be a map, found %s", reflect.TypeOf(actual).Kind())
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
			}
			return true, nil
		}
	}
	return matchScalar(expected, actual)
}

func getKind(value interface{}) reflect.Kind {
	if value == nil {
		return reflect.Invalid
	}
	return reflect.TypeOf(value).Kind()
}

func toNumber(value reflect.Value) (float64, bool) {
	if value.CanFloat() {
		return value.Float(), true
	}
	if value.CanInt() {
		return float64(value.Int()), true
	}
	if value.CanUint() {
		return float64(value.Uint()), true
	}
	return 0, false
}

func matchScalar(expected, actual interface{}) (bool, error) {
	if actual == nil && expected == nil {
		return true, nil
	}
	if actual == expected {
		return true, nil
	}
	// if they are the same type we can use reflect.DeepEqual
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) {
		return reflect.DeepEqual(expected, actual), nil
	}
	e := reflect.ValueOf(expected)
	a := reflect.ValueOf(actual)
	if !a.IsValid() && !e.IsValid() {
		return true, nil
	}
	if a.CanComplex() && e.CanComplex() {
		return a.Complex() == e.Complex(), nil
	}
	if a.CanFloat() && e.CanFloat() {
		return a.Float() == e.Float(), nil
	}
	if a.CanInt() && e.CanInt() {
		return a.Int() == e.Int(), nil
	}
	if a.CanUint() && e.CanUint() {
		return a.Uint() == e.Uint(), nil
	}
	if a, ok := toNumber(a); ok {
		if e, ok := toNumber(e); ok {
			return a == e, nil
		}
	}
	return false, fmt.Errorf("types are not comparable, %s - %s", getKind(expected), getKind(actual))
}
