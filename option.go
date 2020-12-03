package agollo

type OpOption func(op *operation)

// WithNamespace set namespace for operation
func WithNamespace(namespace string) OpOption {
	return func(op *operation) {
		op.namespace = namespace
	}
}

type ClientOption func(*client)

// WithLogger set client logger
func WithLogger(logger Logger) ClientOption {
	return func(c *client) {
		c.logger = logger
	}
}

func SkipLocalCache() ClientOption {
	return func(c *client) {
		c.skipLocalCache = true
	}
}