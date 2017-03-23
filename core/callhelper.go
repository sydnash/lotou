package core

import (
	"errors"
	"github.com/sydnash/lotou/log"
	"reflect"
)

//CallHelper help to call functions where the comein params is like interface{} or []interface{}
//avoid to use type assert(a.(int))
//it's not thread safe
const (
	ReplyFuncPosition = 1
)

type callbackDesc struct {
	cb          reflect.Value
	isAutoReply bool
}
type CallHelper struct {
	funcMap map[CmdType]*callbackDesc
}

type ReplyFunc func(data ...interface{})

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper() *CallHelper {
	ret := &CallHelper{}
	ret.funcMap = make(map[CmdType]*callbackDesc)
	return ret
}

func (c *CallHelper) AddFunc(cmd CmdType, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[cmd] = &callbackDesc{f, true}
}

func (c *CallHelper) AddMethod(cmd CmdType, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type.")
	c.funcMap[cmd] = &callbackDesc{f, true}
}

func (c *CallHelper) setIsAutoReply(cmd CmdType, isAutoReply bool) {
	cb := c.findCallbackDesc(cmd)
	cb.isAutoReply = isAutoReply
	if !isAutoReply {
		t := reflect.New(cb.cb.Type().In(ReplyFuncPosition))
		_ = t.Elem().Interface().(ReplyFunc)
	}
}

func (c *CallHelper) getIsAutoReply(cmd CmdType) bool {
	return c.findCallbackDesc(cmd).isAutoReply
}

func (c *CallHelper) findCallbackDesc(cmd CmdType) *callbackDesc {
	cb, ok := c.funcMap[cmd]
	if !ok {
		log.Fatal("func: %v is not found", cmd)
	}
	return cb
}

func (c *CallHelper) Call(cmd CmdType, src ServiceID, param ...interface{}) []interface{} {
	cb := c.findCallbackDesc(cmd)
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", cmd, err)
		}
	}()

	p := make([]reflect.Value, len(param)+1)
	p[0] = reflect.ValueOf(src) //append src service id
	HelperFunctionToUseReflectCall(cb.cb, p, 1, param)

	ret := cb.cb.Call(p)

	out := make([]interface{}, len(ret))
	for i, v := range ret {
		out[i] = v.Interface()
	}
	return out
}

func (c *CallHelper) CallWithReplyFunc(cmd CmdType, src ServiceID, replyFunc ReplyFunc, param ...interface{}) {
	cb := c.findCallbackDesc(cmd)
	p := make([]reflect.Value, len(param)+2)
	p[0] = reflect.ValueOf(src) //append src service id
	p[1] = reflect.ValueOf(replyFunc)

	HelperFunctionToUseReflectCall(cb.cb, p, 2, param)
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", cmd, err)
		}
	}()
	cb.cb.Call(p)
}
