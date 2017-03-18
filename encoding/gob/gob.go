package gob

func Pack(data []interface{}) []byte {
	encoder := NewEncoder()
	encoder.Reset()
	encoder.Encode(data)
	encoder.UpdateLen()
	buf := encoder.Buffer()
	return buf
}
func Unpack(data []byte) interface{} {
	decoder := NewDecoder()
	decoder.SetBuffer(data)
	sdata, ok := decoder.Decode()
	if !ok {
		panic("gob unpack failed")
	}
	return sdata
}
