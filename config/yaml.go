package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type YamlParser struct {
	BaseParser
}

func (jp *YamlParser) Init() BaseParser {
	return jp
}

// 解析文件
func (jp *YamlParser) ParserToMap(filePath string) (map[string]interface{}, error) {
	// 从配置文件中读取json字符串
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("load config conf file '"+filePath+"'failed: %s", err.Error())
		return nil, errors.New(errorMsg)
	}

	// 存储map
	anyMap := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &anyMap)
	if err != nil {
		errorMsg := fmt.Sprintf("decode config file '"+filePath+"' failed: %s", err.Error())
		return nil, errors.New(errorMsg)
	}
	return anyMap, nil
}
