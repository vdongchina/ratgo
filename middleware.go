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
)

// 中间件容器
type MiddleWareStorage struct {
	global []gin.HandlerFunc
	group  map[string][]gin.HandlerFunc
}

// 初始化
var MiddleWare *MiddleWareStorage

func init() {
	MiddleWare = &MiddleWareStorage{
		global: make([]gin.HandlerFunc, 0),
		group:  map[string][]gin.HandlerFunc{},
	}
}

// 设置全局中间件
func (mws *MiddleWareStorage) SetGlobal(handlerFunc ...gin.HandlerFunc) {
	mws.global = append(mws.global, handlerFunc...)
}

// 批量设置全局中间件
func (mws *MiddleWareStorage) BatchSetGlobal(handlerFuncSlice []gin.HandlerFunc) {
	mws.global = append(mws.global, handlerFuncSlice...)
}

// 获取全局中间件
func (mws *MiddleWareStorage) GetGlobal() []gin.HandlerFunc {
	return mws.global
}

// 设置分组中间件
func (mws *MiddleWareStorage) SetGroup(groupName string, handlerFunc ...gin.HandlerFunc) {
	mws.group[groupName] = append(mws.group[groupName], handlerFunc...)
}

// 批量设置分组中间件
func (mws *MiddleWareStorage) BatchSetGroup(groupName string, handlerFuncSlice []gin.HandlerFunc) {
	mws.group[groupName] = append(mws.group[groupName], handlerFuncSlice...)
}

// 获取分组组中间件
func (mws *MiddleWareStorage) GetGroup() map[string][]gin.HandlerFunc {
	return mws.group
}
