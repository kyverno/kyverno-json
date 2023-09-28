package match

type matchOptions struct {
	wildcard bool
}

type option func(matchOptions) matchOptions

func WithWildcard(enabled bool) option {
	return func(o matchOptions) matchOptions {
		o.wildcard = enabled
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
