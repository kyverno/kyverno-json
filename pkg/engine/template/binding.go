package template

import (
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type resolverFunc = func() (any, error)

type lazyBinding struct {
	resolver resolverFunc
}

func (b *lazyBinding) Value() (any, error) {
	return b.resolver()
}

func NewLazyBinding(resolver resolverFunc) binding.Binding {
	binding := &lazyBinding{}
	lock := &sync.Mutex{}
	binding.resolver = func() (any, error) {
		lock.Lock()
		defer lock.Unlock()
		value, err := resolver()
		binding.resolver = func() (any, error) {
			return value, err
		}
		return binding.resolver()
	}
	return binding
}

func NewLazyBindingWithValue(value any) binding.Binding {
	return NewLazyBinding(func() (any, error) {
		return value, nil
	})
}
