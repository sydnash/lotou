package core

import (
	"errors"
	"sync"
)

//definition for node id
const (
	NODE_ID_MASK       = 0xFF000000
	NODE_ID_OFF        = 24
	INVALID_SERVICE_ID = (0XFF << NODE_ID_OFF) & NODE_ID_MASK
	DEFAULT_NODE_ID    = 0XFF
)

type handleDic map[uint]Service

type handleStorage struct {
	dicMutex sync.Mutex
	dic      handleDic
	nodeId   uint
	curId    uint
}

var (
	h                   *handleStorage
	ServiceNotFindError = errors.New("Service is not find.")
)

func newHandleStorage() *handleStorage {
	h := &handleStorage{}
	h.nodeId = DEFAULT_NODE_ID
	h.dic = make(map[uint]Service)
	h.curId = 0
}

func parseNodeIdFromId(id uint) {
	return (id & NODE_ID_MASK) >> NODE_ID_OFF
}

func checkIsLocalId(id uint) {
	nodeId := ParseNodeId(id)
	if nodeId == NODE_ID_MASK {
		return true
	}
	if nodeId == H.nodeId {
		return true
	}
	return false
}

func checkIsLocalName(name string) {
	PanicWhen(len(name) == 0)
	if name[0] == '.' {
		return true
	}
	return false
}

func init() {
	h = newHandleStorage()
}

func registerService(s Service) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	h.curId++
	id := h.nodeId<<NODE_ID_OFF | h.curId
	h.dic[h.curId] = s
}

func findServiceById(id uint) (s Service, err error) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	s, ok := h.dic[id]
	if !ok {
		err = ServiceNotFindError
	}
	return s, err
}

func findServiceByName(name string) (s Service, err error) {
	PanicWhen(len(name) == 0)
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for key, value := range h.dic {
		if key == name {
			s = value
			return s, nil
		}
	}
	return nil, ServiceNotFindError
}
