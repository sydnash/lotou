package lotou

import (
	"github.com/sydnash/lotou/core"
)

func Start(isStandalone, isMaster bool) {
	core.Init(isStandalone, isMaster)
	if !isStandalone {
	}
}
