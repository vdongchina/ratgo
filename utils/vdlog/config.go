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

const (
	Info    = 1 << iota // 日志级别: 信息
	Debug               // 日志级别: 调试
	Warning             // 日志级别: 警告
	Error               // 日志级别: 错误
)

// 配置
type Config struct {
	Turn       string   // 是否开启 on:是 off:否
	RootPath   string   // 日志文件根路径
	Ext        string   // 日志后缀
	Mode       []string // 模式 [Info,Warning,Debug,Error]
	ServerName string   // 服务名称
	ServerIP   string   // 服务IP
	Platform   string   // 平台类型
	Format     string   // 日志格式化
}

var defaultConfig *Config

func init() {
	defaultConfig = &Config{
		Turn:     "on",
		Format:   "Datetime,LogId,FilePath",
		Ext:      "log",
		RootPath: `/var/logs`,
	}
}
