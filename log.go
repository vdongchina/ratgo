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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vdongchina/ratgo/utils/vdlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

/**
分析:
go原生的:仅提供格式化字符串输出到logger
gin框架带的: 仅提供文本格式生成到文本
logrus开源的:  提供 文本、json格式存储,syslog存储

日志模式:   debug支持list;release
日志存储:   本地  syslog redis es mongo
日志格式:   文本  json

日志分割:当每天文件足够大到一个文件放不下的时候，分成多文件存放日志，后期再添加

1.日志配置
2.日志方法注入
3.日志存储切割
4.文件操作
*/

// 日志初始化
func LoggerInit() gin.HandlerFunc {
	switch Config.Storage {
	case "syslog":
		return LoggerToSyslog()
	case "es":
		return LoggerToES()
	default:
		return LoggerToFile()
	}
}

// 注册日志中间件
func RegisterLogMiddleWare() error {
	MiddleWare.SetGlobal(func(context *gin.Context) {
		requestTime := time.Now()                                               // 请求时间
		requestUri := context.Request.RequestURI                                // 请求路径
		requestId := Md5(requestUri, fmt.Sprintf("%d", requestTime.UnixNano())) // 生成 requestId
		logger := vdlog.Clone().SetLogId(requestId)                             // 克隆 logger并设置 requestId

		// 构造入参数据
		request := map[string]interface{}{
			"time":     requestTime,
			"method":   context.Request.Method,
			"query":    requestUri,
			"header":   context.Request.Header.Clone(),
			"body":     map[string]interface{}{},
			"clientIP": context.ClientIP(),
		}

		// 处理 body 数据
		if context.Request.Method == http.MethodPost {
			rawData, _ := context.GetRawData()
			switch context.ContentType() {
			case "application/json":
				var body map[string]interface{}
				_ = json.Unmarshal(rawData, &body)
				request["body"] = body
			case "application/x-www-form-urlencoded":
				urlValue, _ := url.ParseQuery(string(rawData))
				request["body"] = urlValue
			}
			// body回放
			context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
		}

		// 记录请求日志
		logger.Info(request)          // 记录请求数据
		context.Set("logger", logger) // 存储日志

		// 处理请求
		context.Next()

		// 响应数据
		response, _ := context.Get("response") // 响应数据
		logger.Info(response)                  // 记录响应数据
	})
	return nil
}

// [info]日志写入到文件
func LoggerToFile() gin.HandlerFunc {
	// 文件存储
	return func(c *gin.Context) {
		fileName := CreateFile("info")
		logger := Logger(fileName)
		//contents := c.Request.Body
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 请求头
		userAgent, _ := json.Marshal(c.Request.Header)
		// post请求
		_ = c.Request.ParseForm()
		requestInfo, _ := json.Marshal(c.Request.PostForm)

		//strings.EqualFold(Config.Pattern,"debug")
		if Config.Pattern == "debug" {
			//  时间;唯一id;项目名称;模块名;结果
			var onlyStr = Md5(latencyTime.String(), Config.ProjectName)
			logger.WithFields(WriteFields(onlyStr, "req_uri", reqUri)).Info()
			logger.WithFields(WriteFields(onlyStr, "request_ip", clientIP)).Info()
			logger.WithFields(WriteFields(onlyStr, "post_request", string(requestInfo))).Info()
			logger.WithFields(WriteFields(onlyStr, "user_agent", string(userAgent))).Info()
			SqlContext, ok := c.Get("SQL")
			if ok {
				for _, value := range SqlContext.([]interface{}) {
					sqlValue := value.(map[string]string)
					logger.WithFields(WriteFields(onlyStr, "runsql", sqlValue["runsql"]+" runtime"+sqlValue["runtime"])).Info()
				}
			}
			// logger.WithFields(WriteFields(onlyStr,"response_info",string(backInfo))).Info()
		} else {
			logger.WithFields(logrus.Fields{
				"status_code":  statusCode,
				"latency_time": latencyTime,
				"client_ip":    clientIP,
				"req_method":   reqMethod,
				"req_uri":      reqUri,
			}).Info()
		}
	}
}

// 写入错误
func ErrorWrite(errorInfo string) {
	fmt.Println(errorInfo)
	fileName := CreateFile("error")
	logger := Logger(fileName)
	fmt.Println(errorInfo)
	logger.WithFields(logrus.Fields{
		"ErrorInfo": errorInfo,
	}).Error()
}

// 日志列表-记录格式
func WriteFields(onlyStr string, key string, keyInfo string) logrus.Fields {
	return logrus.Fields{
		"only_key":     onlyStr,
		"project_name": Config.ProjectName,
		key:            keyInfo,
	}
}

//  执行的sql或者返回结果存储,日志动态输出
func LogWrite(c *gin.Context, LogKey string, info interface{}) {
	switch LogKey {
	case "POST":
		c.Set(LogKey, info)
	case "GET":
		c.Set(LogKey, info)
	case "SQL":
		SqlContext, ok := c.Get("SQL")
		if !ok {
			var SqlRs []interface{}
			SqlRs = append(SqlRs, info)
			c.Set(LogKey, SqlRs)
		} else {
			var SqlCenter []interface{}
			for _, value := range SqlContext.([]interface{}) {
				SqlCenter = append(SqlCenter, value)
			}
			SqlRs := make([]interface{}, len(SqlCenter), (cap(SqlCenter))*2)
			copy(SqlRs, SqlCenter)
			SqlRs = append(SqlRs, info)
			c.Set(LogKey, SqlRs)
		}
	}
}

// 日志句柄初始化
func Logger(fileName string) *logrus.Logger {
	//写入文件
	fileHanddle, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = fileHanddle

	//设置日志级别
	//logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	if Config.StoreFormat == 2 {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	return logger
}

// 字符串加密
func Md5(str ...string) string {
	md5Ptr := md5.New()
	_, _ = md5Ptr.Write([]byte(strings.Join(str, "")))
	return fmt.Sprintf("%x", md5Ptr.Sum(nil))
}

// 创建文件
func CreateFile(level string) string {
	// 每天写入新文件需要处理
	now := time.Now()
	logFilePath := ""
	fmt.Println()
	logFilePath = Config.RuntimePath + "/log/"
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + "." + level + ".log"

	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	return fileName
}

// 日志写入 syslog
func LoggerToSyslog() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志写入 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
