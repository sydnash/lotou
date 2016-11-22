package core

import (
	"github.com/sydnash/lotou/log"
	"sync"
)

var (
	isMaster         bool = false
	registerNodeChan chan uint
	selfNodeId       uint
	once             sync.Once
)

func init() {
	registerNodeChan = make(chan uint)
	selfNodeId = NATIVE_PRE
}

func RegisterNode() {
	once.Do(func() {
		if !isMaster {
			if selfNodeId != NATIVE_PRE {
				return
			}
			SendName(".slave", NATIVE_CORE_ID, "registerNode")
			selfNodeId = <-registerNodeChan
			log.Info("registerNode success: %d", selfNodeId)
		}
	})
}

func GlobalName(id uint, name string) bool {
	if isMaster {
		//sync name
		SendName(".master", NATIVE_CORE_ID, "syncName", id, name)
	} else {
		//name to master
		SendName(".slave", NATIVE_CORE_ID, "registerName", id, name)
	}
	return true
}
func SyncName(serviceId uint, serviceName string) {
	c.nameMutex.Lock()
	defer c.nameMutex.Unlock()
	log.Info("sync name: %s, id:%d", serviceName, serviceId)
	c.nameDic[serviceName] = serviceId
}

func RegisterNodeRet(id uint) {
	registerNodeChan <- id
}
func SetAsMaster() {
	selfNodeId = 0
	isMaster = true
}

func GenerateNodeId() uint {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.nodeId++
	return c.nodeId
}
func TopologyMSG() {
}
