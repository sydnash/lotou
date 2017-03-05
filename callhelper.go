package lotou

import (
	"errors"
	"github.com/sydnash/lotou/core"
	"reflect"
)

//CallHelper help to call functions where the comein params is like interface{} or []interface{}
//avoid to use type assert(a.(int))
//it's not thread safe
type CallHelper struct {
	funcMap map[string]reflect.Value
}

var (
	FuncNotFound = errors.New("func not found.")
)

func NewCallHelper() *CallHelper {
	ret := &CallHelper{}
	ret.funcMap = make(map[string]reflect.Value)
	return ret
}

func (c CallHelper) AddFunc(name string, f reflect.Value) {
	core.PanicWhen(f.Kind() != reflect.Func)
	c.funcMap[name] = f
}

func (c CallHelper) Call(name string, param ...interface{}) []reflect.Value {
	f, ok := c.funcMap[name]
	if !ok {
		return []reflect.Value{reflect.ValueOf(FuncNotFound)}
	}
	len := len(param)
	p := make([]reflect.Value, len, len)
	for i, v := range param {
		p[i] = reflect.ValueOf(v)
	}
	return f.Call(p)
}
