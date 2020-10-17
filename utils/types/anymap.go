package types

// Built-in type: map[string]interface{}.
type AnyMap map[string]interface{}

// Set value.
func (am *AnyMap) Set(args string, value interface{}) {
	link := Set(args, value, *am)
	*am = link.Reassignment().(map[string]interface{})
}

// Get value.
func (am *AnyMap) Get(args ...string) *AnyValue {
	return Get(*am, args...)
}