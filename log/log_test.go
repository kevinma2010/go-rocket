package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger.Errorf("test: %s", "Errorf()")
	logger.Errorln("test: ", "Errorf()")
	logger.Warnf("test: %s", "Warnf()")
	logger.Warnln("test: ", "Warnf()")
	logger.Infof("test: %s", "Infof()")
	logger.Infoln("test: ", "Infof()")
	logger.Tracef("test: %s", "Tracef()")
	logger.Traceln("test: ", "Tracef()")
	logger.Fatalf("test: %s", "Fatalf()")
	logger.Fatalln("test: ", "Fatalf()")
}
