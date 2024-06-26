package logger

import (
	"log/slog"

	"github.com/hjwalt/runway/trusted"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

var globalSugar *zap.SugaredLogger

func init() {
	Default()
}

func Default() {
	DefaultZap(false, "", zap.DebugLevel)
}

func UseZap(logger *zap.Logger) {
	globalLogger = logger
	globalSugar = globalLogger.Sugar()
}

func DefaultZap(isProduction bool, logFile string, level zapcore.Level) {
	var err error

	// set timestamp layout
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")

	config := zap.Config{}

	// set default encoders
	if isProduction {
		config = zap.NewProductionConfig()
		enccoderConfig := zap.NewProductionEncoderConfig()
		enccoderConfig.StacktraceKey = "" // to hide stacktrace info
		config.EncoderConfig = enccoderConfig
	} else {
		config = zap.NewDevelopmentConfig()
		enccoderConfig := zap.NewDevelopmentEncoderConfig()
		enccoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig = enccoderConfig
	}

	// set log file
	if len(logFile) > 0 {
		config.OutputPaths = []string{logFile}
	}

	// set level
	config.Level.SetLevel(level)

	// skip one level from the helper function
	logger, err := config.Build(zap.AddCallerSkip(1))
	logger = trusted.Must(logger, err)

	slogger, err := config.Build(zap.AddCallerSkip(3))
	logger = trusted.Must(logger, err)

	slog.SetDefault(slog.New(&ZapHandler{config: config, logger: slogger}))

	UseZap(logger)
}

func Debug(message string, fields ...zap.Field) {
	globalLogger.Debug(message, fields...)
}

func Debugf(message string, fields ...interface{}) {
	globalSugar.Debugf(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	globalLogger.Info(message, fields...)
}

func Infof(message string, fields ...interface{}) {
	globalSugar.Infof(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	globalLogger.Warn(message, fields...)
}

func Warnf(message string, fields ...interface{}) {
	globalSugar.Warnf(message, fields...)
}

func WarnErr(message string, err error) {
	globalLogger.Warn(message, zap.Error(err))
}

func Error(message string, fields ...zap.Field) {
	globalLogger.Error(message, fields...)
}

func Errorf(message string, fields ...interface{}) {
	globalSugar.Errorf(message, fields...)
}

func ErrorErr(message string, err error) {
	globalLogger.Error(message, zap.Error(err))
}

// Conditional logging

func InfoIfTrue(condition bool, message string, fields ...zap.Field) {
	if condition {
		globalLogger.Info(message, fields...)
	}
}

func WarnIfTrue(condition bool, message string, fields ...zap.Field) {
	if condition {
		globalLogger.Warn(message, fields...)
	}
}

func ErrorIfErr(message string, err error) {
	if err != nil {
		globalLogger.Error(message, zap.Error(err))
	}
}
