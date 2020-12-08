package log

import (
	"fmt"
	"log"
)

type Logger struct {
	log.Logger
	FilePath string
	FileName string



}

func NewLogger() *Logger{
	logger := &Logger{}



	return logger
}


func (lw *Logger) Write(p []byte) (n int, err error) {





	fmt.Println(string(p))

	return
}





