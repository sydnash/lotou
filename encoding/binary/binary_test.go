package binary_test

import (
	"fmt"
	"github.com/sydnash/lotou/encoding/binary"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	enc := binary.NewEncoder()

	a := make([]interface{}, 0, 10)
	a = append(a, int(-5))
	a = append(a, int8(-1))
	a = append(a, int16(-2))
	a = append(a, int32(-3))
	a = append(a, int64(-4))
	a = append(a, uint(6))
	a = append(a, uint8(7))
	a = append(a, uint16(8))
	a = append(a, uint32(9))
	a = append(a, uint64(10))
	a = append(a, float32(0.99999))
	a = append(a, float64(0.9999999999))
	a = append(a, "this is a string")
	a = append(a, "这也是一个string")
	t1 := struct {
		A string
		B []byte
	}{A: "结构体的字符串", B: []byte("结构体的字符串")}
	a = append(a, t1)
	for _, v := range a {
		enc.Encode(v)
	}

	dec := binary.NewDecoder()
	dec.SetBuffer(enc.Buffer())

	for i := 0; i < len(a); i++ {
		typ := reflect.TypeOf(a[i])
		v := reflect.New(typ)
		dec.Decode(v.Interface())
		fmt.Println(v.Elem().Interface())
		if !reflect.DeepEqual(v.Elem().Interface(), a[i]) {
			t.Errorf("%v is not equal to %v at idx %v", v.Elem().Interface(), a[i], i)
		}
	}
}
