package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	log, err := zap.NewDevelopment(zap.AddStacktrace(zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l == zapcore.PanicLevel })))

	if err != nil {
		panic(err)
	}

	restoreStdLog := zap.RedirectStdLog(log)
	defer restoreStdLog()

	logger = log
}

// Info ...
func Info(msg string, fields ...zapcore.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Debug ...
func Debug(msg string, fields ...zapcore.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Error ...
func Error(msg string, fields ...zapcore.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}
