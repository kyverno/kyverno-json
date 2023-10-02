package template

import (
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type resolverFunc = func() (interface{}, error)

type lazyBinding struct {
	resolver resolverFunc
}

func (b *lazyBinding) Value() (interface{}, error) {
	return b.resolver()
}

func NewLazyBinding(resolver resolverFunc) binding.Binding {
	binding := &lazyBinding{}
	lock := &sync.Mutex{}
	binding.resolver = func() (interface{}, error) {
		lock.Lock()
		defer lock.Unlock()
		value, err := resolver()
		binding.resolver = func() (interface{}, error) {
			return value, err
		}
		return binding.resolver()
	}
	return binding
}

func NewLazyBindingWithValue(value interface{}) binding.Binding {
	return NewLazyBinding(func() (interface{}, error) {
		return value, nil
	})
}
