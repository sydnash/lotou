package tcp

import (
	"bufio"
	"github.com/sydnash/majiang/core"
	"github.com/sydnash/majiang/log"
	"net"
	"time"
)

/*
	recieve message from network
	split message into package
	send package to dest service

	reciev other inner message and process(such as write to network; close; change dest service)

	agent has tow goroutine:
		one to read message from tcpConnector and send to dest
		one read inner message chan and process

	this also has a timeout for first data arrived after agent create.
*/

type Agent struct {
	*core.Base
	Con       *net.TCPConn
	Dest      uint
	timeout   *time.Timer
	inbuffer  *bufio.Reader
	outbuffer *bufio.Writer
}

func NewAgent(con *net.TCPConn, dest uint) *Agent {
	a := &Agent{Con: con, Dest: dest, Base: core.NewBase()}
	a.timeout = time.AfterFunc(time.Second*5, func() {
		log.Info("there is no data comming, close this connect.")
		core.Close(a.Id(), a.Id())
	})
	a.inbuffer = bufio.NewReader(a.Con)
	a.outbuffer = bufio.NewWriter(a.Con)
	return a
}

func (self *Agent) Run() {
	core.RegisterService(self)
	go func() {
		for {
			m, ok := <-self.In()
			if ok {
				if m.Type == core.MSG_TYPE_CLOSE {
					self.close()
					break
				} else if m.Type == core.MSG_TYPE_NORMAL {
					data := m.Data[0].([]byte)
					_, err := self.outbuffer.Write(data)
					if err != nil {
						log.Error("agent write msg failed: %s", err)
						core.Close(self.Id(), self.Id())
					}
					err = self.outbuffer.Flush()
					if err != nil {
						log.Error("agent write msg failed: %s", err)
						core.Close(self.Id(), self.Id())
					}
				}
			} else {
				self.close()
				break
			}
		}
	}()
	go func() {
		for {
			//need to do split package.
			a := make([]byte, 8192)
			len, err := self.inbuffer.Read(a)
			if err != nil {
				log.Error("agent read msg failed: %s", err)
				core.Close(self.Id(), self.Id())
				break
			}
			if len > 0 {
				self.timeout.Stop()
				nt := make([]byte, len)
				copy(nt, a[:len])
				core.Send(self.Dest, self.Id(), nt)
			} else {
				log.Info("agent read msg len 0")
			}
		}
	}()
}

func (self *Agent) close() {
	log.Info("close agent. %v", self.Con.RemoteAddr())
	self.Con.Close()
	self.timeout.Stop()
	self.Base.Close()
}
