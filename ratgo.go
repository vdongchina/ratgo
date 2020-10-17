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

// 应用容器
var AppStorage map[string]interface{}

// 初始化
func init() {
	AppStorage = map[string]interface{}{}
}

// 运行web服务
func RunWebServer() {
	webServer := NewWebServer()         // 获取WebServer指针
	AppStorage["WebServer"] = webServer // 存储WebServer
	webServer.Init()                    // WebServer初始化
	webServer.Run()                     // 运行WebServer
}

// 获取 *WebServer
func GetWebServer() *WebServer {
	if webServer, ok := AppStorage["WebServer"]; ok {
		return webServer.(*WebServer)
	}
	return nil
}

// 运行cmd服务
func RunCmd() {
	// ...
}

// 运行websocket服务
func RunWebsocket() {
	// ...
}
