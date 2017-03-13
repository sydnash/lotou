package core

import (
	"errors"
	"sync"
)

//definition for node id
const (
	NODE_ID_MASK                 = 0xFF000000
	NODE_ID_OFF                  = 24
	INVALID_SERVICE_ID           = (0XFF << NODE_ID_OFF) & NODE_ID_MASK
	DEFAULT_NODE_ID              = 0XFF
	MASTER_NODE_ID               = 0
	INIT_SERVICE_ID    ServiceID = 10
)

type handleDic map[uint]*service

type handleStorage struct {
	dicMutex sync.Mutex
	dic      handleDic
	nodeId   uint
	curId    uint
}

var (
	h                   *handleStorage
	ServiceNotFindError = errors.New("service is not find.")
	exitGroup           sync.WaitGroup
)

func newHandleStorage() *handleStorage {
	h := &handleStorage{}
	h.nodeId = DEFAULT_NODE_ID
	h.dic = make(map[uint]*service)
	h.curId = uint(INIT_SERVICE_ID)
	return h
}

func parseNodeIdFromId(id ServiceID) uint {
	return (uint(id) & NODE_ID_MASK) >> NODE_ID_OFF
}

func checkIsLocalId(id ServiceID) bool {
	nodeId := parseNodeIdFromId(id)
	if nodeId == DEFAULT_NODE_ID {
		return true
	}
	if nodeId == h.nodeId {
		return true
	}
	return false
}

func checkIsLocalName(name string) bool {
	if len(name) == 0 {
		return true
	}
	if name[0] == '.' {
		return true
	}
	return false
}

func init() {
	h = newHandleStorage()
}

func registerService(s *service) ServiceID {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	h.curId++
	id := h.nodeId<<NODE_ID_OFF | h.curId
	h.dic[id] = s
	sid := ServiceID(id)
	s.setId(sid)
	exitGroup.Add(1)
	return ServiceID(sid)
}

func unregisterService(s *service) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	id := uint(s.getId())
	if _, ok := h.dic[id]; !ok {
		return
	}
	delete(h.dic, id)
	exitGroup.Done()
}

func findServiceById(id ServiceID) (s *service, err error) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	s, ok := h.dic[uint(id)]
	if !ok {
		err = ServiceNotFindError
	}
	return s, err
}

func findServiceByName(name string) (s *service, err error) {
	PanicWhen(len(name) == 0)
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for _, value := range h.dic {
		if value.getName() == name {
			s = value
			return s, nil
		}
	}
	return nil, ServiceNotFindError
}
