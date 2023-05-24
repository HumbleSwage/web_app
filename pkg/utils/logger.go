package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

var LogrusObj *logrus.Logger

type LogFormatter struct {
}

func (lf *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// 定位行号和调用者
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	var err error
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		_, err = fmt.Fprintf(buffer, "[%s] [%s] %s %s %s\n", timestamp, entry.Level, funcVal, fileVal, entry.Message)
	} else {
		_, err = fmt.Fprintf(buffer, "[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}
	return buffer.Bytes(), err
}

func InitLog() {
	if LogrusObj != nil {
		file, err := setOutFile()
		if err != nil {
			panic(err)
		}
		// 设置输出
		LogrusObj.Out = file
		return
	}
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel) // 设置最低的level
	logger.SetReportCaller(true)       // 开启返回函数和行号
	logger.SetFormatter(&LogFormatter{})
}

func setOutFile() (*os.File, error) {
	now := time.Now()
	// 获取当前目录
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs"
	}
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		if err := os.Mkdir(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	// 设置日志文件
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if _, err = os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	// 打开日志文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return file, nil
}
