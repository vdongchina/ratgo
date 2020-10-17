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
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Error info.
func Error(recover interface{}) error {
	var array [4096]byte
	buf := array[:]
	runtime.Stack(buf, false)
	stackString := Stack(recover, buf)
	return errors.New(stackString)
}

// Handle Stack message.
func Stack(err interface{}, buf []byte) string {
	var stackSlice []string
	var msg string

	// Filename、line
	if _, srcName, line, ok := runtime.Caller(4); ok {
		msg = fmt.Sprintf("[Error]File:%s Line:%s Msg:%s \nMethod Stack meassage:", srcName, strconv.Itoa(line), err)
		stackSlice = append(stackSlice, msg)
	} else {
		msg = fmt.Sprintf("[Error]Msg:%s \nMethod Stack meassage:", err)
		stackSlice = append(stackSlice, msg)
	}

	// Handle data.
	stringStack := string(buf)                   // Convert to string
	tmpStack := strings.Split(stringStack, "\n") // Split to []string by \n
	var receiveStack []string
	for _, v := range tmpStack {
		receiveStack = append(receiveStack, strings.TrimSpace(v))
	}
	for i, j, k := 0, 0, len(receiveStack)-1; i < k; i += 2 {
		stackSlice = append(stackSlice, "["+strconv.Itoa(j)+"]"+receiveStack[i]+" "+receiveStack[i+1])
		if j == 10 {
			stackSlice = append(stackSlice, "...")
			break
		}
		j++
	}
	return strings.Join(stackSlice, "\n")
}