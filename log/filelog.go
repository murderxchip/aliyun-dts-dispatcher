package log

import (
	"encoding/json"
	"github.com/murderxchip/aliyun-dts-dispatcher/log/log4go"
	"os"
	"runtime"
	"strconv"
)

//level CRITICAL>ERROR>WARNING>INFO>TRACE>DEBUG
func Info(tag string, v ...interface{}) {
	Log(log4go.INFO, tag, v...)
}

func Debug(tag string, v ...interface{}) {
	Log(log4go.DEBUG, tag, v...)
}

func Error(tag string, v ...interface{}) {
	Log(log4go.ERROR, tag, v...)
}

func Trace(tag string, v ...interface{}) {
	Log(log4go.TRACE, tag, v...)
}

func Warning(tag string, v ...interface{}) {
	Log(log4go.WARNING, tag, v...)
}

func Fatal(tag string, v ...interface{}) {
	Log(log4go.CRITICAL, tag, v...)
}

func Log(Level log4go.Level, tag string, v ...interface{}) {
	var logInfo string
	message := getLogMessage(v...)
	source := getPosition()
	if message != "" {
		logInfo = tag + "\t" + message
	} else {
		logInfo = tag
	}
	switch Level {
	case log4go.DEBUG:
		Logger.Log(log4go.DEBUG, source, logInfo)
	case log4go.TRACE:
		Logger.Log(log4go.TRACE, source, logInfo)
	case log4go.INFO:
		Logger.Log(log4go.INFO, source, logInfo)
	case log4go.WARNING:
		Logger.Log(log4go.WARNING, source, logInfo)
	case log4go.ERROR:
		Logger.Log(log4go.ERROR, source, logInfo)
	case log4go.CRITICAL:
		Logger.Log(log4go.CRITICAL, source, logInfo)
		panic(message)
	}
}

func getLogMessage(messages ...interface{}) string {
	if len(messages) < 1 {
		return ""
	}
	var logs []interface{}
	for _, v := range messages {
		logs = append(logs, v)
	}
	s, _ := json.Marshal(logs)
	return string(s)
}

func getPosition() string {
	position := ""
	_, file, line, ok := runtime.Caller(3)
	if ok {
		position = file + ":" + strconv.Itoa(line)
	}
	hostName, _ := os.Hostname()
	source := hostName + ":" + position
	return source
}
