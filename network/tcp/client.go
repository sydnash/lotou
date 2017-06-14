package tcp

import (
	"bufio"
	"bytes"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/helper"
	"github.com/sydnash/lotou/log"
	"net"
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
	inbuffer        []byte
	outbuffer       *bufio.Writer
	parseCache      *ParseCache
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
	//status modify on main goroutinue
	c.status = CLIENT_STATUS_NOT_CONNECT
	c.bufferForOutMsg = bytes.NewBuffer([]byte{})
	return c
}

func (c *Client) OnInit() {
	c.inbuffer = make([]byte, DEFAULT_RECV_BUFF_LEN)
}

func (c *Client) OnDestroy() {
	c.isNeedExit = true
	if c.Con != nil {
		c.Con.Close()
		c.Con = nil
	}
}
func (c *Client) onConnect(n int) {
	c.connect(n)
}

func (c *Client) sendBufferOutMsgAndData(data []byte) {
	if c.Con != nil {
		c.Con.SetWriteDeadline(time.Now().Add(time.Second * 5))
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
	//status check and modify on main goroutinue
	if c.status != CLIENT_STATUS_CONNECTED {
		if c.status == CLIENT_STATUS_NOT_CONNECT {
			c.status = CLIENT_STATUS_CONNECTING
			c.isNeedExit = false
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
	//sendBufferOutMsgAndData on main goroutinue
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
	case CLIENT_SELF_CONNECTED:
		//status modify on main goroutinue
		c.status = CLIENT_STATUS_CONNECTED
		//sendBufferOutMsgAndData on main goroutinue
		c.sendBufferOutMsgAndData(nil)
	case CLIENT_SELF_DISCONNECTED:
		//status modify on goroutinue
		c.status = CLIENT_STATUS_NOT_CONNECT
		c.OnDestroy()
	}
}

func (c *Client) connect(n int) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("recover: stack: %v\n, %v", helper.GetStack(), err)
		}
	}()
	i := 0
	for {
		if c.isNeedExit {
			return
		}
		if c.Con == nil {
			var err error
			c.Con, err = net.DialTCP("tcp", nil, c.RemoteAddress)
			if err != nil {
				//log.Error("client connect failed: %s", err)
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
		if c.outbuffer == nil {
			c.outbuffer = bufio.NewWriter(c.Con)
			c.parseCache = &ParseCache{}
		} else {
			c.outbuffer.Reset(c.Con)
			c.parseCache.reset()
		}
		c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_CONNECTED) //connect success
		//send normal msg to tell self tcp is connected.
		c.RawSend(c.Id, core.MSG_TYPE_NORMAL, CLIENT_SELF_CONNECTED)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Error("recover: stack: %v\n, %v", helper.GetStack(), err)
				}
			}()
			for {
				//split package
				pack, err := Subpackage(c.inbuffer, c.Con, c.parseCache)
				if err != nil {
					log.Error("client read msg failed: %s", err)
					c.onConError()
					break
				}
				for _, v := range pack {
					c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_DATA, v) //recv message
				}
			}
		}()
	}
}

func (c *Client) onConError() {
	c.sendToHost(core.MSG_TYPE_SOCKET, CLIENT_DISCONNECTED) //disconnected
	//send normal msg to tell self tcp is diconnected.
	c.RawSend(c.Id, core.MSG_TYPE_NORMAL, CLIENT_SELF_DISCONNECTED)
}

func (c *Client) sendToHost(msgType core.MsgType, cmd core.CmdType, data ...interface{}) {
	c.RawSend(c.hostService, msgType, cmd, data...)
}
