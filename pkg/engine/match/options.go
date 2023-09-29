package match

import (
	"github.com/eddycharly/json-kyverno/pkg/engine/template"
)

type matchOptions struct {
	wildcard bool
	template template.Template
}

type option func(matchOptions) matchOptions

func ConfigurehWildcard(enabled bool) option {
	return func(o matchOptions) matchOptions {
		o.wildcard = enabled
		return o
	}
}

func WithWildcard() option {
	return ConfigurehWildcard(true)
}

func WithoutWildcard() option {
	return ConfigurehWildcard(false)
}

func WithTemplate(template template.Template) option {
	return func(o matchOptions) matchOptions {
		o.template = template
		return o
	}
}

func newMatchOptions(options ...option) matchOptions {
	var matchOptions matchOptions
	for _, option := range options {
		if option != nil {
			matchOptions = option(matchOptions)
		}
	}
	return matchOptions
}
