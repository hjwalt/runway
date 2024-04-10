package logger_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func TestAllFunctions(t *testing.T) {
	logger.Debug("test", zap.String("test", "test"))
	logger.Info("test", zap.String("test", "test"))
	logger.Warn("test", zap.String("test", "test"))
	logger.Error("test", zap.String("test", "test"))

	logger.Debugf("test %s %d", "test", 0)
	logger.Infof("test %s %d", "test", 0)
	logger.Warnf("test %s %d", "test", 0)
	logger.Errorf("test %s %d", "test", 0)

	logger.WarnErr("test", errors.New("test error"))
	logger.ErrorErr("test", errors.New("test error"))

	logger.InfoIfTrue(true, "test", zap.String("test", "test"))
	logger.InfoIfTrue(false, "test", zap.String("test", "test"))
	logger.WarnIfTrue(true, "test", zap.String("test", "test"))
	logger.WarnIfTrue(false, "test", zap.String("test", "test"))
	logger.ErrorIfErr("test", errors.New("test error"))
	logger.ErrorIfErr("test", nil)

	slog.Info("test", "key", "value")
	slog.Warn("test", "key", "value")

	localLogger := slog.With("attr", "attr-val").WithGroup("test").With("groupattr", "group-val")
	localLogger.Info("test", "key", "value")
	localLogger.Warn("test", "key", "value")

	// assert.Fail(t, "test")
}

func TestProductionSettings(t *testing.T) {
	logger.DefaultZap(true, "test.log", zap.DebugLevel)

	logger.Debug("test", zap.String("test", "test"))
	logger.Info("test", zap.String("test", "test"))
	logger.Warn("test", zap.String("test", "test"))
	logger.Error("test", zap.String("test", "test"))

	logger.Debugf("test %s %d", "test", 0)
	logger.Infof("test %s %d", "test", 0)
	logger.Warnf("test %s %d", "test", 0)
	logger.Errorf("test %s %d", "test", 0)

	logger.WarnErr("test", errors.New("test error"))
	logger.ErrorErr("test", errors.New("test error"))

	slog.Info("test", "key", "value")

	localLogger := slog.With("attr", "attr-val").WithGroup("test").With("groupattr", "group-val")
	localLogger.Info("test", "key", "value")
}
