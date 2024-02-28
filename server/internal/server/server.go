package server

import (
	"net"
	"os"

	"github/adminsemy/UDPSendler/server/internal/logger"
)

type Server struct {
	address string
	conn    *net.UDPConn
	logger  *logger.Logger
}

func NewServer(address string, log *logger.Logger) *Server {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Error("Error resolve address:", "Err:", err)
		os.Exit(1)
	}
	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error("Error listen address:", "Err:", err)
		os.Exit(1)
	}
	return &Server{address: address, conn: conn, logger: log}
}

func (s *Server) Close() {
	s.conn.Close()
}

func (s *Server) GetConnection() *net.UDPConn {
	return s.conn
}
