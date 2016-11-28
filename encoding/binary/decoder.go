package binary

import (
	"math"
	"reflect"
)

type Decoder struct {
	b    []byte
	r, w int
}

func NewDecoder() *Decoder {
	a := &Decoder{}
	return a
}

func (dec *Decoder) SetBuffer(b []byte) {
	dec.b = b
	dec.reset()
}

func (dec *Decoder) reset() {
	dec.r = 4
	dec.w = len(dec.b)
}

var (
	decoderMap map[reflect.Kind]func(*Decoder, reflect.Value)
)

func (dec *Decoder) Decode(i interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		panic("must ptr")
	}
	v = v.Elem()
	if v.Kind() == reflect.Ptr {
		panic("must only one depth of ptr")
	}
	dec.decodeValue(v)
}

func (enc *Decoder) decodeValue(v reflect.Value) {
	typ := v.Type()
	decodeFunc := findDecoder(typ)
	decodeFunc(enc, v)
}

func readUInt32(dec *Decoder) (a uint32) {
	idx := dec.r
	var i uint
	for i = 0; i < 4; i++ {
		a = a | uint32(dec.b[idx+int(i)])<<(i*8)
	}
	dec.r += 4
	return
}
func readUInt16(dec *Decoder) (a uint16) {
	idx := dec.r
	var i uint
	for i = 0; i < 2; i++ {
		a = a | uint16(dec.b[idx+int(i)])<<(i*8)
	}
	dec.r += 2
	return
}
func readUInt64(dec *Decoder) (a uint64) {
	idx := dec.r
	var i uint
	for i = 0; i < 8; i++ {
		a = a | uint64(dec.b[idx+int(i)])<<(i*8)
	}
	dec.r += 8
	return
}

func decodeInt32(dec *Decoder, v reflect.Value) {
	var a int32 = int32(readUInt32(dec))
	v.SetInt(int64(a))
}
func decodeInt16(dec *Decoder, v reflect.Value) {
	var a int16 = int16(readUInt16(dec))
	v.SetInt(int64(a))
}
func decodeInt64(dec *Decoder, v reflect.Value) {
	var a int64 = int64(readUInt64(dec))
	v.SetInt(a)
}
func decodeInt8(dec *Decoder, v reflect.Value) {
	var a int8 = int8(dec.b[dec.r])
	dec.r += 1
	v.SetInt(int64(a))
}
func decodeUInt32(dec *Decoder, v reflect.Value) {
	var a uint64 = uint64(readUInt32(dec))
	v.SetUint(a)
}
func decodeUInt16(dec *Decoder, v reflect.Value) {
	var a uint64 = uint64(readUInt16(dec))
	v.SetUint(a)
}
func decodeUInt64(dec *Decoder, v reflect.Value) {
	var a uint64 = uint64(readUInt64(dec))
	v.SetUint(a)
}
func decodeUInt8(dec *Decoder, v reflect.Value) {
	var a uint64 = uint64(dec.b[dec.r])
	dec.r += 1
	v.SetUint(a)
}
func decodeFloat32(dec *Decoder, v reflect.Value) {
	a := readUInt32(dec)
	f := math.Float32frombits(a)
	v.SetFloat(float64(f))
}
func decodeFloat64(dec *Decoder, v reflect.Value) {
	a := readUInt64(dec)
	ua := math.Float64frombits(a)
	v.SetFloat(ua)
}

func decodeString(dec *Decoder, v reflect.Value) {
	var count int = int(readUInt16(dec))
	str := string(dec.b[dec.r : dec.r+count])
	dec.r += count
	v.SetString(str)
}
func decodeBytes(dec *Decoder, v reflect.Value) {
	var count int = int(readUInt32(dec))
	str := string(dec.b[dec.r : dec.r+count])
	dec.r += count
	v.SetBytes([]byte(str))
}
func decodeStruct(dec *Decoder, v reflect.Value) {
	num := v.NumField()
	for i := 0; i < num; i++ {
		f := v.Field(i)
		dec.decodeValue(f)
	}
}
func decodeInterface(dec *Decoder, v reflect.Value) {
	rv := v.Elem()
	dec.decodeValue(rv)
}
func decodeBool(dec *Decoder, v reflect.Value) {
	a := dec.b[dec.r]
	dec.r += 1
	var b bool = false
	if a != 0 {
		b = true
	}
	v.SetBool(b)
}
func decodeInt(dec *Decoder, v reflect.Value) {
	decodeInt32(dec, v)
}
func decodeUInt(dec *Decoder, v reflect.Value) {
	decodeUInt32(dec, v)
}
func decodeSlice(dec *Decoder, v reflect.Value) {
	count := int(readUInt32(dec))
	v.SetLen(count)
	for i := 0; i < count; i++ {
		elem := v.Index(i)
		dec.decodeValue(elem)
	}
}
func decodeMap(dec *Decoder, v reflect.Value) {
	count := int(readUInt32(dec))
	typ := v.Type()
	for i := 0; i < count; i++ {
		key := reflect.New(typ.Key()).Elem()
		dec.decodeValue(key)
		value := reflect.New(typ.Elem()).Elem()
		dec.decodeValue(value)
		v.SetMapIndex(key, value)
	}
}

func init() {
	decoderMap = make(map[reflect.Kind]func(*Decoder, reflect.Value))
	decoderMap[reflect.Bool] = decodeBool
	decoderMap[reflect.Int] = decodeInt
	decoderMap[reflect.Int8] = decodeInt8
	decoderMap[reflect.Int16] = decodeInt16
	decoderMap[reflect.Int32] = decodeInt32
	decoderMap[reflect.Int64] = decodeInt64
	decoderMap[reflect.Uint] = decodeUInt
	decoderMap[reflect.Uint8] = decodeUInt8
	decoderMap[reflect.Uint16] = decodeUInt16
	decoderMap[reflect.Uint32] = decodeUInt32
	decoderMap[reflect.Uint64] = decodeUInt64
	decoderMap[reflect.Float32] = decodeFloat32
	decoderMap[reflect.Float64] = decodeFloat64
	decoderMap[reflect.Struct] = decodeStruct
	decoderMap[reflect.String] = decodeString
	decoderMap[reflect.Slice] = decodeSlice
	decoderMap[reflect.Map] = decodeMap
}

func findDecoder(typ reflect.Type) func(*Decoder, reflect.Value) {
	if typ == reflect.TypeOf([]byte{}) {
		return decodeBytes
	}
	return decoderMap[typ.Kind()]
}
