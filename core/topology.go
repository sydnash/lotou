package core

import (
	"sync"
)

type nameRet struct {
	id   ServiceID
	ok   bool
	name string
}

var (
	once                   sync.Once
	isStandalone, isMaster bool
	registerNodeChan       chan uint64
	nameChanMap            map[uint]chan *nameRet
	nameMapMutex           sync.Mutex
	nameRequestId          uint
	beginNodeId            uint64
)

func init() {
	registerNodeChan = make(chan uint64)

	nameChanMap = make(map[uint]chan *nameRet)
	nameRequestId = 0

	isStandalone = true

	beginNodeId = 10
}

func InitNode(_isStandalone, _isMaster bool) {
	isStandalone = _isStandalone
	isMaster = _isMaster
	if !isStandalone && isMaster {
		h.nodeId = MASTER_NODE_ID
	}
}

//RegisterNode : register slave node to master, and get a node id.
func RegisterNode() {
	once.Do(func() {
		if !isStandalone && !isMaster {
			route("registerNode")
			h.nodeId = <-registerNodeChan
		}
	})
}

func DispatchRegisterNodeRet(id uint64) {
	registerNodeChan <- id
}

//globalName regist name to master
//it will notify all exist service through distribute msg.
func globalName(id ServiceID, name string) {
	route("registerName", uint64(id), name)
}

//route send msg to master
//if node is not a master node, it send to .slave node first, .slave will forward msg to master.
func route(data ...interface{}) {
	if !isStandalone {
		router, err := findServiceByName(".router")
		if err != nil {
			return
		}
		localSendWithoutMutex(INVALID_SERVICE_ID, router, MSG_TYPE_NORMAL, MSG_ENC_TYPE_NO, data...)
	}
}

//NameToId couldn't guarantee get the correct id for name.
//it will return err if the named server is until now.
func NameToId(name string) (ServiceID, error) {
	ser, err := findServiceByName(name)
	if err == nil {
		return ser.getId(), nil
	}
	if !checkIsLocalName(name) {
		nameMapMutex.Lock()
		nameRequestId++
		tmp := nameRequestId
		nameMapMutex.Unlock()

		ch := make(chan *nameRet)
		nameMapMutex.Lock()
		nameChanMap[tmp] = ch
		nameMapMutex.Unlock()

		route("getIdByName", name, tmp)
		ret := <-ch
		close(ch)
		if !ret.ok {
			return INVALID_SERVICE_ID, ServiceNotFindError
		}
		return ret.id, nil
	}
	return INVALID_SERVICE_ID, ServiceNotFindError
}

func DispatchGetIdByNameRet(id ServiceID, ok bool, name string, rid uint) {
	nameMapMutex.Lock()
	ch := nameChanMap[rid]
	delete(nameChanMap, rid)
	nameMapMutex.Unlock()
	ch <- &nameRet{id, ok, name}
}

func GenerateNodeId() uint64 {
	ret := beginNodeId
	beginNodeId++
	return ret
}
