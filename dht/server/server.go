package server

import "net"

type Server struct {
	conn *net.UDPConn
}

type Options struct {
	Addr string
}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
}
