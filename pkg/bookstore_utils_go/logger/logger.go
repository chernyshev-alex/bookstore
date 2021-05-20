package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

var log Logger

type BookstoreLogger interface {
	Printf(format string, v ...interface{})
}

type Logger struct {
	wrapped *zap.Logger
}

func init() {

	logConfig := zap.Config{
		OutputPaths: []string{getLogOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.wrapped, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func getLevel() zapcore.Level {
	switch strings.TrimSpace(os.Getenv(envLogLevel)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	}
	return zap.InfoLevel
}

func getLogOutput() (logPath string) {
	if logPath := strings.TrimSpace(os.Getenv(envLogOutput)); logPath == "" {
		return "stdout"
	}
	return
}

func GetLogger() BookstoreLogger {
	return log
}

// elastic compatible interface
func (l Logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func (l Logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

func Info(msg string, tags ...zap.Field) {
	log.wrapped.Info(msg, tags...)
	log.wrapped.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.wrapped.Error(msg, tags...)
	log.wrapped.Sync()

}
