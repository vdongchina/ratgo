## types ##
#### 功能支持
- [获取类型](#获取类型)
- [类型转换](#类型转换)
- 内置魔法类型, 对复杂的数据结构字段值实现增删改查
	- [types.AnyMap](#types.AnyMap)
	- [types.AnySlice](#types.AnySlice)

#### 支持类型
```go
// Define const.
const (
	TInt            = "Int"            // 类型: int
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
```
##### <a id="获取类型">获取类型</a> (对应上面定义的类型字符串)
```go
func GetType(value interface{}) string
```
##### 转换成对应类型 参数 src:被转换的类型值 dstType:将要转换的类型(参考上面定义类型)
```go
func ToType(src interface{}, dstType string) interface{}

// 示例:
types.ToType("10", types.TInt) // 字符串10转换成int
```
```go
Note: 该方法为底层转换方法, 转化后的值还是以 interface{} 形式返回. 如要实现真正转换, 请参考下面 *types.AnyValue
```
#### <a id="类型转换">类型转换</a> 使用 *types.AnyValue（对于不确定类型interfa{}比较适用)
##### 1. 获取 *types.AnyValue. 参数 value: interface{} (即被转换的值, 可传任意类型值)
```go
func Eval(value interface{}) *AnyValue

示例:
any := types.Eval(map[string]interface{}{})
```
##### 2. 返回原值类型(对应上面支持类型, 例如:Int 字符串)
```go
func (av *AnyValue) Type() string
```
##### 3. 返回原值
```go
func (av *AnyValue) Value() interface{}
```
##### 4. 返回错误信息
```go
func (av *AnyValue) ToError() error
```
##### 5. 转成int类型
```go
func (av *AnyValue) ToInt() int
```
##### 6. 转成byte类型
```go
func (av *AnyValue) ToByte() byte
```
##### 7. 转成string类型
```go
func (av *AnyValue) ToString() string
```
##### 8. 转成bool类型
```go
func (av *AnyValue) ToBool() bool
```
##### 9. 转成map[string]string类型
```go
func (av *AnyValue) ToStringMap() map[string]string
```
##### 10. 更多方法...
```go
Note: *types.AnyValue.Method()
```
##### 完整调用示例:
```go
intValue := types.Eval("10").ToInt()
stringValue := types.Eval(10).ToString()
stringMap := types.Eval("master": map[string]interface{}{
	"host": "127.0.0.1",
	"port": "3306",
}).ToStringMap()
```
#### <a id="types.AnyMap">内置类型 types.AnyMap</a>
##### 1. 使用方法
```go
第一种方法:
anyMap := types.AnyMap{
	"xxx":"xxx",
}

第二种方法:
anyMap := types.AnyMap(map[string]interface{}{
	"xxx": "xxx",
})
```
```go
anyMap的数据结构如下(使用json演示):
{
	"database": {
		"master": {
			"host": "127.0.0.1",
			"port": "3306"
		},
		"salve": {
			"host": "127.0.0.1",
			"port": "3307"
		}
	}
}
```
##### 2. 查询字段值 参数 args: 任意 . 拼接字符串 (为空则返回 value = 自身值的 *AnyValue)
```go
func (am *AnyMap) Get(args ...string) *AnyValue
```
```go
调用示例:
v := anyMap.Get("database.master.port")
fmt.Println(v.ToString()) // 输出 3306

v := anyMap.Get()
fmt.Println(v.ToAnyMap()) // 输出 map[database:map[master:map[host:127.0.0.1 port:3306] salve:map[host:127.0.0.1 port:3307]]]
```
##### 3. 插入或更新字段值 参数 args: 任意 . 拼接字符串 value: 要插入或更新的值
```go
func (am *AnyMap) Set(args string, value interface{})
```
```go
调用示例:
anyMap.Set("database.master.port",3389)
v := anyMap.Get("database.master.port")
fmt.Println(v.ToString()) // 输出 3389

v := anyMap.Get()
fmt.Println(v.ToAnyMap()) // 输出 map[database:map[master:map[host:127.0.0.1 port:3389] salve:map[host:127.0.0.1 port:3307]]]
```
#### <a id="types.AnySlice">内置类型 types.AnySlice</a>
```go
与 types.AnyMap使用方法一致
```
```go
anySlice的数据结构如下(使用json演示):
[
	{
		"master": {
			"host": "127.0.0.1",
			"port": "3306"
		},
		"salve": {
			"host": "127.0.0.1",
			"port": "3307"
		}
	}
]
```
```go
Note: key值索引直接使用字符串数值, 例:anySlice.Get("0.master.port")
```
#### 内置类型 更多持续更新中...