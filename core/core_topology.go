package core

var (
	isMaster         bool = true
	registerNodeChan chan uint
	selfNodeId       uint
)

func init() {
	registerNodeChan = make(chan uint)
	selfNodeId = NATIVE_PRE
}

func RegisterNode() {
	if !isMaster {
		if selfNodeId != NATIVE_PRE {
			return
		}
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if selfNodeId != NATIVE_PRE {
			return
		}
		core.SendName(".slave", NATIVE_CORE_ID, "registerNode")
		selfNodeId = <-registerNodeChan
	}
}

func TopologyMSG() {
}
