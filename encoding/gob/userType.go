package gob

func init() {
	registerStructType(T1{})
	registerStructType(T2{})
}

type T1 struct {
	A uint
	B string
	C float64
	E int32
}

type T2 struct {
	T1
	F string
}
