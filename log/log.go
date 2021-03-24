package log

import (
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/log/log4go"
	"strings"
	"time"
)

var Logger log4go.Logger

func NewLogConfig() {
	Logger = make(log4go.Logger)
	logName := config.Config().Log.LogPath
	level := log4go.INFO
	switch strings.ToTitle(config.Config().Log.Level) {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "FATAL":
		level = log4go.CRITICAL
	}

	if strings.ToLower(logName) == "stdout" {
		logConfig := log4go.NewConsoleLogWriter()
		Logger.AddFilter("logfile", level, logConfig)
	} else {
		logFileName := logName[:len(logName)-4] + "-" + time.Now().Format("2006-01-02") + ".log"
		logConfig := log4go.NewFileLogWriter(
			logFileName,
			config.Config().Log.Rotate,
		)
		logConfig.SetRotateDaily(config.Config().Log.RotateDaily)
		Logger.AddFilter("logfile", level, logConfig)
	}
}
