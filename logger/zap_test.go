package logger_test

import (
	"errors"
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
}
