// Copyright 2020 ratgo Author. All Rights Reserved.
// Licensed under the Apache License, Version 1.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package vdlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var DefaultLogger *Logger

// 初始化
func init() {
	DefaultLogger = NewLogger(*defaultConfig)
}

// 日志结构体
type Logger struct {
	Date        string         // 日期
	LogId       interface{}    // 日志id
	Config      *Config        // 日志配置
	lock        sync.Mutex     // 排斥锁
	InfoLogger  *logrus.Logger // logrus.Logger
	WarnLogger  *logrus.Logger // logrus.Logger
	ErrorLogger *logrus.Logger // logrus.Logger
	DebugLogger *logrus.Logger // logrus.Logger
}

// 更新默认logger并返回
func Default(outConfig interface{}) *Logger {
	DefaultLogger = NewLogger(outConfig)
	return DefaultLogger
}

// 根据类型获取
func NewLogger(outConfig interface{}) *Logger {
	var config Config
	switch outConfig.(type) {
	case Config:
		config = outConfig.(Config)
	case map[string]interface{}:
		jsonMap, _ := json.Marshal(outConfig)
		_ = json.Unmarshal(jsonMap, &config)
	default:
		config = *defaultConfig
	}

	// 创建实例 logrus.logger
	infoLogger := logrus.New()
	infoLogger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	//infoLogger.SetReportCaller(true)
	warnLogger := logrus.New()
	warnLogger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	debugLogger := logrus.New()
	debugLogger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	errorLogger := logrus.New()
	errorLogger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})

	// 创建实例 vdlog.logger
	curtDate := time.Now().Format("2006010215")
	logger := &Logger{
		Date:        curtDate,
		Config:      &config,
		lock:        sync.Mutex{},
		InfoLogger:  infoLogger,
		WarnLogger:  warnLogger,
		DebugLogger: debugLogger,
		ErrorLogger: errorLogger,
	}

	// 设置 io.writer
	logger.setIoWriter(logger.InfoLogger, "info", curtDate)
	logger.setIoWriter(logger.WarnLogger, "warn", curtDate)
	logger.setIoWriter(logger.DebugLogger, "debug", curtDate)
	logger.setIoWriter(logger.ErrorLogger, "error", curtDate)
	return logger
}

// 克隆
func Clone() *Logger {
	logger := &Logger{
		Date:        time.Now().Format("2006010215"),
		Config:      DefaultLogger.Config,
		InfoLogger:  DefaultLogger.InfoLogger,
		WarnLogger:  DefaultLogger.WarnLogger,
		ErrorLogger: DefaultLogger.ErrorLogger,
		DebugLogger: DefaultLogger.DebugLogger,
	}

	// 设置 io.writer
	logger.setIoWriter(logger.InfoLogger, "info", logger.Date)
	logger.setIoWriter(logger.WarnLogger, "warn", logger.Date)
	logger.setIoWriter(logger.DebugLogger, "debug", logger.Date)
	logger.setIoWriter(logger.ErrorLogger, "error", logger.Date)
	return logger
}

// 设置logId
func (ls *Logger) SetLogId(logId interface{}) *Logger {
	ls.LogId = logId
	return ls
}

// 记录info日志
func (ls *Logger) Info(content interface{}) {
	ls.InfoLogger.WithFields(ls.formatMap(content)).Info()
}

// 记录警告信息
func (ls *Logger) Warning(content interface{}) {
	ls.WarnLogger.WithFields(ls.formatMap(content)).Warning()
}

// 记录警告信息
func (ls *Logger) Debug(content interface{}) {
	ls.DebugLogger.WithFields(ls.formatMap(content)).Debug()
}

// 记录警告信息
func (ls *Logger) Error(content interface{}) {
	ls.ErrorLogger.WithFields(ls.formatMap(content)).Error()
}

// 设置 io.writer
func (ls *Logger) setIoWriter(logger *logrus.Logger, level string, date string) *logrus.Logger {
	if ioWriter, err := ls.ioWriter(level, date); err != nil {
		fmt.Printf("create %s io.writer failed. error:%v\n", level, err)
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(ioWriter)
	}
	return logger
}

