package agollo

type Option func(op *operation)

// WithNamespace set namespace for operation
func WithNamespace(namespace string) Option {
	return func(op *operation) {
		op.namespace = namespace
	}
}
