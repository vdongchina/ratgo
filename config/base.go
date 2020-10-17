package config

type BaseParser interface {
	Init() BaseParser
	ParserToMap(filePath string) (anyMap map[string]interface{}, err error)
}