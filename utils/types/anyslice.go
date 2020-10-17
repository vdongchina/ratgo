package types

// Built-in type: []interface{}.
type AnySlice []interface{}

// Set value.
func (as *AnySlice) Set(args string, value interface{}) {
	link := Set(args, value, *as)
	*as = link.Reassignment().([]interface{})
}

// Get value.
func (as *AnySlice) Get(args ...string) *AnyValue {
	return Get(*as, args...)
}
