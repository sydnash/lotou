package tcp

import (
	"bufio"
	"bytes"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"net"
	"sync/atomic"
	"time"
)

const (
	CLIENT_STATUS_NOT_CONNECT = iota
	CLIENT_STATUS_CONNECTING
	CLIENT_STATUS_CONNECTED
)

type Client struct {
	*core.Skeleton
	Con             *net.TCPConn
	RemoteAddress   *net.TCPAddr
	hostService     core.ServiceID
	inbuffer        *bufio.Reader
	outbuffer       *bufio.Writer
	status          int32
	bufferForOutMsg *bytes.Buffer
	isNeedExit      bool
}

func NewClient(host, port string, hostID core.ServiceID) *Client {
	c := &Client{
		Skeleton:    core.NewSkeleton(0),
		hostService: hostID,
	}
	address := net.JoinHostPort(host, port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error("tcp resolve failed.")
		return nil
	}
	c.RemoteAddress = tcpAddress
	c.status = CLIENT_STATUS_NOT_CONNECT
	c.bufferForOutMsg = bytes.NewBuffer([]byte{})
	return c
}

func (c *Client) OnInit() {
}

func (c *Client) OnDestroy() {
	c.isNeedExit = true
	if c.Con != nil {
		c.Con.Close()
	}
}
func (c *Client) onConnect(n int) {
	c.connect(n)
}

func (c *Client) sendBufferOutMsgAndData(data []byte) {
	if c.Con != nil {
		c.Con.SetWriteDeadline(time.Now().Add(time.Second * 20))
		if c.bufferForOutMsg.Len() > 0 {
			_, err := c.outbuffer.Write(c.bufferForOutMsg.Bytes())
			if err != nil {
				log.Error("client onSend tmp out msg err: %s", err)
				c.onConError()
			}
			c.bufferForOutMsg.Reset()
		}
		if data != nil {
			_, err := c.outbuffer.Write(data)
			if err != nil {
				log.Error("client onSend writebuff err: %s", err)
				c.onConError()
			}
		}
		if c.Con != nil {
			err := c.outbuffer.Flush()
			if err != nil {
				log.Error("client onSend err: %s", err)
				c.onConError()
			}
		}
	}
}

func (c *Client) onSend(src core.ServiceID, param ...interface{}) {
	if c.status != CLIENT_STATUS_CONNECTED {
		if c.status == CLIENT_STATUS_NOT_CONNECT {
			go c.connect(-1)
		}
		data := param[0].([]byte)
		defer func() {
			if err := recover(); err != nil {
				if err == bytes.ErrTooLarge {
					log.Error("client out msg buffer is too large, we will reset it.")
					c.bufferForOutMsg.Reset()
				} else {
					panic(err)
				}
			}
		}()
		c.bufferForOutMsg.Write(data)
		return
	}
	c.sendBufferOutMsgAndData(param[0].([]byte))
}

func (c *Client) OnNormalMSG(msg *core.Message) {
	src := msg.Src
	cmd := msg.Cmd
	param := msg.Data
	switch cmd {
	case CLIENT_CMD_CONNECT:
		n := param[0].(int)
		c.onConnect(n)
	case CLIENT_CMD_SEND:
		c.onSend(src, param...)
	}
}

func (c *Client) connect(n int) {
	if !atomic.CompareAndSwapInt32(&c.status, CLIENT_STATUS_NOT_CONNECT, CLIENT_STATUS_CONNECTING) {
		return
	}
	i := 0
	for {
		if c.isNeedExit {
			return
		}
		if c.Con == nil {
			var err error
			c.Con, err = net.DialTCP("tcp", nil, c.RemoteAddress)
			if err != nil {
				log.Error("client connect failed: %s", err)
			} else {
				break
			}
		}
		time.Sleep(time.Second * 2)
		i++
		if n > 0 && i >= n {
			break
		}
	}
	if c.Con == nil {
		c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_CONNECT_FAILED) //connect failed
	} else {
		if c.inbuffer == nil && c.outbuffer == nil {
			c.inbuffer = bufio.NewReader(c.Con)
			c.outbuffer = bufio.NewWriter(c.Con)
		} else {
			c.inbuffer.Reset(c.Con)
			c.outbuffer.Reset(c.Con)
		}
		c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_CONNECTED) //connect success
		c.sendBufferOutMsgAndData(nil)
		go func() {
			for {
				//split package
				pack, err := Subpackage(c.inbuffer)
				if err != nil {
					log.Error("client read msg failed: %s", err)
					c.onConError()
					break
				}
				c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_DATA, pack) //recv message
			}
		}()
	}
	atomic.StoreInt32(&c.status, CLIENT_STATUS_CONNECTED)
}

func (c *Client) onConError() {
	c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_DISCONNECTED) //disconnected
	c.OnDestroy()
	if c.Con != nil {
		c.Con.Close()
	}
	c.Con = nil
	atomic.StoreInt32(&c.status, CLIENT_STATUS_NOT_CONNECT)
}

func (c *Client) sendToHost(msgType core.MsgType, cmd core.CmdType, data ...interface{}) {
	c.RawSend(c.hostService, msgType, cmd, data...)
}
