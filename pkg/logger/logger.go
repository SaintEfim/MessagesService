package logger

import (
	"os"

	"MessagesService/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logInfo *config.Config) *zap.Logger {
	var logger *zap.Logger

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var level zapcore.LevelEnabler
	switch logInfo.Logs.Level {
	case "error":
		level = zapcore.ErrorLevel
	case "info":
		level = zapcore.InfoLevel
	case "debug":
		level = zapcore.DebugLevel
	case "warning":
		level = zapcore.WarnLevel
	default:
		level = zapcore.ErrorLevel
	}

	if logInfo.Logs.Path != "" {
		fileEncoder := zapcore.NewJSONEncoder(cfg)
		logRotation := &lumberjack.Logger{
			Filename:   logInfo.Logs.Path,
			MaxBackups: logInfo.Logs.MaxBackups,
			MaxAge:     logInfo.Logs.MaxAge,
		}

		writer := zapcore.AddSync(logRotation)

		core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, level))
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		return logger
	}

	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
