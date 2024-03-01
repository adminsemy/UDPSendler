package reader

import (
	"encoding/binary"
	"log/slog"
	"net"

	"github/adminsemy/UDPSendler/server/internal/logger"
	"github/adminsemy/UDPSendler/server/internal/server"
	"github/adminsemy/UDPSendler/server/internal/writer"
)

var answerOk = []byte("OK\n")

type Reader struct {
	conn   *net.UDPConn
	log    *logger.Logger
	dataCh chan []byte
}

func NewReader(s *server.Server, log *logger.Logger) *Reader {
	return &Reader{
		conn:   s.GetConnection(),
		log:    log,
		dataCh: make(chan []byte),
	}
}

func (r *Reader) Read() {
	for {
		var buf [264]byte
		n, addr, err := r.conn.ReadFromUDP(buf[0:])
		if err != nil {
			r.log.Error("Error from read data:", "Err:", err)
			return
		}
		size := int64(binary.LittleEndian.Uint64(buf[0:8]))
		fileName := buf[8:n]
		r.log.Info("New file:", slog.String("name", string(fileName)), slog.Int64("size", size))
		r.conn.WriteToUDP(answerOk, addr)
		file, err := writer.New(string(fileName), size, r.log)
		if err != nil {
			r.log.Error("Error create file:", "Err:", err)
			return
		}
	endData:
		for {
			select {
			case data := <-r.dataCh:
				file.Write(data)
			default:
				break endData
			}
		}
		r.log.Info("File saved:", slog.String("name", string(fileName)))
		file.Close()
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
