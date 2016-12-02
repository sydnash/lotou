package gob

import (
	"fmt"
	"math"
	"reflect"
)

type Decoder struct {
	b    []byte
	r, w int
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (dec *Decoder) SetBuffer(b []byte) {
	dec.b = b
	dec.reset()
}

func (dec *Decoder) reset() {
	dec.r = 4
	dec.w = len(dec.b)
}

func (dec *Decoder) Decode() (ret interface{}, ok bool) {
	if dec.r >= dec.w {
		return nil, false
	}
	v := dec.decodeValue()
	return v.Interface(), true
}

func (dec *Decoder) decodeValue() (value reflect.Value) {
	kind, depth, structId, typ := dec.decodeType(false)
	value = dec.decodeConcreteValue(kind, depth, structId, typ)
	return
}

func (dec *Decoder) decodeConcreteValue(kind, depth, structId uint, typ reflect.Type) (value reflect.Value) {
	_ = structId
	tk := reflect.Kind(kind)
	switch tk {
	case reflect.Int8:
		t := int8(dec.decodeInt())
		value = reflect.ValueOf(t)
	case reflect.Int16:
		t := int16(dec.decodeInt())
		value = reflect.ValueOf(t)
	case reflect.Int32:
		t := int32(dec.decodeInt())
		value = reflect.ValueOf(t)
	case reflect.Int64:
		t := dec.decodeInt()
		value = reflect.ValueOf(t)
	case reflect.Int:
		t := int(dec.decodeInt())
		value = reflect.ValueOf(t)
	case reflect.Uint8:
		t := uint8(dec.decodeUInt())
		value = reflect.ValueOf(t)
	case reflect.Uint16:
		t := uint16(dec.decodeUInt())
		value = reflect.ValueOf(t)
	case reflect.Uint32:
		t := uint32(dec.decodeUInt())
		value = reflect.ValueOf(t)
	case reflect.Uint64:
		t := dec.decodeUInt()
		value = reflect.ValueOf(t)
	case reflect.Uint:
		t := uint(dec.decodeUInt())
		value = reflect.ValueOf(t)
	case reflect.Float32:
		t := dec.decodeUInt()
		f32 := float32(math.Float64frombits(t))
		value = reflect.ValueOf(f32)
	case reflect.Float64:
		t := dec.decodeUInt()
		f64 := math.Float64frombits(t)
		value = reflect.ValueOf(f64)
	case reflect.Bool:
		v := dec.decodeUInt()
		var b bool = false
		if v != 0 {
			b = true
		}
		value = reflect.ValueOf(b)
	case reflect.String:
		str := dec.decodeString()
		value = reflect.ValueOf(str)
	case reflect.Struct:
		value = dec.decodeStruct(structId)
	case reflect.Slice:
		value = dec.decodeSlice(typ)
	case reflect.Array:
		value = dec.decodeArray(typ)
	case reflect.Map:
		value = dec.decodeMap(typ)
	default:
		panic("not support type")
	}
	for i := 0; i < int(depth); i++ {
		typ := value.Type()
		valuep := reflect.New(typ)
		valuep.Elem().Set(value)
		value = valuep
	}
	return
}

func (dec *Decoder) decodeType(hasDepth bool) (uint, uint, uint, reflect.Type) {
	id := uint16(dec.decodeUInt())
	kind, depth, structId := parseTypeId(id)
	var count uint
	var key reflect.Type
	var elem reflect.Type

	switch reflect.Kind(kind) {
	case reflect.Slice:
		_, _, _, elem = dec.decodeType(true)
	case reflect.Map:
		_, _, _, key = dec.decodeType(true)
		_, _, _, elem = dec.decodeType(true)
	case reflect.Array:
		_, _, _, elem = dec.decodeType(true)
		count = uint(dec.decodeUInt())
	/*	case reflect.Chan:
		et := typ.Elem()
		enc.encodeType(et)
	*/
	default:
	}
	var typ reflect.Type
	typ = createType(&TypeDesc{kind, depth, structId, count, key, elem, hasDepth})

	return kind, depth, structId, typ
}

func (dec *Decoder) decodeUInt() uint64 {
	var x uint64
	var s uint
	t := 0
	for i, b := range dec.b[dec.r:] {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				panic("decode uint error.")
				return 0
			}
			t = i + 1
			dec.r = dec.r + t
			return x | uint64(b)<<s
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	panic("decode uint error.")
	return 0
}

func (dec *Decoder) decodeInt() int64 {
	ux := dec.decodeUInt() // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

func (dec *Decoder) decodeString() string {
	l := dec.decodeUInt()
	str := string(dec.b[dec.r : dec.r+int(l)])
	dec.r = dec.r + int(l)
	return str
}

func (dec *Decoder) decodeStruct(structId uint) (value reflect.Value) {
	typ, ok := typeIdToUt[structId]
	if !ok {
		panic(fmt.Sprintf("struct type %v is not registed.", structId))
	}

	pv := reflect.New(typ)
	value = pv.Elem()

	num := value.NumField()
	for i := 0; i < num; i++ {
		f := value.Field(i)
		if f.CanInterface() {
			v := dec.decodeValue()
			f.Set(v)
		}
	}

	return
}

func (dec *Decoder) decodeSlice(typ reflect.Type) (value reflect.Value) {
	count := int(dec.decodeUInt())
	if typ == reflect.TypeOf([]byte{}) {
		v := reflect.New(typ)
		sli := make([]byte, count)
		copy(sli, dec.b[dec.r:])
		dec.r += count
		v.Elem().Set(reflect.ValueOf(sli))
		return v.Elem()
	} else {
		value = reflect.MakeSlice(typ, 0, count)
		for i := 0; i < count; i++ {
			v := dec.decodeValue()
			value = reflect.Append(value, v)
		}
	}
	return
}

func (dec *Decoder) decodeArray(typ reflect.Type) (value reflect.Value) {
	count := int(dec.decodeUInt())
	value = reflect.New(typ)
	value = value.Elem()
	for i := 0; i < count; i++ {
		v := dec.decodeValue()
		value.Index(i).Set(v)
	}
	return
}

func (dec *Decoder) decodeMap(typ reflect.Type) (value reflect.Value) {
	count := int(dec.decodeUInt())
	value = reflect.MakeMap(typ)
	for i := 0; i < count; i++ {
		key := dec.decodeValue()
		elem := dec.decodeValue()
		value.SetMapIndex(key, elem)
	}
	return
}
