package topology_test

/*
import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
	"time"
)

func TestSlaveb(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	log.Debug("start slave")
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	core.Name(1, "我是服务1")

	game := &Game{core.NewBase()}
	core.RegisterService(game)
	core.Name(game.Id(), "game1")
	game.SetDispatcher(game)
	go func() {
		for msg := range game.In() {
			game.DispatchM(msg)
		}
	}()

	for {
		time.Sleep(time.Minute * 10)
	}

}

*/
