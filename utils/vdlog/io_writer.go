package vdlog

import "encoding/json"

// ioWriter接口
type IoWriterInterface interface {
	Init(config map[string]interface{})
	Write(p []byte) (n int, err error)
}

// 文件ioWriter
type FileWriter struct {
	RootPath   string `json:"rootPath"`   // 日志文件根路径
	MiddlePath string `json:"middlePath"` // 日志文件中间路径 例如: module-list
	Ext        string `json:"ext"`        // 日志后缀
}

// 初始化
func (fw *FileWriter) Init(config map[string]interface{}) {
	byteConfig, _ := json.Marshal(config)
	_ = json.Unmarshal(byteConfig, fw)
	return
}

// 设置日志中间路径
func (fw *FileWriter) SetMiddlePath(middlePath string) *FileWriter {
	fw.MiddlePath = middlePath
	return fw
}

// 日志输出
func (fw *FileWriter) Write(p []byte) (n int, err error) {
	return
}
