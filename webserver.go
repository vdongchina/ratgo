// Copyright 2020 ratgo Author. All Rights Reserved.
// Licensed under the Apache License, Version 1.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ratgo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vdongchina/ratgo/ext"
	"github.com/vdongchina/ratgo/extend"
	"github.com/vdongchina/ratgo/extend/cache"
	"net/http"
	"reflect"
)

// Error handle.
type ErrorHandle func()

// Web server.
type WebServer struct {
	gin *gin.Engine
}

// 获取web实例
func NewWebServer() *WebServer {
	return &WebServer{
		gin: gin.Default(),
	}
}

// 初始化
func (ws *WebServer) Init() {
	Config.Init() // 配置初始化
	// 初始化数据库
	if Config.InitDb == true {
		extend.Gorm.Init(Config.Get("database").ToAnyMap())
	}
	// 初始化redis
	if Config.InitRedis == true {
		cache.Redis.Init(Config.Get("redis").ToAnyMap())
		ext.Redis.Init(Config.Get("redis").ToAnyMap())
	}
	ws.gin.Use(LoggerInit()) // 日志注入
	// 执行用户挂载函数
	if len(UserFuncArray) > 0 {
		for _, function := range UserFuncArray {
			err := function()
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

// 运行server
func (ws *WebServer) Run() {
	_ = ws.registerMiddleWare()     // 注册中间件
	_ = ws.registerRouter()         // 注册路由
	_ = ws.registerStatic()         // 注册静态文件
	_ = ws.gin.Run(Config.HTTPAddr) // 运行gin
}

// 获取原生gin
func (ws *WebServer) Gin() *gin.Engine {
	return ws.gin
}

// 注册全局中间件
func (ws *WebServer) registerMiddleWare() error {
	if middleWare := MiddleWare.GetGlobal(); len(middleWare) > 0 {
		ws.gin.Use(middleWare...)
	}
	return nil
}

// 注册路由
func (ws *WebServer) registerRouter() error {
	if Router.Mode == "General" { // 简易模式
		if generalMap := Router.GetGeneralPath(); len(generalMap) > 0 {
			for _, value := range generalMap {
				for _, v := range value {
					ws.gin.HEAD(v, ws.generalHandle)    // 注册 HEAD handle
					ws.gin.GET(v, ws.generalHandle)     // 注册 GET handle
					ws.gin.POST(v, ws.generalHandle)    // 注册 POST handle
					ws.gin.OPTIONS(v, ws.generalHandle) // 注册 OPTIONS handle
				}
			}
		}
	} else if Router.Mode == "Restful" { // Restful模式
		// ...
	}
	return nil
}

// 注册静态文件
func (ws *WebServer) registerStatic() error {
	// 静态路径
	if staticMap := Router.GetStatic(); len(staticMap) > 0 {
		for key, value := range staticMap {
			// ws.gin.Static(key, value)
			ws.gin.StaticFS(key, http.Dir(value))
		}
	}

	// 静态文件
	if staticFileMap := Router.GetStaticFile(); len(staticFileMap) > 0 {
		for key, value := range staticFileMap {
			ws.gin.StaticFile(key, value)
		}
	}

	// html模板
	if htmlGlob := Router.GetHTMLGlob(); len(htmlGlob) > 0 {
		for _, value := range htmlGlob {
			ws.gin.LoadHTMLGlob(value)
		}
	}
	return nil
}

// 简易路由 handle
func (ws *WebServer) generalHandle(context *gin.Context) {
	defer ws.errorCatch(context)
	// 开始时间
	// startTime := time.Now()

	// 获取
	path := context.Request.URL.Path
	generalCtrl := Router.GetGeneral(path)
	if generalCtrl == nil {
		panic(fmt.Sprintf("get controller failed by path '%s'", path))
	}

	// 克隆结构体
	ctrlType := reflect.TypeOf(generalCtrl).Elem()
	controller, ok := reflect.New(ctrlType).Interface().(ControllerInterface)
	if !ok {
		panic("controller is not ControllerInterface")
	} else {
		controller.Init(context, NewResult())
	}

	// 执行 BeforeExec()
	controller.BeforeExec()
	if result := controller.Result(); result.Status != 200 {
		ws.Response(context, result)
		return
	}

	// 执行 Exec()
	controller.Exec()
	if result := controller.Result(); result.Status != 200 {
		ws.Response(context, result)
		return
	}
}

// 异常响应
func (ws *WebServer) Response(context *gin.Context, result *Result) {
	if result.Status == 900 {
		result.Status = 200
	}
	switch result.Type {
	case "String":
		context.String(result.Status, result.Msg)
	case "Json":
		context.JSON(result.Status, result.Data)
	case "Html":
		context.HTML(result.Status, result.Msg, result.Data)
	}
}

// 异常捕获
func (ws *WebServer) errorCatch(context *gin.Context) {
	if r := recover(); r != nil {
		context.String(200, Error(r).Error())
	}
}