// 创建文件 io.Writer
func (ls *Logger) ioWriter(level string, date string) (osFile *os.File, err error) {
	// 判断根目录是否存在
	if _, ok := IsExists(ls.Config.RootPath); !ok {
		return nil, errors.New(fmt.Sprintf("log path '%s' is not exists.", ls.Config.RootPath))
	}

	// 文件路径
	paths := []string{date, level}
	baseName := fmt.Sprintf("%s.%s", strings.Join(paths, "-"), ls.Config.Ext)
	filename := path.Join(ls.Config.RootPath, baseName)
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
}

// 记录警告信息
func (ls *Logger) write(level string, content interface{}) {
	// 判断目录是否存在
	filePath := ls.Config.RootPath
	if !FileExists(filePath) {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			fmt.Printf("create dir '%s' failed, error:%v\n", filePath, err)
			return
		}
	}

	// 日志文件
	paths := []string{ls.get("Date").(string)}
	paths = append(paths, level)
	baseName := fmt.Sprintf("%s.%s", strings.Join(paths, "-"), ls.Config.Ext)
	filename := path.Join(filePath, baseName)

	// 文件句柄
	_, isFile := IsFile(filename)
	var osFile *os.File
	var err error
	if isFile { // 打开文件，
		osFile, _ = os.OpenFile(filename, os.O_APPEND, 0777)
	} else { // 新建文件
		osFile, err = os.Create(filename)
	}

	//使用完毕，需要关闭文件
	defer func() {
		err = osFile.Close()
		if err != nil {
			fmt.Printf("close file '%s' failed. error:%v\n", filename, err)
		}
	}()
	if err != nil {
		fmt.Printf("handle file '%s' failed. error:%v\n", filename, err)
		return
	}

	// 格式化日志内容并写入文件
	_, err = osFile.WriteString(ls.format(level, content))
	if err != nil {
		fmt.Printf("write file '%s' failed. error:%v\n", filename, err)
	}
}

// 判断路径是否存在
func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

// 判断所给路径是否为文件夹
func IsDir(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && f.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && !f.IsDir()
}

// 判断文件是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// json序列化
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// 获取格式化数据
func (ls *Logger) format(level string, content interface{}) (format string) {
	format = ls.Config.Format
	format = strings.ReplaceAll(format, "{Datetime}", fmt.Sprintf("%v", ls.get("Datetime")))
	format = strings.ReplaceAll(format, "{ServerName}", fmt.Sprintf("%v", ls.get("ServerName")))
	format = strings.ReplaceAll(format, "{Level}", strings.ToUpper(level))
	format = strings.ReplaceAll(format, "{LogId}", fmt.Sprintf("%v", ls.get("LogId")))
	format = strings.ReplaceAll(format, "{FilePath}", fmt.Sprintf("%v", ls.get("FilePath")))
	jsonContent, _ := JSONMarshal(content)
	format = strings.ReplaceAll(format, "{Content}", fmt.Sprintf("%s", string(jsonContent)))
	format = strings.ReplaceAll(format, "{Platform}", fmt.Sprintf("%v", ls.get("Platform")))
	format = strings.ReplaceAll(format, "{ServerIP}", fmt.Sprintf("%v", ls.get("ServerIP")))
	format = fmt.Sprintf("%s\r\n", format)
	return
}

// 获取格式化数据
func (ls *Logger) formatMap(content interface{}) (format map[string]interface{}) {
	format = map[string]interface{}{}
	formatSlice := strings.Split(ls.Config.Format, ",")
	for _, v := range formatSlice {
		format[v] = ls.get(v)
	}
	format["content"] = content
	return
}

// 记录警告信息
func (ls *Logger) get(key string) (value interface{}) {
	curtTime := time.Now()
	switch key {
	case "Date":
		value = curtTime.Format("20060102")
	case "Datetime":
		value = curtTime.Format("2006-01-02 15:04:05")
	case "ServerName":
		value = ls.Config.ServerName
	case "Platform":
		value = ls.Config.Platform
	case "ServerIP":
		value = ls.Config.ServerIP
	case "LogId":
		value = ls.LogId
	case "FilePath":
		callerFile := CallTrack("vdlog")
		value = fmt.Sprintf("%s:%v", callerFile.File, callerFile.Line)
	default:
		value = ""
	}
	return
}
