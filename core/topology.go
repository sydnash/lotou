package core

import (
	"sync"
)

type nameRet struct {
	id   uint
	ok   bool
	name string
}

var (
	once                   sync.Once
	isStandalone, isMaster bool
	registerNodeChan       chan uint
	nameChanMap            map[uint]chan *nameRet
	nameMapMutex           sync.Mutex
	nameRequestId          uint
	beginNodeId            uint
)

func init() {
	registerNodeChan = make(chan uint)

	nameChanMap = make(map[uint]chan *nameRet)
	nameRequestId = 0

	isStandalone = true

	beginNodeId = 10
}

func InitNode(_isStandalone, _isMaster bool) {
	isStandalone = _isStandalone
	isMaster = _isMaster
	if !isStandalone && isMaster {
		h.nodeId = INIT_NODE_ID
	}
}

//RegisterNode : register slave node to master, and get a node id.
func RegisterNode() {
	once.Do(func() {
		if !isStandalone && !isMaster {
			sendName(INVALID_SERVICE_ID, ".slave", MSG_TYPE_NORMAL, "registerNode")
			h.nodeId = <-registerNodeChan
		}
	})
}

func RegisterNodeRet(id uint) {
	registerNodeChan <- id
}

func globalName(id uint, name string) {
	sendToMaster("registerName", id, name)
}

func sendToMaster(data ...interface{}) {
	if !isStandalone {
		if isMaster {
			//sync name
			sendName(INVALID_SERVICE_ID, ".master", MSG_TYPE_NORMAL, data...)
		} else {
			//name to master
			sendName(INVALID_SERVICE_ID, ".slave", MSG_TYPE_NORMAL, data...)
		}
	}
}

func NameToId(name string) (uint, error) {
	ser, err := findServiceByName(name)
	if err == nil {
		return ser.getId(), nil
	}
	if !checkIsLocalName(name) {
		nameMapMutex.Lock()
		nameRequestId++
		tmp := nameRequestId
		nameMapMutex.Unlock()

		sendToMaster("getIdByName", name, tmp)
		ch := make(chan *nameRet)

		nameMapMutex.Lock()
		nameChanMap[tmp] = ch
		nameMapMutex.Unlock()
		ret := <-ch
		close(ch)
		if !ret.ok {
			return INVALID_SERVICE_ID, ServiceNotFindError
		}
		return ret.id, nil
	}
	return INVALID_SERVICE_ID, ServiceNotFindError
}

func GetIdByNameRet(id uint, ok bool, name string, rid uint) {
	nameMapMutex.Lock()
	ch := nameChanMap[rid]
	delete(nameChanMap, rid)
	nameMapMutex.Unlock()
	ch <- &nameRet{id, ok, name}
}

func GenerateNodeId() uint {
	ret := beginNodeId
	beginNodeId++
	return ret
}
