package core

type Base struct {
	id uint
	in chan *Message
}

func (self *Base) Id() uint {
	return self.id
}
func (self *Base) Send(m *Message) {
	self.in <- m
}
func (self *Base) SetId(id uint) {
	self.id = id
}
func (self *Base) In() chan *Message {
	return self.in
}
func (self *Base) Close() {
	close(self.in)
}

func NewBase() *Base {
	a := &Base{}
	a.in = make(chan *Message, 1024)
	return a
}
