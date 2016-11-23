package core

import (
	"fmt"
	"github.com/sydnash/lotou/log"
	"sync"
)

type nameRet struct {
	id   uint
	ok   bool
	name string
}

var (
	isMaster         bool = false
	registerNodeChan chan uint
	selfNodeId       uint
	once             sync.Once
	nameMapMutex     sync.Mutex
	nameChanMap      map[uint]chan *nameRet
	nameRequestId    uint
)

func init() {
	registerNodeChan = make(chan uint)
	selfNodeId = NATIVE_PRE
	nameChanMap = make(map[uint]chan *nameRet)
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

func globalName(id uint, name string) bool {
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

func getGlobalIdByName(name string) (uint, bool) {
	if isMaster || CheckIsLocalName(name) {
		log.Warn(fmt.Sprintf("GetIdByName: %s is not exist", name))
		return 0, false
	}
	nameMapMutex.Lock()
	nameRequestId++
	tmp := nameRequestId
	nameMapMutex.Unlock()

	SendName(".slave", NATIVE_CORE_ID, "getIdByName", name, tmp)
	ch := make(chan *nameRet)

	nameMapMutex.Lock()
	nameChanMap[tmp] = ch
	nameMapMutex.Unlock()
	ret := <-ch
	return ret.id, ret.ok
}
func GetIdByNameRet(id uint, ok bool, name string, rid uint) {
	nameMapMutex.Lock()
	ch := nameChanMap[rid]
	delete(nameChanMap, rid)
	nameMapMutex.Unlock()
	ch <- &nameRet{id, ok, name}
}

func CheckIsLocalName(name string) bool {
	if len(name) == 0 {
		panic("name is empty.")
	}
	if name[0] == '.' {
		return true
	}
	return false
}
func ParseNodeId(id uint) uint {
	ret := id >> 24
	return ret
}
func CheckIsLocalServiceId(id uint) bool {
	nodeId := ParseNodeId(id)
	if nodeId == selfNodeId || nodeId == 0xFF {
		return true
	}
	if isMaster && nodeId == 0x0 {
		return true
	}
	return false
}
func sendToMaster(m *Message) {
	if isMaster {
		SendName(".master", NATIVE_CORE_ID, "forward", m)
	} else {
		SendName(".slave", NATIVE_CORE_ID, "forward", m)
	}
}
func NameDic() map[string]uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.nameDic
}

func TopologyMSG() {
}
