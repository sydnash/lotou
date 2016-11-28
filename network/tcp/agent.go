package tcp

import (
	"bufio"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"net"
	"sync"
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
	Con         *net.TCPConn
	Dest        uint
	timeout     *time.Timer
	inbuffer    *bufio.Reader
	outbuffer   *bufio.Writer
	ticker      *time.Ticker
	buffer      [][]byte
	bufferMutxt sync.Mutex
}

const (
	AGENT_CLOSED = iota
	AGENT_DATA
	AGENT_ARRIVE
)
const (
	AGENT_CMD_SEND = iota
)

func NewAgent(con *net.TCPConn, dest uint) *Agent {
	a := &Agent{Con: con, Dest: dest, Base: core.NewBase()}
	a.timeout = time.AfterFunc(time.Second*5, func() {
		log.Info("there is no data comming, close this connect.")
		core.Close(a.Id(), a.Id())
	})
	a.inbuffer = bufio.NewReader(a.Con)
	a.outbuffer = bufio.NewWriter(a.Con)
	a.buffer = make([][]byte, 0, 100)
	return a
}

func (self *Agent) Run() {
	core.RegisterService(self)
	core.SendSocket(self.Dest, self.Id(), AGENT_ARRIVE) //recv message
	self.ticker = time.NewTicker(time.Millisecond * 100)
	go func() {
	OUT:
		for {
			<-self.ticker.C
		SELECT_LOOP:
			for {
				select {
				case m, ok := <-self.In():
					if ok {
						if m.Type == core.MSG_TYPE_CLOSE {
							self.close()
							break OUT
						} else if m.Type == core.MSG_TYPE_NORMAL {
							cmd := m.Data[0].(int)
							if cmd == AGENT_CMD_SEND {
								data := m.Data[1].([]byte)
								_, err := self.outbuffer.Write(data)
								if err != nil {
									log.Error("agent write msg failed: %s", err)
									self.onConnectError()
								}
								err = self.outbuffer.Flush()
								if err != nil {
									log.Error("agent write msg failed: %s", err)
									self.onConnectError()
								}
							}
						}
					} else {
						self.close()
						break OUT
					}
				default:
					break SELECT_LOOP
				}
			}
			self.bufferMutxt.Lock()
			for _, pack := range self.buffer {
				core.SendSocket(self.Dest, self.Id(), AGENT_DATA, pack)
			}
			self.buffer = self.buffer[0:0]
			self.bufferMutxt.Unlock()
		}
	}()
	go func() {
		for {
			pack, err := Subpackage(self.inbuffer)
			if err != nil {
				log.Error("agent read msg failed: %s", err)
				self.onConnectError()
				break
			}
			if self.timeout != nil {
				self.timeout.Stop()
				self.timeout = nil
			}
			self.bufferMutxt.Lock()
			self.buffer = append(self.buffer, pack)
			self.bufferMutxt.Unlock()
		}
	}()
}

func (self *Agent) onConnectError() {
	core.SendSocket(self.Dest, self.Id(), AGENT_CLOSED)
	core.Close(self.Id(), self.Id())
}
func (self *Agent) close() {
	log.Info("close agent. %v", self.Con.RemoteAddr())
	self.Con.Close()
	if self.timeout != nil {
		self.timeout.Stop()
	}
	self.Base.Close()
	self.ticker.Stop()
}
