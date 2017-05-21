package tcp

import (
	"errors"
	"github.com/sydnash/lotou/log"
	"net"
)

const (
	PARSE_STATUS_LEN int = iota
	PARSE_STATUS_MSG
	PARSE_STATUS_END
)

const (
	PACKAGE_LEN_COUNT int = 4
)

var ErrPacketLenExceed = errors.New("packet length exceed")

type ParseCache struct {
	msg           []byte
	msgLen        int
	copyLen       int
	status        int
	msgHeader     [PACKAGE_LEN_COUNT]byte
	copyHeaderLen int
}

func (p *ParseCache) reset() {
	p.msg = nil
	p.msgLen = 0
	p.copyLen = 0
	p.copyHeaderLen = 0
	p.status = PARSE_STATUS_LEN
}

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

func Subpackage(cache []byte, in net.Conn, status *ParseCache) (pack [][]byte, err error) {
READ_LOOP:
	for {
		if len(pack) > 0 {
			return pack, nil
		}
		n, err := in.Read(cache)
		if err != nil {
			return nil, err
		}

		startPos := 0
		for {
			switch status.status {
			case PARSE_STATUS_LEN:
				if len(cache[startPos:n]) < PACKAGE_LEN_COUNT-status.copyHeaderLen {
					copyLen := copy(status.msgHeader[status.copyHeaderLen:], cache[startPos:n])
					status.copyHeaderLen += copyLen
					if len(pack) == 0 {
						continue READ_LOOP
					} else {
						return pack, nil
					}
				}
				if status.copyHeaderLen == 0 {
					status.msgLen = int(ByteSliceToInt(cache[startPos:n]))
					startPos += PACKAGE_LEN_COUNT
				} else {
					copyLen := copy(status.msgHeader[status.copyHeaderLen:], cache[startPos:n])
					startPos += copyLen
					status.msgLen = int(ByteSliceToInt(status.msgHeader[:]))
				}
				if status.msgLen > MAX_PACKET_LEN {
					log.Error("packet length(%v) exceeds the maximum message length %v", status.msgLen, MAX_PACKET_LEN)
					return pack, ErrPacketLenExceed
				}
				tmp := make([]byte, status.msgLen)
				if status.copyHeaderLen != 0 {
					copy(tmp, status.msgHeader[:])
				} else {
					copy(tmp[0:PACKAGE_LEN_COUNT], cache[startPos-PACKAGE_LEN_COUNT:startPos])
				}
				status.status = PARSE_STATUS_MSG
				status.msg = tmp
				status.copyLen = PACKAGE_LEN_COUNT
			case PARSE_STATUS_MSG:
				copyLen := copy(status.msg[status.copyLen:], cache[startPos:n])
				status.copyLen += copyLen
				startPos += copyLen
				if status.copyLen != status.msgLen {
					if len(pack) == 0 {
						continue READ_LOOP
					} else {
						return pack, nil
					}
				}
				status.status = PARSE_STATUS_END
			case PARSE_STATUS_END:
				pack = append(pack, status.msg)
				status.reset()
			}
		}
	}
}
