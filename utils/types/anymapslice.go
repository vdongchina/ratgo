package types

// Built-in type: []map[string]interface{}.
type AnyMapSlice []map[string]interface{}

// Set value.
func (ams *AnyMapSlice) Set(args string, value interface{}) {
	link := Set(args, value, *ams)
	*ams = link.Reassignment().([]map[string]interface{})
}

// Get value.
func (ams *AnyMapSlice) Get(args ...string) *AnyValue {
	return Get(*ams, args...)
}