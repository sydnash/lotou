package network

import (
	"fmt"
	"net"
)

type struct TCPServer {
	Host string
	Port string

	listener *net.TCPListener
}

func newTCPServer(host, port string) *TCPServer {
	return &TCPServer(host, port)
}

func (self *TCPServer) Listen() {
}
