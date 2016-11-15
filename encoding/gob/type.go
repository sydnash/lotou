package gob

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	rtypeToUtMutxt sync.Mutex
	typeId         uint
	rtToUt         map[reflect.Type]*user_type
	typeIdToUt     map[uint]*user_type
	baseRtToTypeId map[reflect.Type]uint
)

type user_type struct {
	base     reflect.Type
	depth    uint
	baseId   uint //base reflect type's id
	encodeId uint //high 8 bit is depth,  low 24 bit is id
}

func registerType(i interface{}) {
	value := reflect.ValueOf(i)
	ut := validType(value)

	id, ok := baseRtToTypeId[ut.base]
	if ok {
		return
	}
	id++
	ut.baseId = id
	ut.encodeId = ut.depth<<24 | ut.baseId
	typeIdToUt[id] = ut
	baseRtToTypeId[ut.base] = id
}

func validType(value reflect.Value) *user_type {
	typ := value.Type()
	ut, ok := rtToUt[typ]
	if ok {
		return ut
	}
	bt, depth := findBaseAndDepth(value)
	ut = &user_type{base: bt, depth: depth}
	rtToUt[typ] = ut
	return ut
}
func findBaseAndDepth(value reflect.Value) (rt reflect.Type, depth uint) {
	for rt = value.Type(); rt.Kind() == reflect.Ptr; {
		value = value.Elem()
		rt = value.Type()
		depth++
	}
	fmt.Println(rt, depth)
	return rt, depth
}

func init() {
	rtToUt = make(map[reflect.Type]*user_type)
	typeIdToUt = make(map[uint]*user_type)
	baseRtToTypeId = make(map[reflect.Type]uint)

	registerType(byte(0))
	registerType(int(0))
	registerType(uint(0))
	registerType(bool(false))

	registerType(int8(0))
	registerType(int16(0))
	registerType(int32(0))
	registerType(int64(0))
	registerType(uint8(0))
	registerType(uint16(0))
	registerType(uint32(0))
	registerType(uint64(0))
}
