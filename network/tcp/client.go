package tcp

import (
	"bufio"
	"github.com/sydnash/majiang/core"
	"github.com/sydnash/majiang/log"
	"net"
	"time"
)

type Client struct {
	*core.Base
	Con           *net.TCPConn
	RemoteAddress *net.TCPAddr
	Dest          uint
	inbuffer      *bufio.Reader
	outbuffer     *bufio.Writer
}

func NewClient(host, port string, dest uint) *Client {
	c := &Client{Base: core.NewBase(), Dest: dest}
	address := net.JoinHostPort(host, port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error("create client error: %s", err)
		return nil
	}
	c.RemoteAddress = tcpAddress
	return c
}

func (self *Client) Run() uint {
	core.RegisterService(self)
	go func() {
		for m := range self.In() {
			if m.Type == core.MSG_TYPE_CLOSE {
				self.Con.Close()
				self.Close()
			} else if m.Type == core.MSG_TYPE_NORMAL {
				cmd := m.Data[0].(int)
				if cmd == 0 { //connect
					n := m.Data[1].(int)
					self.connect(n)
				} else if cmd == 1 { //send
					if self.Con == nil {
						self.connect(2)
					}
					if self.Con != nil {
						data := m.Data[1].([]byte)
						_, err := self.outbuffer.Write(data) // self.Con.Write(data)
						if err != nil {
							log.Error("agent write msg failed: %s", err)
							self.onConError()
						}
						if self.Con != nil {
							err = self.outbuffer.Flush()
							if err != nil {
								log.Error("agent write msg failed: %s", err)
								self.onConError()
							}
						}
					}
				}
			}
		}
	}()
	return self.Id()
}
func (self *Client) connect(n int) {
	for i := 0; i < n; n++ {
		if self.Con == nil {
			var err error
			self.Con, err = net.DialTCP("tcp", nil, self.RemoteAddress)
			if err != nil {
				log.Error("client connect failed: %s", err)
			} else {
				break
			}
		}
		time.Sleep(time.Second * 2)
	}
	if self.Con == nil {
		core.Send(self.Dest, self.Id(), 0) //connect failed
	} else {
		if self.inbuffer == nil && self.outbuffer == nil {
			self.inbuffer = bufio.NewReader(self.Con)
			self.outbuffer = bufio.NewWriter(self.Con)
		} else {
			self.inbuffer.Reset(self.Con)
			self.outbuffer.Reset(self.Con)
		}
		core.Send(self.Dest, self.Id(), 1) //connect success
		go func() {
			for {
				//need to do split package.
				a := make([]byte, 8192)
				len, err := self.inbuffer.Read(a)
				if err != nil {
					log.Error("agent read msg failed: %s", err)
					self.onConError()
					break
				}
				if len > 0 {
					nt := make([]byte, len)
					copy(nt, a[:len])
					core.Send(self.Dest, self.Id(), 3, nt) //recv message
				} else {
					log.Info("agent read msg len 0")
				}
			}
		}()
	}
}
func (self *Client) onConError() {
	core.Send(self.Dest, self.Id(), 2) //recv message
	self.Con = nil
}
