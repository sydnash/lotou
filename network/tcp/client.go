package tcp

import (
	"bufio"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"net"
	"time"
)

type Client struct {
	*core.Skeleton
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
	c := &Client{Skeleton: core.NewSkeleton(), Dest: dest}
	address := net.JoinHostPort(host, port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error("tcp resolve failed.")
		return nil
	}
	c.RemoteAddress = tcpAddress
	return c
}

func (c *Client) OnInit() {
}

func (c *Client) OnDestroy() {
	if c.Con != nil {
		c.Con.Close()
	}
}
func (c *Client) onConnect(n int) {
	c.connect(n)
}
func (c *Client) onSend(src uint, param ...interface{}) {
	if c.Con == nil {
		c.connect(2)
	}
	if c.Con != nil {
		data := param[0].([]byte)
		c.Con.SetWriteDeadline(time.Now().Add(time.Second * 20))
		_, err := c.outbuffer.Write(data)
		if err != nil {
			log.Error("agent write msg failed: %s", err)
			c.onConError()
		}
		if c.Con != nil {
			err = c.outbuffer.Flush()
			if err != nil {
				log.Error("agent write msg failed: %s", err)
				c.onConError()
			}
		}
	}
}

func (c *Client) OnNormalMSG(src uint, param ...interface{}) {
	cmd := param[0].(int)
	param = param[1:]
	if cmd == CLIENT_CMD_CONNECT { //connect
		n := param[0].(int)
		c.onConnect(n)
	} else if cmd == CLIENT_CMD_SEND { //send
		c.onSend(src, param...)
	}
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
		self.RawSend(self.Dest, core.MSG_TYPE_SOCKET, CLIENT_CONNECT_FAILED) //connect failed
	} else {
		if self.inbuffer == nil && self.outbuffer == nil {
			self.inbuffer = bufio.NewReader(self.Con)
			self.outbuffer = bufio.NewWriter(self.Con)
		} else {
			self.inbuffer.Reset(self.Con)
			self.outbuffer.Reset(self.Con)
		}
		self.RawSend(self.Dest, core.MSG_TYPE_SOCKET, CLIENT_CONNECTED) //connect success
		go func() {
			for {
				//split package
				pack, err := Subpackage(self.inbuffer)
				if err != nil {
					log.Error("agent read msg failed: %s", err)
					self.onConError()
					break
				}
				self.RawSend(self.Dest, core.MSG_TYPE_SOCKET, CLIENT_DATA, pack) //recv message
			}
		}()
	}
}
func (self *Client) onConError() {
	self.RawSend(self.Dest, core.MSG_TYPE_NORMAL, CLIENT_DISCONNECTED) //disconnected
	self.OnDestroy()
	self.Con = nil
}
