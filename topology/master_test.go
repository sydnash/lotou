package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
)

func TestMaster(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	core.SetAsMaster()
	topology.StartMaster("127.0.0.1", "4000")

	log.Info("test")
	for {
	}
}
