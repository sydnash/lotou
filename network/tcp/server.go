package tcp

/*
server listen on a tcp port
when a tcp connect comming in
create a agent
*/
import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"net"
	"time"
)

type Server struct {
	Host        string
	Port        string
	hostService core.ServiceID
	listener    *net.TCPListener
}

const (
	TCPServerClosed = "TCPServerClosed"
)

func NewServer(host, port string, hsID core.ServiceID) *Server {
	s := &Server{
		Host:        host,
		Port:        port,
		hostService: hsID,
	}
	return s
}

func (self *Server) Listen() error {
	address := net.JoinHostPort(self.Host, self.Port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error("tcp server: resolve tcp address failed, %s", err)
		return err
	}
	self.listener, err = net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Error("tcp server: listen tcp failed %s", err)
		return err
	}

	go func() {
		var tempDelay time.Duration
		for {
			tcpCon, err := self.listener.AcceptTCP()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Warn("tcp server: accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				log.Error("tcp server: accept err %s, server closed.", err)
				core.Send(self.hostService, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_NO, TCPServerClosed)
				break
			}
			a := NewAgent(tcpCon, self.hostService)
			core.StartService(&core.ModuleParam{
				N: "",
				M: a,
				L: 0,
			})
		}
	}()

	return nil
}

func (self *Server) Close() {
	if self.listener != nil {
		self.listener.Close()
	}
}
