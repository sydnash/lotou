package tcp

import (
	"bufio"
	"errors"
	"github.com/sydnash/lotou/log"
)

var ErrPacketLenExceed = errors.New("packet length exceed")

func IntToByteSlice(v uint32) []byte {
	a := make([]byte, 4)
	a[3] = byte((v >> 24) & 0xFF)
	a[2] = byte((v >> 16) & 0XFF)
	a[1] = byte((v >> 8) & 0XFF)
	a[0] = byte(v & 0XFF)
	return a
}
func ByteSliceToInt(s []byte) (v uint32) {
	v = uint32(s[3])<<24 | uint32(s[2])<<16 | uint32(s[1])<<8 | uint32(s[0])
	return v
}

func Subpackage(in *bufio.Reader) (pack []byte, err error) {
	packageLenCount := 4
	var packLen int
	for {
		pre, err := in.Peek(packageLenCount)
		if err != nil {
			return nil, err
		}
		if len(pre) == packageLenCount {
			packLen = int(ByteSliceToInt(pre))
			break
		}
	}
	if packLen > MAX_PACKET_LEN {
		log.Error("packet length exceeds the maximum range: %v", packLen)
		return pack, ErrPacketLenExceed
	}
	for {
		t, err := in.Peek(packLen)
		if err != nil {
			return nil, err
		}
		if len(t) == packLen {
			pack = make([]byte, packLen)
			rlen, err := in.Read(pack)
			if err != nil {
				return nil, err
			}
			if rlen == packLen {
				return pack, err
			}
		}
	}
}
