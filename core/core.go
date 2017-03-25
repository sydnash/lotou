package core

import (
	"fmt"
	"github.com/sydnash/lotou/log"
	"reflect"
	"runtime/debug"
)

//StartService start a given module with specific name
//call module's OnInit interface after register
//and register name to master if name is a global name
//and start msg loop in an another goroutine
func StartService(name string, m Module) ServiceID {
	s := newService(name)
	s.m = m
	id := registerService(s)
	m.setService(s)
	d := m.getDuration()
	if !checkIsLocalName(name) {
		globalName(id, name)
	}
	if d > 0 {
		s.runWithLoop(d)
	} else {
		s.run()
	}
	return id
}

//HelperFunctionToUseReflectCall help to convert realparam like([]interface{}) to reflect.Call's param([]reflect.Value
//and if param is nil, then use reflect.New to create a empty to avoid crash when reflect.Call invoke.
//and genrate more readable error message when param is not ok.
func HelperFunctionToUseReflectCall(f reflect.Value, callParam []reflect.Value, startNum int, realParam []interface{}) {
	n := len(realParam)
	lastCallParamIdx := f.Type().NumIn() - 1
	isVariadic := f.Type().IsVariadic()
	for i := 0; i < n; i++ {
		paramIndex := i + startNum
		var expectedType reflect.Type
		if isVariadic && paramIndex >= lastCallParamIdx { //variadic function's last param is []T
			expectedType = f.Type().In(lastCallParamIdx)
			expectedType = expectedType.Elem()
		} else {
			expectedType = f.Type().In(paramIndex)
		}
		//if param is nil, create a empty reflect.Value
		if realParam[i] == nil {
			callParam[paramIndex] = reflect.New(expectedType).Elem()
		} else {
			callParam[paramIndex] = reflect.ValueOf(realParam[i])
		}
		actualType := callParam[paramIndex].Type()
		if !actualType.AssignableTo(expectedType) {
			//panic if param is not assignable to Call
			errStr := fmt.Sprintf("InvocationCausedPanic: called with a mismatched parameter type [parameter #%v: expected %v; got %v].", paramIndex, expectedType, actualType)
			panic(errStr)
		}
	}
}

//Parse Node Id parse node id from service id
func ParseNodeId(id ServiceID) uint64 {
	return id.parseNodeId()
}

//Send send a message to dst service no src service.
func Send(dst ServiceID, msgType MsgType, encType EncType, cmd CmdType, data ...interface{}) error {
	return lowLevelSend(INVALID_SERVICE_ID, dst, msgType, encType, 0, cmd, data...)
}

//SendCloseToAll simple send a close msg to all service
func SendCloseToAll() {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for _, ser := range h.dic {
		localSendWithoutMutex(INVALID_SERVICE_ID, ser, MSG_TYPE_CLOSE, MSG_ENC_TYPE_NO, 0, Cmd_None, false)
	}
}

//Wait wait on a sync.WaitGroup, until all service is closed.
func Wait() {
	exitGroup.Wait()
}

//CheckIsLocalServiceId heck a given service id is a local service
func CheckIsLocalServiceId(id ServiceID) bool {
	return checkIsLocalId(id)
}

//SafeGo start a groutine, and handle all panic within it.
func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("recover: stack: %v\n, %v", string(debug.Stack()), err)
			}
			return
		}()
		f()
	}()
}
