package core

import (
	"errors"
	"github.com/sydnash/lotou/log"
	"reflect"
)

//callHelper help to call functions where the comein params is like interface{} or []interface{}
//avoid to use type assert(a.(int))
//it's not thread safe
type callHelper struct {
	funcMap   map[string]reflect.Value
	idFuncMap map[int]reflect.Value
}

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper() *callHelper {
	ret := &callHelper{}
	ret.funcMap = make(map[string]reflect.Value)
	ret.idFuncMap = make(map[int]reflect.Value)
	return ret
}

func (c *callHelper) AddFunc(name string, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[name] = f
}

func (c *callHelper) AddMethod(name string, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type.")
	c.funcMap[name] = f
}

func (c *callHelper) AddFuncInt(id int, fun interface{}) {
	f := reflect.ValueOf(fun)
	PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.idFuncMap[id] = f
}

func (c *callHelper) AddMethodInt(id int, v interface{}, methodName string) {
	self := reflect.ValueOf(v)
	f := self.MethodByName(methodName)
	PanicWhen(f.Kind() != reflect.Func, "method must be a function type")
	c.idFuncMap[id] = f
}

func (c *callHelper) Call(id interface{}, src ServiceID, encType int32, param ...interface{}) []interface{} {
	var cb reflect.Value
	var ok bool
	switch key := id.(type) {
	case int:
		cb, ok = c.idFuncMap[key]
	case string:
		cb, ok = c.funcMap[key]
	default:
		log.Fatal("methodid: %v is not register", id)
	}
	if !ok {
		log.Fatal("func: %v is not found", id)
	}
	p := []reflect.Value{}
	p = append(p, reflect.ValueOf(src)) //append src service id
	for _, v := range param {
		p = append(p, reflect.ValueOf(v))
	}
	ret := cb.Call(p)

	out := make([]interface{}, len(ret))
	for i, v := range ret {
		out[i] = v.Interface()
	}
	return out
}
