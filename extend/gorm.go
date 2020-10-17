// Copyright 2020 ratgo Author. All Rights Reserved.
// Licensed under the Apache License, Version 1.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package extend

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/vdongchina/ratgo/utils/types"
	"sync"
	"time"
)

// *gorm.DB 存储容器
type DbStorage struct {
	lock   *sync.RWMutex
	anyMap types.AnyMap
	dbMap  map[string]*gorm.DB
}

// *gorm.DB 连接对象配置
type DbConfig struct {
	DriverName     string
	DataSourceName string
	MaxOpenConn    int
	MaxIdleConn    int
	MaxLifetime    time.Duration
	SingularTable  string
}

// 定义
var Gorm *DbStorage

func init() {
	Gorm = &DbStorage{
		lock:   new(sync.RWMutex),
		anyMap: types.AnyMap{},
		dbMap:  map[string]*gorm.DB{},
	}
}

// 配置初始化 - map形式
func (ds *DbStorage) Init(config types.AnyMap) {
	ds.anyMap = config
}

// 获取配置
func (ds *DbStorage) GetConfig(key string) map[string]DbConfig {
	config := ds.anyMap.Get(key).ToStringMap()

	// 读取数据库配置并校验
	configParam := map[string]string{"driverName": "", "dataSourceName": "", "maxIdleConn": "", "maxOpenConn": "", "maxLifetime": ""}
	for k := range configParam {
		if v, ok := config[k]; !ok {
			panic("The database config's " + k + " is not set.")
		} else {
			configParam[k] = v
		}
	}

	// SingularTable 设置
	if value, ok := config["SingularTable"]; ok {
		configParam["SingularTable"] = value
	} else {
		configParam["SingularTable"] = "false"
	}

	// 赋值
	dc := map[string]DbConfig{key: {
		DriverName:     configParam["driverName"],
		DataSourceName: configParam["dataSourceName"],
		MaxOpenConn:    types.Eval(configParam["maxIdleConn"]).ToInt(),
		MaxIdleConn:    types.Eval(configParam["maxIdleConn"]).ToInt(),
		MaxLifetime:    1800,
		SingularTable:  configParam["SingularTable"],
	}}
	return dc
}

// 配置初始化 - 结构体map形式
func (ds *DbStorage) InitDb(config map[string]DbConfig) {
	// 存储Db对象
	for key, value := range config {
		ds.dbMap[key] = ds.gormDb(value)
	}
}

// 使用 *gorm.DB
func (ds *DbStorage) Db(key string) *gorm.DB {
	if db, ok := ds.dbMap[key]; ok {
		return db
	} else {
		ds.lock.Lock()
		defer ds.lock.Unlock()
		if db, ok := ds.dbMap[key]; ok {
			return db
		}
		// add *gorm.DB
		ds.InitDb(ds.GetConfig(key))
		if db, ok := ds.dbMap[key]; ok {
			return db
		}
	}
	return nil
}

// 根据配置获取 *gorm.DB
func (ds *DbStorage) gormDb(config DbConfig) *gorm.DB {
	gormDB, err := gorm.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		panic(fmt.Sprintf("DB connect faild err: %v", err))
	}

	// 连接池设置
	gormDB.DB().SetMaxIdleConns(config.MaxIdleConn)    // SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	gormDB.DB().SetMaxOpenConns(config.MaxOpenConn)    // SetMaxOpenConns 设置打开数据库连接的最大数量。
	gormDB.DB().SetConnMaxLifetime(config.MaxLifetime) // SetConnMaxLifetime 设置了连接可复用的最大时间。

	// SingularTable 设置
	if config.SingularTable == "true" {
		gormDB.SingularTable(true)
	}
	return gormDB
}
