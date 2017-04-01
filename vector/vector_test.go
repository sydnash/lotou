package vector

import (
	"reflect"
	"sort"
	"testing"
)

func TestLenAndCap(t *testing.T) {
	v := New()

	if v.Len() != 0 || v.Cap() != 0 {
		t.Error("len and cap must be 0")
	}

	v = NewCap(10)
	if v.Len() != 0 || v.Cap() != 10 {
		t.Error("len must be 0 and cap must 10")
	}
}

func errorOnDeepEqual(a []interface{}, b []interface{}, s string, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Error(s)
	}
}

func TestAppend(t *testing.T) {
	v := New()
	v.Append([]interface{}{1, 2, 3}...)
	if !reflect.DeepEqual(v.s, []interface{}{1, 2, 3}) {
		t.Error("append")
	}
	v.Append(4)
	if !reflect.DeepEqual(v.s, []interface{}{1, 2, 3, 4}) {
		t.Error("append")
	}
	v.AppendVec(v)
	if !reflect.DeepEqual(v.s, []interface{}{1, 2, 3, 4, 1, 2, 3, 4}) {
		t.Error("append vec")
	}
}

func TestClone(t *testing.T) {
	v := NewCap(10)
	v.Append([]interface{}{1, 2, 3, 4}...)
	v1 := v.Clone()
	if !reflect.DeepEqual(v.s, v1.s) {
		t.Error("clone")
	}
}

func TestDelete(t *testing.T) {
	v := NewCap(10)
	v.Append([]interface{}{1, 2, 3, 4}...)

	v.Delete(0)

	if !reflect.DeepEqual(v.s, []interface{}{2, 3, 4}) {
		t.Error("delete")
	}

	v.Insert(0, 5)
	errorOnDeepEqual(v.s, []interface{}{5, 2, 3, 4}, "insert", t)

	v.InsertVariant(2, 6, 6)
	errorOnDeepEqual(v.s, []interface{}{5, 2, 6, 6, 3, 4}, "Insert", t)

	e2 := v.At(1)
	if e2 != 2 {
		t.Error("at")
	}

	v.Reverse()
	errorOnDeepEqual(v.s, []interface{}{4, 3, 6, 6, 2, 5}, "reverse", t)

	v.PushFront(0)
	v.Push(0)

	errorOnDeepEqual(v.s, []interface{}{0, 4, 3, 6, 6, 2, 5, 0}, "push", t)

	v.Pop()
	v.PopFront()

	errorOnDeepEqual(v.s, []interface{}{4, 3, 6, 6, 2, 5}, "reverse", t)

	n := New()
	n.Copy(v)
	errorOnDeepEqual(n.s, []interface{}{4, 3, 6, 6, 2, 5}, "Copy", t)

	n.DeleteByValue(4)
	errorOnDeepEqual(n.s, []interface{}{3, 6, 6, 2, 5}, "delete by value", t)
}
