package template

import (
	"sync"
)

type lazyBinding func() (any, error)

func (b lazyBinding) Value() (any, error) {
	return b()
}

func NewLazyBinding(resolver func() (any, error)) lazyBinding {
	return sync.OnceValues(resolver)
}
