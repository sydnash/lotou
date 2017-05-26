package gob_test

import (
	"fmt"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"testing"
)

func TestPackMore(t *testing.T) {
	specificType := core.MSG_TYPE_NORMAL
	m1 := &core.Message{
		Type: specificType,
		Cmd:  "", //~ This line can not be ignored.
	}
	bytes := gob.Pack(m1)
	fmt.Println(bytes)
	m2, _ := gob.Unpack(bytes)
	m3 := m2.([]interface{})[0].(*core.Message)
	if m3.Type != specificType {
		t.Error("Not the same")
	}
}
