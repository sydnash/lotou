package binary

import (
	"math"
	"reflect"
)

//Encoder encode value to binary start with 4 byte for binary len.
//and only 2 byte to encode map, slice's len.
type Encoder struct {
	b    []byte
	r, w int
}

func (enc *Encoder) Reset() {
	if enc.b == nil {
		enc.b = make([]byte, 1024, 1024)
	}
	enc.r = 0
	enc.w = 4
	enc.b = enc.b[:]
}

func NewEncoder() *Encoder {
	a := &Encoder{}
	a.Reset()
	return a
}

var (
	encoderMap map[reflect.Kind]func(*Encoder, reflect.Value)
)

func (enc *Encoder) Buffer() []byte {
	return enc.b[:enc.w]
}

func (enc *Encoder) Encode(i interface{}) {
	enc.encodeValue(reflect.ValueOf(i))
}

func (enc *Encoder) encodeValue(v reflect.Value) {
	typ := v.Type()
	encodeFunc := findEncoder(typ)
	encodeFunc(enc, v)
}

func (enc *Encoder) grow(n int) {
	if enc.w+n > cap(enc.b) {
		buf := make([]byte, 2*cap(enc.b), 2*cap(enc.b))
		copy(buf, enc.b[:enc.w])
		enc.b = buf[:]
	}
}
func intToByteSlice(v uint32) []byte {
	a := make([]byte, 4)
	a[3] = byte((v >> 24) & 0xFF)
	a[2] = byte((v >> 16) & 0XFF)
	a[1] = byte((v >> 8) & 0XFF)
	a[0] = byte(v & 0XFF)
	return a
}
func (enc *Encoder) UpdateLen() {
	l := enc.w
	b := intToByteSlice(uint32(l))
	copy(enc.b[:4], b[:])
}

func encodeInt32(enc *Encoder, v reflect.Value) {
	a := uint32(v.Interface().(int32))
	enc.grow(4)
	idx := enc.w
	for i := 0; i < 4; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 4
}
func encodeInt16(enc *Encoder, v reflect.Value) {
	a := uint16(v.Interface().(int16))
	enc.grow(2)
	idx := enc.w
	for i := 0; i < 2; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 2
}
func encodeInt64(enc *Encoder, v reflect.Value) {
	a := uint64(v.Interface().(int64))
	enc.grow(8)
	idx := enc.w
	for i := 0; i < 8; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 8
}
func encodeInt8(enc *Encoder, v reflect.Value) {
	a := uint8(v.Interface().(int8))
	enc.grow(1)
	enc.b[enc.w] = a
	enc.w += 1
}
func encodeUInt32(enc *Encoder, v reflect.Value) {
	a := v.Interface().(uint32)
	enc.grow(4)
	idx := enc.w
	for i := 0; i < 4; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 4
}
func encodeUInt16(enc *Encoder, v reflect.Value) {
	a := v.Interface().(uint16)
	enc.grow(2)
	idx := enc.w
	for i := 0; i < 2; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 2
}
func encodeUInt64(enc *Encoder, v reflect.Value) {
	a := v.Interface().(uint64)
	enc.grow(8)
	idx := enc.w
	for i := 0; i < 8; i++ {
		enc.b[idx+i] = byte(a >> (uint(i) * 8))
	}
	enc.w += 8
}
func encodeUInt8(enc *Encoder, v reflect.Value) {
	a := v.Interface().(uint8)
	enc.grow(1)
	enc.b[enc.w] = a
	enc.w += 1
}
func encodeFloat32(enc *Encoder, v reflect.Value) {
	a := v.Interface().(float32)
	ua := math.Float32bits(a)
	encodeUInt32(enc, reflect.ValueOf(ua))
}
func encodeFloat64(enc *Encoder, v reflect.Value) {
	a := v.Interface().(float64)
	ua := math.Float64bits(a)
	encodeUInt64(enc, reflect.ValueOf(ua))
}

func encodeString(enc *Encoder, v reflect.Value) {
	str := v.Interface().(string)
	encodeBytes(enc, reflect.ValueOf([]byte(str)))
}
func encodeBytes(enc *Encoder, v reflect.Value) {
	bytes := v.Bytes()
	count := len(bytes)
	encodeInt16(enc, reflect.ValueOf(int16(count)))
	enc.grow(count)
	copy(enc.b[enc.w:], bytes[:])
	enc.w += count
}
func encodeStruct(enc *Encoder, v reflect.Value) {
	num := v.NumField()
	for i := 0; i < num; i++ {
		f := v.Field(i)
		if f.CanInterface() {
			enc.encodeValue(f)
		}
	}
}
func encodeInterface(enc *Encoder, v reflect.Value) {
	rv := v.Elem()
	enc.encodeValue(rv)
}
func encodeBool(enc *Encoder, v reflect.Value) {
	b := v.Bool()
	var a uint8 = 0
	if b {
		a = 1
	}
	encodeUInt8(enc, reflect.ValueOf(a))
}
func encodeInt(enc *Encoder, v reflect.Value) {
	var a uint32 = uint32(v.Interface().(int))
	encodeUInt32(enc, reflect.ValueOf(a))
}
func encodeUInt(enc *Encoder, v reflect.Value) {
	var a uint32 = uint32(v.Interface().(uint))
	encodeUInt32(enc, reflect.ValueOf(a))
}
func encodeSlice(enc *Encoder, v reflect.Value) {
	count := v.Len()
	encodeInt32(enc, reflect.ValueOf(int32(count)))
	for i := 0; i < count; i++ {
		elem := v.Index(i)
		enc.encodeValue(elem)
	}
}
func encodeMap(enc *Encoder, v reflect.Value) {
	count := v.Len()
	encodeInt32(enc, reflect.ValueOf(int32(count)))
	keys := v.MapKeys()
	for _, k := range keys {
		enc.encodeValue(k)
		elem := v.MapIndex(k)
		enc.encodeValue(elem)
	}
}

func init() {
	encoderMap = make(map[reflect.Kind]func(*Encoder, reflect.Value))
	encoderMap[reflect.Bool] = encodeBool
	encoderMap[reflect.Int] = encodeInt
	encoderMap[reflect.Int8] = encodeInt8
	encoderMap[reflect.Int16] = encodeInt16
	encoderMap[reflect.Int32] = encodeInt32
	encoderMap[reflect.Int64] = encodeInt64
	encoderMap[reflect.Uint] = encodeUInt
	encoderMap[reflect.Uint8] = encodeUInt8
	encoderMap[reflect.Uint16] = encodeUInt16
	encoderMap[reflect.Uint32] = encodeUInt32
	encoderMap[reflect.Uint64] = encodeUInt64
	encoderMap[reflect.Float32] = encodeFloat32
	encoderMap[reflect.Float64] = encodeFloat64
	encoderMap[reflect.Struct] = encodeStruct
	encoderMap[reflect.String] = encodeString
	encoderMap[reflect.Slice] = encodeSlice
	encoderMap[reflect.Map] = encodeMap
}

func findEncoder(typ reflect.Type) func(*Encoder, reflect.Value) {
	if typ == reflect.TypeOf([]byte{}) {
		return encodeBytes
	}
	return encoderMap[typ.Kind()]
}
