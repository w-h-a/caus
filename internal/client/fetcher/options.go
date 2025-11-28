package fetcher

import "context"

type Option func(*Options)

type Options struct {
	Location string
	ApiKey   string
	AppKey   string
	Context  context.Context
}

func WithLocation(loc string) Option {
	return func(o *Options) {
		o.Location = loc
	}
}

func WithApiKey(key string) Option {
	return func(o *Options) {
		o.ApiKey = key
	}
}

func WithAppKey(key string) Option {
	return func(o *Options) {
		o.AppKey = key
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
