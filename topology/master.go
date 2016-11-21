package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type master struct {
	*core.Base
	decoder *gob.Decoder
	encoder *gob.Encoder
}

func StartMaster(ip, port) {
	m := &master{Base: core.NewBase()}
	m.decoder = gob.NewDecoder
	m.encoder = gob.NewEncoder
	core.RegisterService(m)
	core.Name(".master", m.Id())
	s := tcp.New(ip, port, m.Id())
	s.Listen()
}

func (m *master) run() {
	m.SetSelf(m)
	m.RegisterBaseCB(MSG_TYPE_CLOSE, (*master).close, true)
	m.RegisterBaseCB(MSG_TYPE_NORMAL, (*master).normalMSG, false)
	for s := range m.In() {
		m.DispatchM(msg)
	}
}

func (m *master) normalMSG(dest, src uint, msgEncode string, data ...interface{}) {
	if msgEncode == "go" {
	} else if msgEncode == "socket" {
	}
}

func (m *master) close(dest, src uint) {
	_, _ = dest, src
	m.Close()
}
