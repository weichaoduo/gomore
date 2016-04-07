// main loop

package golog

import (
	"fmt"
	"gomore/global"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/weekface/mgorus"
	//"github.com/rifflock/lfshook"
)

var Log *logrus.Logger

// 初始化日志设置
func InitLogger() {

	if runtime.GOOS != "windows" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	Log = logrus.New()

	fmt.Println("LogBehindType", global.Config.MyLog.LogBehindType)
	if global.Config.MyLog.LogBehindType == "mongodb" {

		mongodb_server := global.Config.MyLog.MongodbHost + ":" + global.Config.MyLog.MongodbPort
		hooker, err := mgorus.NewHooker(mongodb_server, "db", "collection")
		if err == nil {
			Log.Hooks.Add(hooker)
		} else {
			fmt.Println("mongodb err:", err)
		}
	}

	// init logger
	loglevel := global.Config.MyLog.LogLevel
	if loglevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if loglevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	if loglevel == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	}
	if loglevel == "warn" {
		logrus.SetLevel(logrus.WarnLevel)
	}
	if loglevel == "fatal" {
		logrus.SetLevel(logrus.FatalLevel)
	}
	if loglevel == "panic" {
		logrus.SetLevel(logrus.PanicLevel)
	}

	fmt.Println("logger status : ", loglevel, runtime.GOOS)

}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	Log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	Log.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	Log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}
