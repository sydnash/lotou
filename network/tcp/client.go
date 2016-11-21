package tcp

import (
	"bufio"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
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

const (
	CLIENT_CONNECT_FAILED = iota
	CLIENT_CONNECTED
	CLIENT_DISCONNECTED
	CLIENT_DATA
)
const (
	CLIENT_CMD_CONNECT = iota
	CLIENT_CMD_SEND
)

func NewClient(host, port string, dest uint) *Client {
	c := &Client{Base: core.NewBase(), Dest: dest}
	address := net.JoinHostPort(host, port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error("create client error: %s", err)
		return nil
	}
	c.RemoteAddress = tcpAddress

	c.SetSelf(c)
	c.RegisterBaseCB(core.MSG_TYPE_CLOSE, (*Client).Close, true)
	c.RegisterBaseCB(core.MSG_TYPE_NORMAL, (*Client).Normal, true)
	return c
}

func (self *Client) Close(dest, src uint, encodeType string) {
	self.Base.Close()
	if self.Con != nil {
		self.Con.Close()
	}
}
func (self *Client) Normal(dest, src uint, encodeType string, cmd int, param ...interface{}) {
	if cmd == CLIENT_CMD_CONNECT { //connect
		n := param[0].(int)
		self.connect(n)
	} else if cmd == CLIENT_CMD_SEND { //send
		if self.Con == nil {
			self.connect(2)
		}
		if self.Con != nil {
			data := param[0].([]byte)
			_, err := self.outbuffer.Write(data)
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

func (self *Client) Run() uint {
	core.RegisterService(self)
	go func() {
		for m := range self.In() {
			self.DispatchM(m)
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
		core.SendSocket(self.Dest, self.Id(), CLIENT_CONNECT_FAILED) //connect failed
	} else {
		if self.inbuffer == nil && self.outbuffer == nil {
			self.inbuffer = bufio.NewReader(self.Con)
			self.outbuffer = bufio.NewWriter(self.Con)
		} else {
			self.inbuffer.Reset(self.Con)
			self.outbuffer.Reset(self.Con)
		}
		core.SendSocket(self.Dest, self.Id(), CLIENT_CONNECTED) //connect success
		go func() {
			for {
				//split package
				pack, err := Subpackage(self.inbuffer)
				if err != nil {
					log.Error("agent read msg failed: %s", err)
					self.onConError()
					break
				}
				core.SendSocket(self.Dest, self.Id(), CLIENT_DATA, pack) //recv message
			}
		}()
	}
}
func (self *Client) onConError() {
	core.SendSocket(self.Dest, self.Id(), CLIENT_DISCONNECTED) //disconnected
	self.Con = nil
}
