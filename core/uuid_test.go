package core_test

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/sydnash/lotou/core"
	"testing"
)

func TestUUID(t *testing.T) {
	pp := core.UUID()
	var p []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(p, pp)

	str := base64.StdEncoding.EncodeToString(p)

	tt := md5.Sum([]byte(str))

	str = base64.StdEncoding.EncodeToString(tt[:])
	fmt.Println(string(str))
}
