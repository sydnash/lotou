package gob

import (
	"errors"
	"fmt"
)

//Pack pack data to bytes, if has error, it will panic
func Pack(data ...interface{}) []byte {
	encoder := NewEncoder()
	encoder.Reset()
	encoder.Encode(data)
	encoder.UpdateLen()
	buf := encoder.Buffer()
	return buf
}

//Pack pack data to bytes, if has error, it will return a error
func PackWithErr(data ...interface{}) (ret []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
			ret = nil
		}
	}()

	ret = Pack(data...)
	return ret, nil
}

//Unpack unpack bytes to value, return error
func Unpack(data []byte) (ret interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
			ret = nil
		}
	}()
	decoder := NewDecoder()
	decoder.SetBuffer(data)
	ret, ok := decoder.Decode()
	if !ok {
		return nil, errors.New("decode error.")
	}
	return ret, nil
}
