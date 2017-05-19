package core

import (
	"errors"
	"github.com/sydnash/lotou/vector"
	"sync"
)

//definition for node id
//if system is standalone, then it's node id is DEFAULT_NODE_ID
//if system is multi node, master's node id is MASTER_NODE_ID, slave's node is allocation by master service.
const (
	NODE_ID_OFF                  = 64 - 16
	NODE_ID_MASK                 = 0xFFFF << NODE_ID_OFF
	INVALID_SERVICE_ID           = NODE_ID_MASK
	DEFAULT_NODE_ID              = 0XFFFF
	MASTER_NODE_ID               = 0
	INIT_SERVICE_ID    ServiceID = 10
)

type handleDic map[uint64]*service

//a storage that stores all local services
type handleStorage struct {
	dicMutex           sync.Mutex
	dic                handleDic
	nodeId             uint64
	curId              uint64
	baseServiceIdCache *vector.Vector
}

var (
	h                   *handleStorage
	ServiceNotFindError = errors.New("service is not find.")
	exitGroup           sync.WaitGroup
)

func newHandleStorage() *handleStorage {
	h := &handleStorage{}
	h.nodeId = DEFAULT_NODE_ID
	h.dic = make(map[uint64]*service)
	h.curId = uint64(INIT_SERVICE_ID)
	h.baseServiceIdCache = vector.NewCap(1000)
	return h
}

//checkIsLocalId checks a given service id is a local service's id
//a serviceId's node id is equal to DEFAULT_NODE_ID or nodeId is a local service's id
func checkIsLocalId(id ServiceID) bool {
	nodeId := id.parseNodeId()
	if nodeId == DEFAULT_NODE_ID {
		return true
	}
	if nodeId == h.nodeId {
		return true
	}
	return false
}

//checkIsLocalName checks a given name is a local name.
//a name start with '.' or empty is a local name. others a all global name
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

//registerService register a service and allocate a service id to the given service.
func registerService(s *service) ServiceID {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	var baseServiceId uint64
	//if h.baseServiceIdCache.Empty() {
	h.curId++
	baseServiceId = h.curId
	//} else {
	//	baseServiceId = h.baseServiceIdCache.Pop().(uint64)
	//}
	id := h.nodeId<<NODE_ID_OFF | baseServiceId
	h.dic[id] = s
	sid := ServiceID(id)
	s.setId(sid)
	exitGroup.Add(1)
	return ServiceID(sid)
}

//unregisterService delete a service and put it's to cache which can be resued again when register
func unregisterService(s *service) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	id := uint64(s.getId())
	if _, ok := h.dic[id]; !ok {
		return
	}
	delete(h.dic, id)
	h.baseServiceIdCache.Push((ServiceID(id)).parseBaseId())
	exitGroup.Done()
}

//findServiceById return a service by service id
func findServiceById(id ServiceID) (s *service, err error) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	s, ok := h.dic[uint64(id)]
	if !ok {
		err = ServiceNotFindError
	}
	return s, err
}

//findServiceByName return a service by service name, it only return local service.
func findServiceByName(name string) (s *service, err error) {
	PanicWhen(len(name) == 0, "name must not empty.")
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
