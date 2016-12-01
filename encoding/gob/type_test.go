package gob_test

import "github.com/sydnash/lotou/encoding/gob"
import "fmt"
import "testing"
import "reflect"

func TestType(t *testing.T) {
	enc := gob.NewEncoder()

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
	a = append(a, &gob.T1{10, "哈哈，这都可以？", 1.5, -100})
	a = append(a, &gob.T2{gob.T1{10, "哈哈，这都可以？", 1.5, -100}, "那么这样还可以吗？"})
	a = append(a, gob.T1{10, "哈哈，这都可以？", 1.5, -100})
	a = append(a, gob.T2{gob.T1{10, "哈哈，这都可以？", 1.5, -100}, "那么这样还可以吗？"})
	a = append(a, true)
	a = append(a, false)
	a = append(a, [3]int{1, 2, 3})
	a = append(a, []byte{})
	m := make(map[int]string)
	m[1] = "map的第一个元素"
	m[1] = "map的第二个元素"
	a = append(a, m)
	s := make([]string, 0, 2)
	s = append(s, "这是slice的元素")
	a = append(a, s)
	str := "这是一个[]byte"
	s1 := []byte(str)
	a = append(a, s1)

	b := make([]interface{}, 0, 10)
	b = append(b, m)
	b = append(b, s)
	b = append(b, s1)
	a = append(a, b)
	a = append(a, a)
	for _, v := range a {
		enc.Encode(v)
	}
	dec := gob.NewDecoder()
	dec.SetBuffer(enc.Buffer())

	var ok bool = true
	var r interface{}
	idx := 0
	for ok {
		r, ok = dec.Decode()
		fmt.Println(r, reflect.TypeOf(r), ok)
		if ok {
			if !reflect.DeepEqual(r, a[idx]) {
				t.Errorf("%v is not equal to %v at idx %v", r, a[idx], idx)
			}
			if reflect.TypeOf(r) != reflect.TypeOf(a[idx]) {
				t.Errorf("%v is not equal to %v at idx %v", reflect.TypeOf(r), reflect.TypeOf(a[idx]), idx)
			}
			idx++
		}
	}

}
