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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vdongchina/ratgo/utils/encrypt"
	"github.com/vdongchina/ratgo/utils/vdlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 注册日志中间件
func RegisterLogMiddleWare() error {
	MiddleWare.SetGlobal(func(context *gin.Context) {
		requestTime := time.Now()                                                       // 请求时间
		requestUri := context.Request.RequestURI                                        // 请求路径
		requestId := encrypt.Md5(requestUri, fmt.Sprintf("%d", requestTime.UnixNano())) // 生成 requestId

		// 构造入参数据
		request := map[string]interface{}{
			"time":     requestTime,
			"method":   context.Request.Method,
			"query":    requestUri,
			"header":   context.Request.Header,
			"body":     map[string]interface{}{},
			"clientIP": context.ClientIP(),
		}
		// 处理 body 数据
		if context.Request.Method == http.MethodPost {
			rawData, _ := context.GetRawData()
			switch context.ContentType() {
			case "application/x-www-form-urlencoded":
				urlValue, _ := url.ParseQuery(string(rawData))
				request["body"] = urlValue
			case "application/json":
				var body map[string]interface{}
				_ = json.Unmarshal(rawData, &body)
				request["body"] = body
			case "application/xml", "text/xml":
				request["body"] = string(rawData)
			}
			// body回放
			context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
		}

		// 错误接收
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[log-middleWare] a fatal error occurred during the runtime, error: %v\n", r)
			}
		}()

		// 记录请求日志
		logger := vdlog.Clone().SetLogId(requestId) // 克隆 logger并设置 requestId
		logger.Info(request)                        // 记录请求数据
		context.Set("logger", logger)               // 存储日志
		context.Set("requestId", requestId)         // 设置 requestId

		// 处理请求
		context.Next()

		// 响应数据
		response, _ := context.Get("response") // 响应数据
		logger.Info(response)                  // 记录响应数据
	})
	return nil
}
