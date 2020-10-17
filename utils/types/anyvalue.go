package types

// Any type data storage, use it can convert to target type.
type AnyValue struct {
	value interface{}
}

// Accept an any type value.
func Eval(value interface{}) *AnyValue {
	return &AnyValue{value: value}
}

// Get type of origin value.
func (av *AnyValue) Type() string {
	return GetType(av.value)
}

// Get origin value.
func (av *AnyValue) Value() interface{} {
	return av.value
}

// Convert value to type of error.
func (av *AnyValue) ToError() error {
	switch GetType(av.value) {
	case TError:
		return av.value.(error)
	}
	return nil
}

// Convert value to type of int.
func (av *AnyValue) ToInt() int {
	if dstValue := ToType(av.value, TInt); dstValue != nil {
		return dstValue.(int)
	}
	return 0
}

// Convert value to type of int.
func (av *AnyValue) ToInt64() int64 {
	if dstValue := ToType(av.value, TInt64); dstValue != nil {
		return dstValue.(int64)
	}
	return 0
}

// Convert value to type of byte.
func (av *AnyValue) ToByte() byte {
	if dstValue := ToType(av.value, TByte); dstValue != nil {
		return dstValue.(byte)
	}
	return byte(0)
}

// Convert value to type of string.
func (av *AnyValue) ToString() string {
	if dstValue := ToType(av.value, TString); dstValue != nil {
		return dstValue.(string)
	}
	return ""
}

// Convert value to type of bool.
func (av *AnyValue) ToBool() bool {
	if dstValue := ToType(av.value, TBool); dstValue != nil {
		return dstValue.(bool)
	}
	return false
}

// Convert value to type of []int.
func (av *AnyValue) ToIntSlice() []int {
	value := make([]int, 0)
	v := ToType(av.value, TIntSlice)
	if v != nil {
		value = v.([]int)
	}
	return value
}

// Convert value to type of []byte.
func (av *AnyValue) ToByteSlice() []byte {
	value := make([]byte, 0)
	v := ToType(av.value, TByteSlice)
	if v != nil {
		value = v.([]byte)
	}
	return value
}

// Convert value to type of []string.
func (av *AnyValue) ToStringSlice() []string {
	value := make([]string, 0)
	v := ToType(av.value, TStringSlice)
	if v != nil {
		value = v.([]string)
	}
	return value
}

// Convert value to type of map[string]string.
func (av *AnyValue) ToStringMap() map[string]string {
	value := map[string]string{}
	v := ToType(av.value, TStringMap)
	if v != nil {
		value = v.(map[string]string)
	}
	return value
}

// Convert value to type of map[string]interface{}.
func (av *AnyValue) ToAnyMap() map[string]interface{} {
	value := map[string]interface{}{}
	v := ToType(av.value, TAnyMap)
	if v != nil {
		value = v.(map[string]interface{})
	}
	return value
}

// Convert value to type of []map[string]string .
func (av *AnyValue) ToStringMapSlice() []map[string]string {
	value := make([]map[string]string, 0)
	switch GetType(av.value) {
	case TStringMapSlice:
		return av.value.([]map[string]string)
	case TAnyMapSlice:
		for k, v := range av.value.([]map[string]interface{}) {
			value[k] = ToType(v, TStringMap).(map[string]string)
		}
	}
	return value
}

// Convert value to type of []interface{}.
func (av *AnyValue) ToAnySlice() []interface{} {
	value := make([]interface{}, 0)
	v := ToType(av.value, TAnySlice)
	if v != nil {
		value = v.([]interface{})
	}
	return value
}

// Convert value to type of []map[string]interface{} .
func (av *AnyValue) ToAnyMapSlice() []map[string]interface{} {
	value := make([]map[string]interface{}, 0)
	switch GetType(av.value) {
	case TStringMapSlice:
		for k, v := range av.value.([]map[string]string) {
			value[k] = ToType(v, TAnyMap).(map[string]interface{})
		}
	case TAnyMapSlice:
		value = av.value.([]map[string]interface{})
	}
	return value
}