package types

import (
	"strconv"
	"strings"
)

// Set value.
func Set(args string, value interface{}, src interface{}) *LinkedList {
	switch src.(type) {
	case AnyMap:
		src = map[string]interface{}(src.(AnyMap))
	case AnySlice:
		src = []interface{}(src.(AnySlice))
	default:
		return nil
	}
	// Handle.
	argsGroup := strings.Split(args, ".")
	length := len(argsGroup) - 1
	breakFor := false
	link := &LinkedList{
		Value:   src,
		NextKey: argsGroup[0],
		Next:    nil,
	}
	for i := 0; i <= length; i++ {
		switch GetType(src) {
		case TAnyMap:
			if v, ok := src.(map[string]interface{})[argsGroup[i]]; ok {
				src = v
			} else {
				breakFor = true
			}
		case TStringMap:
			if v, ok := src.(map[string]string)[argsGroup[i]]; ok {
				src = v
			} else {
				breakFor = true
			}
		case TAnySlice:
			if intKey, err := strconv.Atoi(argsGroup[i]); err != nil {
				breakFor = true
			} else if intKey < len(src.([]interface{})) {
				src = src.([]interface{})[intKey]
			} else {
				breakFor = true
			}
		case TStringSlice:
			if intKey, err := strconv.Atoi(argsGroup[i]); err != nil {
				breakFor = true
			} else if intKey < len(src.([]string)) {
				src = src.([]string)[intKey]
			} else {
				breakFor = true
			}
		default:
			breakFor = true
		}
		if breakFor && i < length {
			break
		}
		if i == length {
			AddLinkedList(link, &LinkedList{
				Value:   value,
				NextKey: "",
				Next:    nil,
			})
		} else {
			AddLinkedList(link, &LinkedList{
				Value:   src,
				NextKey: argsGroup[i+1],
				Next:    nil,
			})
		}
	}
	return link
}

// Get value.
func Get(src interface{}, args ...string) *AnyValue {
	switch src.(type) {
	case AnyMap:
		src = map[string]interface{}(src.(AnyMap))
	case AnySlice:
		src = []interface{}(src.(AnySlice))
	default:
		return nil
	}
	if args == nil || len(args) == 0 {
		return Eval(src)
	}

	// Data handle.
	argsGroup := strings.Split(args[0], ".")
	for _, key := range argsGroup {
		switch GetType(src) {
		case TStringMap:
			if v, ok := src.(map[string]string)[key]; ok {
				src = v
			} else {
				src = nil
				break
			}
		case TAnyMap:
			if v, ok := src.(map[string]interface{})[key]; ok {
				src = v
			} else {
				src = nil
				break
			}
		case TStringSlice:
			if intKey, err := strconv.Atoi(key); err != nil {
				src = nil
				break
			} else if intKey < len(src.([]string)) {
				src = src.([]string)[intKey]
			} else {
				src = nil
				break
			}
		case TAnySlice:
			if intKey, err := strconv.Atoi(key); err != nil {
				src = nil
				break
			} else if intKey < len(src.([]interface{})) {
				src = src.([]interface{})[intKey]
			} else {
				src = nil
				break
			}
		default:
			break
		}
	}
	return Eval(src)
}