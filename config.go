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
	"github.com/gin-gonic/gin"
	"github.com/vdongchina/ratgo/config"
	"github.com/vdongchina/ratgo/utils/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ConfigStorage.
type ConfigStorage struct {
	AppName        string // Application name
	AppVersion     string // Application version
	RunMode        string // Running Mode: dev | test | prod
	AppPath        string
	ConfigPath     string
	RuntimePath    string
	RuntimeLogPath string
	EnableHTTP     bool
	HTTPAddr       string
	EnableHTTPS    bool
	HTTPSAddr      string
	HTTPSCertFile  string
	HTTPSKeyFile   string
	RecoverPanic   bool
	RecoverFunc    func(c *gin.Context)
	DefinedConfig  types.AnyMap
	Error          error
	Pattern        string // debug:list; release
	StoreFormat    int    // 1:text; 2:json;
	Storage        string // local; syslog; redis; es; mongo
	ProjectName    string // 项目名称
	InitDb         bool   // 是否初始化 gorm db
	InitRedis      bool   // 是否初始化 redis
	HandleFunc     func(config *ConfigStorage) error
}

var (
	// BConfig is the default config for Application
	Config *ConfigStorage

	// configProvider is the provider for the config, default is ini
	configProvider = "yml"
)

// 初始化 ratgo配置
func init() {
	Config = &ConfigStorage{
		AppName:        "ratgo",
		RunMode:        "dev",
		AppPath:        "",
		ConfigPath:     "",
		RuntimePath:    "",
		RuntimeLogPath: "",
		HTTPAddr:       ":8080",
		HTTPSAddr:      ":10443",
		HTTPSCertFile:  "",
		HTTPSKeyFile:   "",
		DefinedConfig:  types.AnyMap{},
		Pattern:        "debug",
		StoreFormat:    1,
		Storage:        "local",
		ProjectName:    "lianxilog",
		InitDb:         true,
		InitRedis:      true,
		HandleFunc:     nil,
	}
}

// 初始化
func (cs *ConfigStorage) Init(outAnyMap ...map[string]interface{}) {
	// Run Mode Include: dev、test、prod
	if cs.RunMode = os.Getenv("RATGO_RUNMODE"); cs.RunMode == "" {
		cs.RunMode = "dev"
	}
	// Application Path ... flag get param yet
	if cs.AppPath, cs.Error = os.Getwd(); cs.Error != nil {
		panic(cs.Error)
	}

	// 路径map
	pathMap := map[string]string{
		"ConfigPath":     "config." + Config.RunMode,
		"RuntimePath":    "runtime",
		"RuntimeLogPath": "runtime.log",
	}
	for key, value := range pathMap {
		realpath := filepath.Join(cs.AppPath, filepath.Join(strings.Split(value, ".")...))
		if cs.Error = cs.DirExists(realpath); cs.Error != nil { // 不存在则创建
			if err := os.MkdirAll(realpath, os.ModePerm); err != nil {
				panic(err)
			}
		}
		pathMap[key] = realpath
	}

	// Scan dir list & Load config
	if outAnyMap != nil {
		cs.DefinedConfig = outAnyMap[0]
	} else {
		configFiles := cs.ScanDir(pathMap["ConfigPath"])
		for _, file := range configFiles {
			fileName := file.Name()
			filePath := filepath.Join(pathMap["ConfigPath"], fileName)
			fileSplit := strings.Split(fileName, ".")
			if fileSplit[1] == configProvider { // 配置类型过滤后缀名为 configProvider
				if unitConfig, err := cs.ParserToAnyMap(configProvider, filePath); err != nil {
					panic(err.Error())
				} else {
					cs.DefinedConfig[fileSplit[0]] = unitConfig
				}
			}
		}
	}

	// App配置
	appConfig := types.AnyMap(cs.Get("ratgo.main").ToAnyMap())
	if len(appConfig) > 0 {
		delete(cs.DefinedConfig, "ratgo")
	}

	// 更新服务配置
	cs.MapMerge(appConfig, pathMap)
	pt := reflect.TypeOf(cs).Elem()
	pv := reflect.ValueOf(cs).Elem()
	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		// 字段重新赋值
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			if v := appConfig.Get(name).Value(); v != nil {
				pf.SetString(appConfig.Get(name).ToString())
			}
		case reflect.Int, reflect.Int64:
			if v := appConfig.Get(name).Value(); v != nil {
				pf.SetInt(appConfig.Get(name).ToInt64())
			}
		case reflect.Bool:
			if v := appConfig.Get(name).Value(); v != nil {
				pf.SetBool(appConfig.Get(name).ToBool())
			}
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}

// 配置处理方法
func (cs *ConfigStorage) Handle(fn func(config *ConfigStorage) error) {
	cs.HandleFunc = fn
}

// 使用解析器解析文件至map类型配置数据
func (cs *ConfigStorage) ParserToAnyMap(configProvider string, filePath string) (map[string]interface{}, error) {
	var configParser config.BaseParser
	switch configProvider {
	case "ini":
		configParser = &config.IniParser{}
	case "json":
		configParser = &config.JsonParser{}
	case "yml":
		configParser = &config.YamlParser{}
	default:
		configParser = &config.IniParser{}
	}
	return configParser.ParserToMap(filePath)
}

// 合并map, 不支持递归合并
func (cs *ConfigStorage) MapMerge(dstMap map[string]interface{}, merged ...interface{}) {
	if merged == nil {
		return
	}
	for _, value := range merged {
		switch value.(type) {
		case map[string]interface{}:
			for k, v := range value.(map[string]interface{}) {
				dstMap[k] = v
			}
		case map[string]string:
			for k, v := range value.(map[string]string) {
				dstMap[k] = v
			}
		default:
			continue
		}
	}
}

// Get app config.
func (cs *ConfigStorage) Get(args ...string) *types.AnyValue {
	return cs.DefinedConfig.Get(args...)
}

// Set value.
func (cs *ConfigStorage) Set(args string, value interface{}) {
	cs.DefinedConfig.Set(args, value)
}

// Judge dir is exist.
func (cs *ConfigStorage) DirExists(path string) (err error) {
	_, err = os.Stat(path)
	return err
}

// 遍历目录
func (cs *ConfigStorage) ScanDir(dirPath string) (fileSlice []os.FileInfo) {
	// 判断目录是否存在
	if err := cs.DirExists(dirPath); err != nil {
		panic(err)
	}

	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(err)
	}
	defer func() { _ = dir.Close() }()

	// 读取目录文件
	if files, err := dir.Readdir(-1); err != nil {
		panic(err)
	} else {
		fileSlice = files
	}
	return fileSlice
}
