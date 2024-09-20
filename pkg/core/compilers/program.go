package compilers

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type Program func(any, binding.Bindings) (any, error)
