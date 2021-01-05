package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger.Fatalf("test: %s", "Fatalf()")
	logger.Errorf("test: %s", "Errorf()")
	logger.Warnf("test: %s", "Warnf()")
	logger.Infof("test: %s", "Infof()")
	logger.Tracef("test: %s", "Tracef()")
}
