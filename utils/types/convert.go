package types

import (
	"encoding/json"
	"strconv"
)

// Define const.
const (
	TInt            = "Int"            // 类型: int
	TInt64          = "Int64"          // 类型: int
	TByte           = "Byte"           // 类型: byte
	TRune           = "Rune"           // 类型: rune
	TFloat64        = "Float64"        // 类型: float64
	TString         = "String"         // 类型: string
	TBool           = "Bool"           // 类型: bool
	TIntSlice       = "IntSlice"       // 类型: []int
	TByteSlice      = "ByteSlice"      // 类型: []byte
	TStringSlice    = "StringSlice "   // 类型: []string
	TAnySlice       = "AnySlice "      // 类型: []interface{}
	TStringMapSlice = "StringMapSlice" // 类型: []map[string]string
	TAnyMapSlice    = "AnyMapSlice"    // 类型: []map[string]interface{}
	TStringMap      = "StringMap"      // 类型: map[string]string
	TAnyMap         = "AnyMap"         // 类型: map[string]interface{}
	TError          = "Error"          // 类型: error
	TNil            = "Nil"            // 类型: nil
	TUnknown        = "unknown"        // 类型: unknown
)

// Get value type.
func GetType(value interface{}) string {
	switch value.(type) {
	case int:
		return TInt
	case int64:
		return TInt64
	case byte:
		return TByte
	case rune:
		return TRune
	case float64:
		return TFloat64
	case string:
		return TString
	case bool:
		return TBool
	case []int:
		return TIntSlice
	case []byte:
		return TByteSlice
	case []string:
		return TStringSlice
	case []interface{}:
		return TAnySlice
	case []map[string]string:
		return TStringMapSlice
	case []map[string]interface{}:
		return TAnyMapSlice
	case map[string]string:
		return TStringMap
	case map[string]interface{}, AnyMap:
		return TAnyMap
	case error:
		return TError
	case nil:
		return TNil
	}
	return TUnknown
}

