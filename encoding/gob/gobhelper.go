package gob

import (
	"fmt"
	"hash/fnv"
	"reflect"
)

func getTypeFrag(x interface{}) string {
	return getTypeFragStr(reflect.TypeOf(x))
}

//~ Should be private
func getTypeFragStr(t reflect.Type) string {
	var frag string
	switch t.Kind() {
	//~ Primitive types.
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Uintptr, reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.String:
		frag = fmt.Sprintf("%v", t)
		break
	//~ Types that do not proceed
	case reflect.Chan, reflect.Func:
		panic(fmt.Sprintf("Canno process such type %v", t.Kind()))
	case reflect.Ptr:
		frag = getPtrFlatStr(t)
	case reflect.Array:
		frag = getArrayFlatStr(t)
	case reflect.Map:
		frag = getMapFlatStr(t)
	case reflect.Slice:
		frag = getSliceFlatStr(t)
		break
	case reflect.Struct:
		frag = getStructFlatStr(t)
	default:
		panic(fmt.Sprintf("Unknown case of %v", t.Kind()))
	}
	if len(frag) == 0 {
		panic("Cannot get type of this:Something missing")
	}
	return frag
}

func getSliceFlatStr(t reflect.Type) string {
	if t.Kind() != reflect.Slice {
		panic(`Not a slice type`)
	}
	return fmt.Sprintf("[-%v]", getTypeFragStr(t.Elem()))
}

func getMapFlatStr(t reflect.Type) string {
	if t.Kind() != reflect.Map {
		panic(`Not a map as expected`)
	}
	return fmt.Sprintf("{map[%v]%v}", getTypeFragStr(t.Key()), getTypeFragStr(t.Elem()))
}

func getPtrFlatStr(t reflect.Type) string {
	if t.Kind() != reflect.Ptr {
		panic(`Not a ptr as expected`)
	}
	return fmt.Sprintf("*%v", getTypeFragStr(t.Elem()))
}

func getArrayFlatStr(t reflect.Type) string {
	if t.Kind() != reflect.Array {
		panic(`Not an array as expected`)
	}
	return fmt.Sprintf("[(%v)-%v]", t.Len(), getTypeFragStr(t.Elem()))
}

func getStructFlatStr(t reflect.Type) string {
	if t.Kind() != reflect.Struct {
		panic(`Not a struct as expected`)
	}
	n := t.NumField()
	var s string = fmt.Sprintf("%v=Struct{<%v>::", t.Name(), n)
	for i := 0; i < n; i++ {
		f := t.Field(i) //~ StructField
		s += fmt.Sprintf("%v<%v>;", f.Name, getTypeFragStr(f.Type))
	}
	s += "}"
	return s
}

func getStructID(t reflect.Type) uint {
	strrpr := getTypeFragStr(t)
	h := fnv.New64a()
	h.Sum([]byte(strrpr))
	return uint(h.Sum64())
}

/* Uncomment the rest to get the original implementation
var (
	_baseID uint
)

func getStructID() uint {
	_baseID++
	return _baseID
}
*/
