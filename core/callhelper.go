package core

import (
	"errors"
	"github.com/sydnash/lotou/log"
	"reflect"
)

//CallHelper help to call functions where the comein params is like interface{} or []interface{}
//avoid to use type assert(a.(int))
//it's not thread safe
type CallHelper struct {
	funcMap   map[string]reflect.Value
	idFuncMap map[int]reflect.Value
}

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper() *CallHelper {
	ret := &CallHelper{}
	ret.funcMap = make(map[string]reflect.Value)
	ret.idFuncMap = make(map[int]reflect.Value)
	return ret
}

func (c *CallHelper) AddFunc(name string, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[name] = f
}

func (c *CallHelper) AddMethod(name string, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type.")
	c.funcMap[name] = f
}

func (c *CallHelper) AddFuncInt(id int, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.idFuncMap[id] = f
}

func (c *CallHelper) AddMethodInt(id int, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type")
	c.idFuncMap[id] = f
}

func (c *CallHelper) Call(id interface{}, param ...interface{}) []interface{} {
	var f reflect.Value
	var ok bool
	switch key := id.(type) {
	case int:
		f, ok = c.idFuncMap[key]
	case string:
		f, ok = c.funcMap[key]
	default:
		panic("methodid is not register")
	}
	if !ok {
		log.Fatal("func: %v is not found", id)
	}
	paramLen := len(param)
	p := make([]reflect.Value, paramLen, paramLen)
	for i, v := range param {
		p[i] = reflect.ValueOf(v)
	}
	ret := f.Call(p)

	out := make([]interface{}, len(ret))
	for i, v := range ret {
		out[i] = v.Interface()
	}
	return out
}
