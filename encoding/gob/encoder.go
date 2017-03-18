package gob

import (
	"fmt"
	"math"
	"reflect"
)

type Encoder struct {
	b    []byte
	r, w int
}

func NewEncoder() *Encoder {
	a := &Encoder{}
	a.Reset()
	return a
}

func intToByteSlice(v uint32) []byte {
	a := make([]byte, 4)
	a[3] = byte((v >> 24) & 0xFF)
	a[2] = byte((v >> 16) & 0XFF)
	a[1] = byte((v >> 8) & 0XFF)
	a[0] = byte(v & 0XFF)
	return a
}
func ByteSliceToInt(s []byte) (v uint32) {
	v = uint32(s[3])<<24 | uint32(s[2])<<16 | uint32(s[1])<<8 | uint32(s[0])
	return v
}

func (enc *Encoder) SetBuffer(b []byte) {
	enc.b = b
	enc.Reset()
}
func (enc *Encoder) Buffer() []byte {
	return enc.b[:enc.w]
}

func (enc *Encoder) UpdateLen() {
	l := enc.w
	b := intToByteSlice(uint32(l))
	copy(enc.b[:4], b[:])
}

func (enc *Encoder) Reset() {
	if enc.b == nil {
		enc.b = make([]byte, 1024, 1024)
	}
	enc.r = 0
	enc.w = 4
	enc.b = enc.b[:]
}

func (enc *Encoder) grow(n int) {
	if enc.w+n > cap(enc.b) {
		buf := make([]byte, 2*cap(enc.b), 2*cap(enc.b))
		copy(buf, enc.b[:enc.w])
		enc.b = buf[:]
	}
}

func (enc *Encoder) Encode(i interface{}) {
	value := reflect.ValueOf(i)
	enc.encodeValue(value)
}

func (enc *Encoder) encodeValue(v reflect.Value) {
	rt, depth := findBaseAndDepth(v.Type())
	enc.encodeType(v.Type())
	enc.encodeConcreteValue(rt, depth, v)
}

func (enc *Encoder) encodeConcreteValue(rt reflect.Type, depth uint, v reflect.Value) {
	value := v
	for i := 0; i < int(depth); i++ {
		value = value.Elem()
	}
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.encodeInt(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		enc.encodeUInt(value.Uint())
	case reflect.Float64, reflect.Float32:
		tmp := value.Float()
		tmp64 := math.Float64bits(float64(tmp))
		enc.encodeUInt(tmp64)
	case reflect.Bool:
		t := value.Bool()
		var v uint64 = 0
		if t {
			v = 1
		}
		enc.encodeUInt(v)
	case reflect.String:
		str := value.String()
		enc.encodeString(str)
	case reflect.Struct:
		enc.encodeStruct(value)
	case reflect.Array, reflect.Slice:
		enc.encodeArrayLike(value)
	case reflect.Map:
		enc.encodeMap(value)
	case reflect.Interface:
	default:
		panic(fmt.Sprintf("not support type to send. %v", rt))
	}
}

func (enc *Encoder) encodeType(typ reflect.Type) {
	id := valueToId(typ)
	enc.encodeUInt(uint64(id))
	kind, _, _ := parseTypeId(id)
	switch reflect.Kind(kind) {
	case reflect.Slice:
		et := typ.Elem()
		enc.encodeType(et)
	case reflect.Map:
		kt := typ.Key()
		enc.encodeType(kt)
		et := typ.Elem()
		enc.encodeType(et)
	case reflect.Array:
		et := typ.Elem()
		enc.encodeType(et)
		le := typ.Len()
		enc.encodeUInt(uint64(le))
		/*	case reflect.Chan:
			et := typ.Elem()
			enc.encodeType(et)
		*/
	default:
		return
	}
}
func (enc *Encoder) encodeInt(v int64) {
	enc.grow(9)
	ux := uint64(v) << 1
	if v < 0 {
		ux = ^ux
	}
	enc.encodeUInt(ux)
}
func (enc *Encoder) encodeUInt(v uint64) {
	enc.grow(9)
	for v >= 0x80 {
		enc.b[enc.w] = byte(v) | 0x80
		v >>= 7
		enc.w++
	}
	enc.b[enc.w] = byte(v)
	enc.w++
}
func (enc *Encoder) encodeString(str string) {
	l := len(str)
	enc.encodeUInt(uint64(l))

	enc.grow(l)
	b := []byte(str)
	copy(enc.b[enc.w:], b[:])
	enc.w = enc.w + l
}
func (enc *Encoder) encodeStruct(value reflect.Value) {
	num := value.NumField()
	for i := 0; i < num; i++ {
		v := value.Field(i)
		if v.CanInterface() {
			enc.encodeValue(reflect.ValueOf(v.Interface()))
		}
	}
}
func (enc *Encoder) encodeArrayLike(value reflect.Value) {
	num := value.Len()
	enc.encodeUInt(uint64(num))

	if value.Type() == reflect.TypeOf([]byte{}) {
		enc.grow(num)
		b := value.Interface().([]byte)
		copy(enc.b[enc.w:], b[:])
		enc.w = enc.w + num
	} else {
		for i := 0; i < num; i++ {
			v := value.Index(i)
			enc.encodeValue(reflect.ValueOf(v.Interface())) //if pass v direct, it's type kind will be interface when the slice is []interface{}
		}
	}
}

func (enc *Encoder) encodeMap(value reflect.Value) {
	num := value.Len()
	enc.encodeUInt(uint64(num))
	keys := value.MapKeys()
	for _, k := range keys {
		v := value.MapIndex(k)
		enc.encodeValue(reflect.ValueOf(k.Interface()))
		enc.encodeValue(reflect.ValueOf(v.Interface()))
	}
}
