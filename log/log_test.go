package log

import (
	"testing"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	GetLogger().Info("test", zap.String("t", "hi"))
}
