package lotou

import (
	"fmt"
	"testing"
)

func normalFunc(i int, n string) {
	fmt.Println(i, n)
}

type Int int

func (self Int) Add(i int) {
	fmt.Println(self + Int(i))
}

func TestHelper(t *testing.T) {
	callHelper := NewCallHelper()

	callHelper.AddFunc("t1", normalFunc)

	i := Int(1)

	callHelper.AddMethod("t2", i, "Add")

	callHelper.Call("t1", 1, "name")
	callHelper.Call("t2", 2)

	a := []interface{}{1, "name"}
	callHelper.Call("t1", a...)

	b := []interface{}{2}
	callHelper.Call("t2", b...)
}
