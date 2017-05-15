package middleware

// A ContextID represents context key for middleware.
type ContextID struct {
	name string
}

// String returns formatted context key identifier.
func (k *ContextID) String() string {
	return "ottoman/middleware context: " + k.name
}

// ContextKey constructs context key using name supplied.
func ContextKey(name string) *ContextID {
	return &ContextID{name: name}
}
