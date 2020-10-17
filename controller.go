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

import "github.com/gin-gonic/gin"

// 控制器接口
type ControllerInterface interface {
	Init(ctx *gin.Context, result *Result) // 初始化
	Context() *gin.Context                 // 获取gin的 Context
	BeforeExec()                           // 动作前执行方法
	Exec()                                 // 动作方法
	Result() *Result                       // 控制响应(辅助用法,阻断执行直接响应异常)
}

// 控制器
type Controller struct {
	context *gin.Context
	result  *Result
}

// 控制器初始化方法
func (c *Controller) Init(context *gin.Context, result *Result) {
	c.context = context
	c.result = result
}

// 简易模式下执行方法
func (c *Controller) Context() *gin.Context {
	return c.context
}

// 简易模式 - 前置方法
func (c *Controller) BeforeExec() {

}

// 简易模式 - 执行方法
func (c *Controller) Exec() {

}

// 结果方法
func (c *Controller) Result() *Result {
	return c.result
}
