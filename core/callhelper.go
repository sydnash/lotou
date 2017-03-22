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
	funcMap   map[string]*callbackDesc
	idFuncMap map[int]*callbackDesc
}

type ReplyFunc func(data ...interface{})

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper() *CallHelper {
	ret := &CallHelper{}
	ret.funcMap = make(map[string]*callbackDesc)
	ret.idFuncMap = make(map[int]*callbackDesc)
	return ret
}

func (c *CallHelper) AddFunc(name string, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[name] = &callbackDesc{f, true}
}

func (c *CallHelper) AddMethod(name string, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type.")
	c.funcMap[name] = &callbackDesc{f, true}
}

func (c *CallHelper) AddFuncInt(id int, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.idFuncMap[id] = &callbackDesc{f, true}
}

func (c *CallHelper) AddMethodInt(id int, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type")
	c.idFuncMap[id] = &callbackDesc{f, true}
}

func (c *CallHelper) setIsAutoReply(id interface{}, isAutoReply bool) {
	cb := c.findCallbackDesk(id)
	cb.isAutoReply = isAutoReply
	if !isAutoReply {
		t := reflect.New(cb.cb.Type().In(ReplyFuncPosition))
		log.Debug("%v", t.Elem().Interface().(ReplyFunc))
	}
}

func (c *CallHelper) getIsAutoReply(id interface{}) bool {
	return c.findCallbackDesk(id).isAutoReply
}

func (c *CallHelper) findCallbackDesk(id interface{}) *callbackDesc {
	var cb *callbackDesc
	var ok bool
	switch key := id.(type) {
	case int:
		cb, ok = c.idFuncMap[key]
	case string:
		cb, ok = c.funcMap[key]
	default:
		log.Fatal("methodid: %v is not registered; %v", id)
	}
	if !ok {
		log.Fatal("func: %v is not found", id)
	}
	return cb
}

func (c *CallHelper) Call(id interface{}, src ServiceID, param ...interface{}) []interface{} {
	cb := c.findCallbackDesk(id)
	p := []reflect.Value{}
	p = append(p, reflect.ValueOf(src)) //append src service id
	for _, v := range param {
		p = append(p, reflect.ValueOf(v))
	}
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", id, err)
		}
	}()
	ret := cb.cb.Call(p)

	out := make([]interface{}, len(ret))
	for i, v := range ret {
		out[i] = v.Interface()
	}
	return out
}

func (c *CallHelper) CallWithReplyFunc(id interface{}, src ServiceID, replyFunc ReplyFunc, param ...interface{}) {
	cb := c.findCallbackDesk(id)
	p := []reflect.Value{}
	p = append(p, reflect.ValueOf(src)) //append src service id
	p = append(p, reflect.ValueOf(replyFunc))
	for _, v := range param {
		p = append(p, reflect.ValueOf(v))
	}
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("CallHelper.Call err: method: %v %v", id, err)
		}
	}()
	cb.cb.Call(p)
}
