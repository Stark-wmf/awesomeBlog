package core

import (
	"fmt"
	"github.com/cihub/seelog"
)

var logger seelog.LoggerInterface

func InitLog(fileName string) error {
	var err error
	logger, err = seelog.LoggerFromConfigAsFile(fileName)
	if err != nil {
		return err
	}

	if e := seelog.ReplaceLogger(logger); e != nil {
		return e
	}
	return nil
}

func Debug(content string) {
	logger.Debug(content)
}

func Info(content string) {
	logger.Info(content)
}

func Error(content string) {
	logger.Error(content)
	fmt.Println(content)
}

func Errorf(content string, param ...interface{}) {
	logger.Errorf(content, param)
}
func FlushLog() {
	logger.Flush()
}
