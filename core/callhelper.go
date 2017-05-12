package core

import (
	"errors"
	"fmt"
	"github.com/sydnash/lotou/log"
	"reflect"
)

//CallHelper use reflect.Call to invoke a function.
//it's not thread safe
const (
	ReplyFuncPosition = 1
)

type callbackDesc struct {
	cb          reflect.Value
	isAutoReply bool
}
type CallHelper struct {
	funcMap         map[CmdType]*callbackDesc
	hostServiceName string //help to locate which callback is not registered.
}

type ReplyFunc func(data ...interface{})

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper(name string) *CallHelper {
	return &CallHelper{
		hostServiceName: name,
		funcMap:         make(map[CmdType]*callbackDesc),
	}
}

//AddFunc add callback with normal function
func (c *CallHelper) AddFunc(cmd CmdType, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[cmd] = &callbackDesc{f, true}
}

//AddMethod add callback with struct's method by method name
//method name muse be exported
func (c *CallHelper) AddMethod(cmd CmdType, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, fmt.Sprintf("[CallHelper:AddMethod] cmd{%v} method must be a function type.", cmd))
	c.funcMap[cmd] = &callbackDesc{f, true}
}

//setIsAutoReply recode special cmd is auto reply after Call is return
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
		if cb, ok = c.funcMap[Cmd_Default]; ok {
			log.Info("func: <%v>:%d is not found in {%v}, use default cmd handler.", cmd, len(cmd), c.hostServiceName)
		} else {
			log.Fatal("func: <%v>:%d is not found in {%v}", cmd, len(cmd), c.hostServiceName)
		}
	}
	return cb
}

//Call invoke special function for cmd
func (c *CallHelper) Call(cmd CmdType, src ServiceID, param ...interface{}) []interface{} {
	cb := c.findCallbackDesc(cmd)
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", cmd, err)
		}
	}()

	//addition one param for source service id
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

//CallWithReplyFunc invoke special function for cmd with a reply function which is used to reply Call or Request.
func (c *CallHelper) CallWithReplyFunc(cmd CmdType, src ServiceID, replyFunc ReplyFunc, param ...interface{}) {
	cb := c.findCallbackDesc(cmd)
	//addition two param for source service id and reply function
	p := make([]reflect.Value, len(param)+2)
	p[0] = reflect.ValueOf(src)
	p[1] = reflect.ValueOf(replyFunc)

	HelperFunctionToUseReflectCall(cb.cb, p, 2, param)
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", cmd, err)
		}
	}()
	cb.cb.Call(p)
}
