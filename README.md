# ratgo框架 - 快速构建web项目
- 使用gin作为底层框架，在此基础上做进一步封装，同时保留gin原生方法
- 增加简易路由，用以约束开发规范
- 内置扩展包为低耦合设计，均可单独使用
- 使用容器，以长驻内存的方式运行时注入元素，实现复用节省开销
- 内置丰富的工具类，可实现大部分网络功能

### 目录
- [安装与配置](#安装与配置)
- [项目结构](#项目结构)
- [服务配置](#服务配置)
- [快速开始](#快速开始)
- [路由配置](#路由配置)
- [Web应用](#Web应用)
- [数据库](#数据库)
- [Redis](#Redis)
- [Utils工具](#Utils工具)
	- curl
	- 加密算法
	- 文件操作
	- ...

### <a id="安装与配置">安装与配置</a>
#### 1. 安装Go (version 1.10+)，配置使用gitlab私有仓库作为项目的依赖包
```go
1. 获取gitlab的access token
进入Gitlab->Settings->Access Tokens，然后创建一个personal access token，这里权限最好选择只读(read_repository)。有了access token后，我们还需要在git中进行配置，这样才能go get下了私有仓库的包，需要把刚刚的token添加进git的请求头中，操作如下:
	$ git config --global http.extraheader "PRIVATE-TOKEN: YOUR_PRIVATE_TOKEN"

2. 配置git将请求从ssh转换为https:
	$ git config --global url."ssh://git@gitlabj01.vdongchina.com:2022".insteadOf "https://gitlabj01.vdongchina.com"

3. 设置 GOPRIVATE 环境变量:
	$ set GOPRIVATE=gitlabj01.vdongchina.com
```
#### 2. 使用下面命令进行安装 ratgo
	$ go get github.com/vdongchina/ratgo
#### 3. 依赖包安装(已安装可忽略)
	$ go get github.com/gin-gonic/gin
	$ go get github.com/unrolled/secure
	$ go get github.com/go-sql-driver/mysql
	$ go get github.com/jinzhu/gorm
	$ go get github.com/garyburd/redigo/redis
	$ go get gopkg.in/ini.v1
	$ go get github.com/go-touch/mtype
#### 4. 设置系统环境变量
	RATGO_RUNMODE = dev | test | prod
#### 5. 如使用 go mod(自行查找使用方法) 包依赖管理工具，请参考下面命令
##### Windows 下开启 GO111MODULE 并设置 GOPROXY 的命令为：
	$ set GO111MODULE=on
	$ go env -w GOPROXY=https://goproxy.cn,direct
##### MacOS 或者 Linux 下开启 GO111MODULE 并设置 GOPROXY 的命令为：
	$ export GO111MODULE=on
	$ export GOPROXY=https://goproxy.cn

### <a id="项目结构">项目结构</a>
- application
	- modules // Controller存放目录,
    	- demo // 模块名
    - model
        - bussiness // 业务类
        - common // 公共处理函数
        - dao // 抽象数据模型
        - form // 校验数据模型
        - mysql // mysql表字段映射模型
        - redis // redis数据模型
    - router // 注册路由目录
- config
    - dev // 必要配置项(通过系统环境变量选择配置路径)
    	- ratgo.ini // web服务配置
        - database.ini // 数据库配置
        - redis.ini // redis配置
    - test
    - prod
- runtime
    - log // 存储系统错误日志
- main.go // 入口文件

### <a id="路由配置">服务配置</a>
#### 配置项 src/myratgo/config/ratgo.ini
	[main]
	AppName = ratgo // 服务名称
	HTTPAddr = 127.0.0.1:8089 // http地址端口
	HTTPSAddr = 127.0.0.1:10443 // https地址端口
#### 项目中使用配置
	读取方式： (不会直接获取对应值, 而是返回一个*mtype.anyvalue结构体指针, 可实现对应类型转换)
	config ：= ratgo.Config.Get("xx.xx")
	
	调用示例:
	config := ratgo.Config.Get("ratgo.main.AppName").ToString() // 输出 ratgo
	
	备注: ratgo.ini、database.ini、redis.ini等必要配置项,文件和字段名均为ratgo使用,不可修改.
### <a id="快速开始">快速开始</a>
#### 示例说明:
	项目名称：myratgo
	src: GOPATH 下的 src 目录
#### 入口文件: src/myratgo/main.go
```go
package main

import (
	"github.com/vdongchina/ratgo"
	_ "myratgo/application/router"
)

func main() {
	ratgo.RunWebServer()
}
```
##### 编译运行
	$ go run main.go // 编译并运行
	或者:
	$ go build main.go // 编译
	$ ./main.x // 执行编译后的文件

### <a id="路由配置">路由配置</a>
#### 举例: src/myratgo/application/router/demo.go
```go
package router

import (
	"github.com/vdongchina/ratgo"
	"myratgo/application/modules/demo"
	"myratgo/application/modules/demo/sub"
)

func init() {
	// Demo模块
	ratgo.Router.General("demo", ratgo.GeneralMap{
		"index": &demo.Index{}, // 对应控制器: application/modules/demo下的Index结构体
		"list":  &demo.List{}, // 对应控制器: application/modules/demo下的List结构体
		"sub/index": &sub.Index{}, // 对应控制器: application/modules/demo/sub下的Index结构体
	})

	// 静态文件
	ratgo.Router.SetStatic("/apidoc", "./apidoc") // 访问: http://127.0.0.1:port/apidoc
	ratgo.Router.SetStatic("/assets", "./assets") // 访问: http://127.0.0.1:port/assets
}
```
```go
访问:
http://127.0.0.1:port/xx/xx
http://127.0.0.1:port/xx/xx/xx

示例:
http://127.0.0.1:port/demo/index
http://127.0.0.1:port/demo/list
http://127.0.0.1:port/demo/sub/index
```
```go
Note: 该路由模式为简易模式, 对应gin: any("/xx/:param1",func())、any("/xx/:param1/:param2",func())两种模式
,即两段和三段路由无拦截访问。
```

### <a id="Web应用">Web应用</a>
#### 基类Controller的代码示例: ratgo/Controller.go
```go
package ratgo

import "github.com/gin-gonic/gin"

// 控制器
type Controller struct {
	context *gin.Context
	result  *Result
}

// 控制器初始化方法 -- 此方法为系统调用
func (c *Controller) Init(context *gin.Context, result *Result) {
	c.context = context
	c.result = result
}

// 获取 gin 的上下文Context -- 对应一个用户请求,用法可参考gin框架
func (c *Controller) Context() *gin.Context {
	return c.context
}

// 简易模式 - 前置方法 -- 在简易路由模式下,系统会首先执行该方法
func (c *Controller) BeforeExec() {

}

// 简易模式 - 执行方法 -- 在简易路由模式下,系统执行完 BeforeExec() 后会执行该方法
func (c *Controller) Exec() {

}

// 结果方法 -- 控制响应(辅助用法,阻断执行直接响应异常)
func (c *Controller) Result() *Result {
	return c.result
}
```
	Note:
	BeforeExec()方法可用于鉴权、登录验证等公共功能;
	Result()方法获取到一个*ratgo.result可修改状态码、响应类型等, 系统调用执行方法时会根据状态码进行逻辑处理.
#### Controller.result的示例代码: ratgo/Result.go
```go
package ratgo

/**************************************** 数据类型 - 结构体Result ****************************************/
// 定义常量
const (
	RespString = "String"
	RespJson   = "Json"
	RespHtml   = "Html"
)

// 响应结果
type Result struct {
	Status int         // 状态码: [200:OK] [400:Bad Request] [500:Internal Server Error] [900:逻辑异常(状态码200)]
	Type   string      // 响应类型: String、Json、Html 默认Json
	Msg    string      // 消息提示
	Data   interface{} // 响应数据
}

// 实例化 Result
func NewResult() *Result {
	return &Result{
		Status: 200,
		Type:   RespJson,
		Msg:    "",
		Data:   "",
	}
}

// 设置Code
func (r *Result) SetStatus(status int) {
	r.Status = status
}

// 设置Code
func (r *Result) SetType(t string) {
	r.Type = t
}

// 设置Msg
func (r *Result) SetMsg(msg string) {
	r.Msg = msg
}

// 设置Data
func (r *Result) SetData(data interface{}) {
	r.Data = data
}
```
#### 项目Controller的代码示例: myratgo/application/modules/demo/index.go
```go
package demo

import (
	"github.com/vdongchina/ratgo"
)

type Index struct {
	ratgo.Controller
}

/**
 * @api {post} /demo/index  xxx接口
 * @apiDescription 1.0
 * @apiGroup api
 * @apiVersion 1.0.0
 * @apiParam {string} xxx 参数
 * @apiSuccessExample {json} 正确返回值:
 * {"code":200,"msg":"ok","time":1560245913,"data":""}
 */
func (this *Index) Exec() {
	this.Context().String(200, "执行Controller: demo.Index")
}
```
	Note：访问 http://127.0.0.1:port/demo/index // 浏览器输出: 执行Controller: demo.Index
### <a id="数据库">数据库</a>
#### 配置项 myratgo/config/dev/database.ini
	[plus_center] // 配置分组,必填
	; 主库
	master.driverName = mysql // 驱动名称
	master.dataSourceName = root:root@tcp(127.0.0.1:3306)/dbName?charset=utf8 // 连接参数
	master.maxIdleConn = 100 // 空闲连接数
	master.maxOpenConn = 100 // 最大连接数
	master.maxLifetime = 1800 // 连接超时
	
	; 从库
	slave.driverName = mysql
	slave.dataSourceName = root:root@tcp(127.0.0.1:3306)/dbName?charset=utf8
	slave.maxIdleConn = 100
	slave.maxOpenConn = 100
	slave.maxLifetime = 1800
#### 默认使用Gorm(用法请参考 http://gorm.book.jasperxu.com)
```go
func (ds *DbStorage) Db(key string) *gorm.DB // 返回 *gorm.DB

示例:
db := extend.Gorm.Db("plus_center.master") // 对应配置项设置
db.xxx() // gorm的用法
```
### <a id="Redis">Redis</a>

#### 配置项 xxx/config/dev/redis.ini
	[plus_center] // // 配置分组,必填
	master.host = 127.0.0.1:6379 // 主机端口
	master.password = "" // 密码
	master.db = 10 // 库标
	master.MaxIdle = 16 // 空闲连接数
	master.MaxActive = 32 // 最大连接数 
	master.IdleTimeout = 120 // 超时时间

#### RedisModel的示例
```go
package redis

type TestModel struct {}

// Redis库标识
func (b *Base) Identify() string {
	return "plus_center.master"
}
```
```go
(this *Users) Identify() string // 设置redis连接参数,对应Redis配置的key链关系
```
#### Redis的使用示例:
##### 传入一个model获取 Redis Dao 实例
```go
RedisModel(model interface{}) *RedisDao

示例:
redisDao := RedisModel(&TestModel{})
```
##### 获取连接池对象,开发者可通过此返回值
```go
(rd *RedisDao) Pool() *redis.Pool

示例:
pool := RedisModel(&TestModel{}).Pool()
```
##### 执行redis命令,返回\*base.AnyValue,可进行类型转换. 参数name:命令名称 args:该命令对应的参数
```go
(rd *RedisDao) Command(name string, args ...interface{}) *base.AnyValue

示例:
RedisModel(&TestModel{}).Command("SET","username","admin")
RedisModel(&TestModel{}).Command("HSET","user","username","admin")
ret := RedisModel(&TestModel{}).Command("GET","username")
ret.ToError() // 可获取错误信息,如果返回nil,则说明无错误发生
ret.ToAffectedRows() // 返回受影响行数
```
### <a id="Utils工具">Utils工具</a>
#### Curl使用
##### GET请求
##### 获取一个 *curl.GetCaller
```go
func Get() *GetCaller
```
##### 设置header信息
```go
func (gc *GetCaller) Header(header map[string]string)
```
##### 发送一个GET请求
```go
func (gc *GetCaller) Call(url string, args ...map[string]interface{}) *multitype.AnyValue
```
##### GET请求示例:
```go
get := curl.Get()

// 如需设置header, 则使用此方法
get.Header(map[string]string{"Authorization":"Basic MTAwMToxMjM0NTY="})

// 发送请求
get.Call("http://www.baidu.com",map[string]interface{}{
	"user_id":  1,
	"username": "admin",
	"password": "123456",
})
```
##### POST请求
##### 获取一个 *curl.PostCaller
```go
func Post() *PostCaller
```
##### 设置header信息
```go
func (pc *PostCaller) Header(header map[string]string)
```
##### 发送一个POST请求
```go
func (pc *PostCaller) Call(url string, args ...interface{}) *multitype.AnyValue
```
##### POST请求示例:
```go
post := curl.Post()

// 如需设置header, 则使用此方法
post.Header(map[string]string{"Authorization":"Basic MTAwMToxMjM0NTY="})

// 设置 json 请求 header, 默认 header {"Content-Type": "application/x-www-form-urlencoded"}
post.Header(map[string]string{"Content-Type": "application/json"})

// 发送请求(key-value形式)
post.Call("http://www.baidu.com",map[string]interface{}{
	"user_id":  1,
	"username": "admin",
	"password": "123456",
})

// 发送请求(json串形式)
post.Call("http://www.baidu.com",`{"user_id":1,"username":"admin","password":"123456"}`)
```