// Convert src to target type.
func ToType(src interface{}, dstType string) interface{} {
	switch GetType(src) {
	case TInt:
		switch dstType { // int 转 int
		case TInt:
			return src
		case TByte: // int 转 byte
			return byte(src.(int))
		case TString: // int 转 string
			return strconv.Itoa(src.(int))
		case TBool: // int 转 bool
			if src.(int) > 0 {
				return true
			}
			return false
		}
	case TByte:
		switch dstType { // byte 转 int
		case TInt:
			return int(src.(byte))
		case TByte: // byte 转 byte
			return src
		case TString: // byte 转 string
			return string([]byte{src.(byte)})
		case TBool: // byte 转 bool
			if src.(byte) > 0 {
				return true
			}
			return false
		}
	case TFloat64:
		switch dstType { // byte 转 int
		case TInt:
			return int(src.(float64))
		case TByte: // byte 转 byte
			return byte(src.(float64))
		case TString: // byte 转 string
			v := int(src.(float64))
			return ToType(v, TString)
		case TBool: // byte 转 bool
			if src.(float64) > 0 {
				return true
			}
			return false
		}
	case TString:
		switch dstType {
		case TInt: // string 转 int
			if v, err := strconv.Atoi(src.(string)); err != nil {
				return 0
			} else {
				return v
			}
		case TInt64: // string 转 int64
			if v, err := strconv.ParseInt(src.(string), 10, 64); err == nil {
				return v
			}
		case TByte: // string 转 byte
			return src
		case TString: // string 转 string
			return src
		case TBool: // string 转 bool
			if src.(string) == "true" {
				return true
			}
			return false
		case TByteSlice: // string 转 []Byte
			return []byte(src.(string))
		}
	case TBool:
		switch dstType {
		case TInt:
			if src == true {
				return 1
			}
			return 0
		case TByte:
			if src == true {
				return byte(1)
			}
			return byte(0)
		case TString:
			if src == true {
				return "true"
			}
			return "false"
		case TBool:
			return src
		}
	case TIntSlice:
		switch dstType { // []int 转 []int
		case TIntSlice:
			return src
		case TStringSlice: // []int 转 []string
			var v []string
			for _, value := range src.([]int) {
				v = append(v, ToType(value, TString).(string))
			}
			return v
		}
	case TByteSlice:
		switch dstType {
		case TString: // []byte 转 string
			return string(src.([]byte))
		case TIntSlice: // []byte 转 []int
			var v []int
			for _, value := range src.([]byte) {
				v = append(v, int(value))
			}
			return v
		case TStringSlice:
			var v []string
			for _, value := range src.([]byte) {
				v = append(v, ToType(value, TString).(string))
			}
			return v
		case TAnyMap:
			var v map[string]interface{}
			if err := json.Unmarshal(src.([]byte), &v); err == nil {
				return v
			}
		}
	case TStringSlice:
		switch dstType {
		case TIntSlice: // []string 转 []int
			var v []int
			for _, value := range src.([]string) {
				v = append(v, ToType(value, TInt).(int))
			}
			return v
		case TStringSlice: // []string 转 []string
			return src
		case TAnySlice: // []string 转 []string
			var v []interface{}
			for _, value := range src.([]string) {
				v = append(v, value)
			}
			return v
		}
	case TAnySlice:
		switch dstType {
		case TIntSlice: // []interface{} 转 []int
			var v []int

			// 类型断言
			var srcDeclare []interface{}
			if declare, ok := src.(AnySlice); ok {
				srcDeclare = declare
			} else {
				srcDeclare = src.([]interface{})
			}
			for _, value := range srcDeclare {
				if dstValue := ToType(value, TInt); dstValue != nil {
					v = append(v, dstValue.(int))
				} else {
					v = append(v, 0)
				}
			}
			return v
		case TStringSlice: // []interface{} 转 []string
			var v []string
			// 类型断言
			var srcDeclare []interface{}
			if declare, ok := src.(AnySlice); ok {
				srcDeclare = declare
			} else {
				srcDeclare = src.([]interface{})
			}
			for _, value := range srcDeclare {
				if dstValue := ToType(value, TString); dstValue != nil {
					v = append(v, dstValue.(string))
				} else {
					v = append(v, "")
				}
			}
			return v
		case TAnySlice:
			var v []interface{}
			if declare, ok := src.(AnySlice); ok {
				v = declare
			} else {
				v = src.([]interface{})
			}
			return v
		}
	case TStringMap:
		switch dstType {
		case TStringSlice: // map[string]string 转 []string
			var v []string
			for _, value := range src.(map[string]string) {
				v = append(v, value)
			}
			return v
		case TAnySlice: // map[string]string 转 []interface{}
			var v []interface{}
			for _, value := range src.(map[string]string) {
				v = append(v, value)
			}
			return v
		case TStringMap: // // map[string]string 转 map[string]string
			return src
		case TAnyMap: // map[string]string 转 map[string]interface{}
			v := map[string]interface{}{}
			for key, value := range src.(map[string]string) {
				v[key] = value
			}
			return v
		}
	case TAnyMap:
		switch dstType {
		case TStringMap: // map[string]interface{} 转 map[string]string
			v := map[string]string{}

			// 类型断言
			var srcDeclare map[string]interface{}
			if declare, ok := src.(AnyMap); ok {
				srcDeclare = declare
			} else {
				srcDeclare = src.(map[string]interface{})
			}
			for key, value := range srcDeclare {
				if dstValue := ToType(value, TString); dstValue != nil {
					v[key] = dstValue.(string)
				} else {
					v[key] = ""
				}
			}
			return v
		case TAnyMap: // map[string]interface{} 转 map[string]interface{}
			var v map[string]interface{}
			if declare, ok := src.(AnyMap); ok {
				v = declare
			} else {
				v = src.(map[string]interface{})
			}
			return v
		}
	}
	return nil
}
