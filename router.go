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

import "strings"

// 简易模式路由容器
type GeneralMap map[string]ControllerInterface

// 路由存储器
type RouterStorage struct {
	Mode          string
	generalMap    GeneralMap
	generalPath   map[string][]string
	staticMap     map[string]string
	staticFileMap map[string]string
	htmlGlob      []string
	html          []string
}

// router对象
var Router *RouterStorage

func init() {
	Router = &RouterStorage{
		Mode:          "General",
		generalMap:    GeneralMap{},
		generalPath:   map[string][]string{},
		staticMap:     map[string]string{},
		staticFileMap: map[string]string{},
		htmlGlob:      make([]string, 0),
	}
}

// 设置简易路由
func (rs *RouterStorage) General(group string, controllers GeneralMap) {
	if rs.Mode != "General" {
		rs.Mode = "General"
	}
	for key, value := range controllers {
		generalKey := strings.Join([]string{"", group, key}, "/")
		rs.generalMap[generalKey] = value
	}
	if _, ok := rs.generalPath[group]; !ok {
		rs.generalPath[group] = []string{
			strings.Join([]string{"", group, ":param1"}, "/"),
			strings.Join([]string{"", group, ":param1", ":param2"}, "/"),
		}
	}
}

// 获取简易路由对应控制器
func (rs *RouterStorage) GetGeneral(path string) ControllerInterface {
	if controller, ok := rs.generalMap[path]; ok {
		return controller
	}
	return nil
}

// 获取简易路由 relativePath数组
func (rs *RouterStorage) GetGeneralPath() map[string][]string {
	return rs.generalPath
}

// 设置静态路径
func (rs *RouterStorage) SetStatic(relativePath, root string) {
	rs.staticMap[relativePath] = root
}

// 批量设置静态路径 map
func (rs *RouterStorage) BatchSetStatic(staticMap map[string]string) {
	for key, value := range staticMap {
		rs.SetStatic(key, value)
	}
}

// 获取静态路径 map
func (rs *RouterStorage) GetStatic() map[string]string {
	return rs.staticMap
}

// 静态文件
func (rs *RouterStorage) SetStaticFile(relativePath, root string) {
	rs.staticFileMap[relativePath] = root
}

// 批量设置静态文件 map
func (rs *RouterStorage) BatchSetStaticFile(staticFileMap map[string]string) {
	for key, value := range staticFileMap {
		rs.SetStaticFile(key, value)
	}
}

// 获取静态文件 map
func (rs *RouterStorage) GetStaticFile() map[string]string {
	return rs.staticFileMap
}

// 设置 html模板
func (rs *RouterStorage) SetHTMLGlob(pattern string) {
	rs.htmlGlob = append(rs.htmlGlob, pattern)
}

// 获取 html模板
func (rs *RouterStorage) GetHTMLGlob() []string {
	return rs.htmlGlob
}
