package vdlog

import (
	"runtime"
	"strings"
)

var DefaultSkipPackage = "runtime"

type CallerFile struct {
	File string
	Line int
	Func string
}

// 调用追踪
func CallTrack(skipPackage ...string) *CallerFile {
	skipPackage = append(skipPackage, DefaultSkipPackage)
	auto := &CallerFile{}
	for skip := 1; true; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if ok {
			auto.File = file
			auto.Line = line
			auto.Func = runtime.FuncForPC(pc).Name()
			if !InSlice(skipPackage, auto.Func) {
				break
			}
		} else {
			break
		}
	}
	return auto
}

// 判断元素是否存在
func InSlice(stringSlice []string, str string) (inSlice bool) {
	for _, v := range stringSlice {
		if strings.Contains(str, v) {
			inSlice = true
			break
		}
	}
	return
}
