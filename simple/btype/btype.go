package btype

type PHead struct {
	Len     int32
	Flag    int32
	Id      int16
	SType   int32
	OriLen  int32
	Type    int32
	reqId   int32
	session int32
	acid    int32
}
