package core

import (
	"time"
)

type service struct {
	id           uint
	name         string
	msgChan      chan *Message
	loopTicker   *time.Ticker
	loopDuration int //unit is Millisecond
	m            Module
}

func newService(name string) *service {
	s := &service{name: name}
	s.msgChan = make(chan *Message, 1024)
	return s
}

func (s *service) setModule(m Module) {
	s.m = m
}

func (s *service) getName() string {
	return s.name
}

func (s *service) setId(id uint) {
	s.id = id
}

func (s *service) getId() uint {
	return s.id
}

func (s *service) pushMSG(m *Message) {
	s.msgChan <- m
}

func (s *service) dispatchMSG(msg *Message) {
	if msg.EncType == MSG_ENC_TYPE_GO {
		t := unpack(msg.Data[0].([]byte))
		msg.Data = t.([]interface{})
	}
	switch msg.Type {
	case MSG_TYPE_NORMAL:
		s.m.OnNormalMSG(msg.Dst, msg.Data...)
	}
}

func (s *service) loop() {
EXIT:
	for {
		select {
		case msg, ok := <-s.msgChan:
			if !ok {
				break EXIT
			}
			s.dispatchMSG(msg)
		}
	}
	s.m.OnDestroy()
}

func (s *service) loopWithLoop() {
EXIT:
	for {
		select {
		case msg, ok := <-s.msgChan:
			if !ok {
				break EXIT
			}
			s.dispatchMSG(msg)
		case <-s.loopTicker.C:
			s.m.OnMainLoop(s.loopDuration)
		}
	}
	s.loopTicker.Stop()
	s.m.OnDestroy()
}

func (s *service) run() {
	go s.loop()
}

func (s *service) runWithLoop(d int) {
	s.loopTicker = time.NewTicker(time.Duration(d) * time.Millisecond)
	go s.loopWithLoop()
}
