package gob

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	rtypeToUtMutxt    sync.Mutex
	typeIdToUt        map[uint]reflect.Type
	baseRtToTypeId    map[reflect.Type]uint
	kindToReflectType map[reflect.Kind]reflect.Type
)

func RegisterStructType(i interface{}) {
	value := reflect.ValueOf(i)
	rt, depth := findBaseAndDepth(value.Type())
	_ = depth
	if rt.Kind() != reflect.Struct {
		return
	}
	_, ok := baseRtToTypeId[rt]
	if ok {
		return
	}
	typeId := getStructID(rt)
	baseRtToTypeId[rt] = typeId
	typeIdToUt[typeId] = rt
}

func findBaseAndDepth(typ reflect.Type) (rt reflect.Type, depth uint) {
	for rt = typ; rt.Kind() == reflect.Ptr; {
		rt = rt.Elem()
		depth++
	}
	return rt, depth
}

func valueToId(typ reflect.Type) uint16 {
	rt, depth := findBaseAndDepth(typ)
	kind := rt.Kind()
	var structId uint
	var ok bool
	if kind == reflect.Struct {
		structId, ok = baseRtToTypeId[rt]
		if !ok {
			panic(fmt.Sprintf("%v is not register.", rt))
		}
	}
	id := uint16(structId)<<8 | uint16(depth)<<5 | uint16(kind)
	return id
}

func gernerateId(kind, depth, structId uint) uint16 {
	id := uint16(structId)<<8 | uint16(depth)<<5 | uint16(kind)
	return id
}

func parseTypeId(id uint16) (kind uint, depth uint, structId uint) {
	kind = uint(id & 0x1F)
	depth = uint((id >> 5) & 0x07)
	structId = uint((id >> 8) & 0xFF)
	return
}

type TypeDesc struct {
	kind, depth, structId, count uint
	key, elem                    reflect.Type
	hasDepth                     bool
}

func createType(desc *TypeDesc) (typ reflect.Type) {
	var ok bool
	typ, ok = kindToReflectType[reflect.Kind(desc.kind)]
	if !ok {
		switch reflect.Kind(desc.kind) {
		case reflect.Struct:
			id := gernerateId(desc.kind, desc.depth, desc.structId)
			typ = typeIdToUt[uint(id)]
		case reflect.Array:
			typ = reflect.ArrayOf(int(desc.count), desc.elem)
		case reflect.Slice:
			typ = reflect.SliceOf(desc.elem)
		case reflect.Map:
			typ = reflect.MapOf(desc.key, desc.elem)
		}
	}
	if desc.hasDepth {
		for i := 0; i < int(desc.depth); i++ {
			typ = reflect.PtrTo(typ)
		}
	}
	return typ
}

func init() {
	typeIdToUt = make(map[uint]reflect.Type)
	baseRtToTypeId = make(map[reflect.Type]uint)
	kindToReflectType = make(map[reflect.Kind]reflect.Type)

	kindToReflectType[reflect.Bool] = reflect.TypeOf(false)
	kindToReflectType[reflect.Int] = reflect.TypeOf(int(0))
	kindToReflectType[reflect.Int8] = reflect.TypeOf(int8(0))
	kindToReflectType[reflect.Int16] = reflect.TypeOf(int16(0))
	kindToReflectType[reflect.Int32] = reflect.TypeOf(int32(0))
	kindToReflectType[reflect.Int64] = reflect.TypeOf(int64(0))
	kindToReflectType[reflect.Uint] = reflect.TypeOf(uint(0))
	kindToReflectType[reflect.Uint8] = reflect.TypeOf(uint8(0))
	kindToReflectType[reflect.Uint16] = reflect.TypeOf(uint16(0))
	kindToReflectType[reflect.Uint32] = reflect.TypeOf(uint32(0))
	kindToReflectType[reflect.Uint64] = reflect.TypeOf(uint64(0))
	kindToReflectType[reflect.Float32] = reflect.TypeOf(float32(0.0))
	kindToReflectType[reflect.Float64] = reflect.TypeOf(float64(0.0))
	kindToReflectType[reflect.String] = reflect.TypeOf("")
	t := make([]interface{}, 1)
	kindToReflectType[reflect.Interface] = reflect.TypeOf(t).Elem()
}
