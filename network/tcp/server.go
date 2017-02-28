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
	Host     string
	Port     string
	Dest     uint
	listener *net.TCPListener
}

func New(host, port string, dest uint) *Server {
	s := &Server{Host: host, Port: port}
	s.Dest = dest
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
		for {
			tcpCon, err := self.listener.AcceptTCP()
			if err != nil {
				log.Warn("tcp server: accept tcp faield %s", err)
				time.Sleep(time.Second * 3)
				continue
			}
			a := NewAgent(tcpCon, self.Dest)
			core.StartService("", 10, a)
		}
	}()

	return nil
}
