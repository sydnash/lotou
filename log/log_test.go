package log_test

import (
	"github.com/sydnash/majiang/log"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log.Init("test", log.DEBUG_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	for {
		time.Sleep(time.Second)
		log.Debug("hahaha %v, %v", 2, 3)
	}
}
