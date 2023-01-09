package logger

import (
    "Anubis/pkg/config"
    "github.com/sirupsen/logrus"
    "os"
)

var logger = logrus.New()
var logLevels = map[string]logrus.Level{
    "debug": logrus.DebugLevel,
    "ingo":  logrus.InfoLevel,
    "warn":  logrus.WarnLevel,
    "error": logrus.ErrorLevel,
}

func Initial() {
    formatter := &Formatter{
        LogFormat:       "%time% [%lvl%] %msg%",
        TimestampFormat: "2006-01-02 15:04:05",
    }

    conf := config.GlobalConfig
    level, ok := logLevels[conf.LogLevel]
    if !ok {
        level = logrus.InfoLevel
    }

    logger.SetFormatter(formatter)
    logger.SetLevel(level)
    logger.SetOutput(os.Stdout)
}

func Debug(args ...interface{}) {
    logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
    logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
    logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
    logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
    logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
    logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
    logger.Errorf(format, args...)
}

func Panic(args ...interface{}) {
    logrus.Panic(args...)
}

func Fatal(args ...interface{}) {
    logrus.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
    logrus.Fatalf(format, args...)
}
