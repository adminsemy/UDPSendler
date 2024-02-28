package reader

import (
	"encoding/binary"
	"log/slog"
	"net"

	"github/adminsemy/UDPSendler/server/internal/server"
)

type DataFile struct {
	Part uint64 `json:"part"`
	Body []byte `json:"body"`
}

type Reader struct {
	conn *net.UDPConn
}

func NewReader(s *server.Server) *Reader {
	return &Reader{
		conn: s.GetConnection(),
	}
}

func (r *Reader) readData(readCh chan<- DataFile) {
	var buf [507]byte
	var data DataFile
	for {
		n, addr, err := r.conn.ReadFromUDP(buf[0:])
		if err != nil {
			slog.Error("Error from read data:", "Err:", err)
			return
		}
		data.Part = binary.LittleEndian.Uint64(buf[0:8])
		data.Body = buf[8:n]
		slog.Info("New part:", "part:", string(data.Body))
		r.conn.WriteToUDP([]byte("OK\n"), addr)
		readCh <- data
	}
}
