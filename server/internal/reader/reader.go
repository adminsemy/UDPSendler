package reader

import (
	"encoding/binary"
	"log/slog"
	"net"
	"os"

	"github/adminsemy/UDPSendler/server/internal/logger"
	"github/adminsemy/UDPSendler/server/internal/server"
)

var answerOk = []byte("OK\n")

type DataFile struct {
	Part uint64 `json:"part"`
	Body []byte `json:"body"`
}

type Reader struct {
	conn *net.UDPConn
	log  *logger.Logger
}

func NewReader(s *server.Server, log *logger.Logger) *Reader {
	return &Reader{
		conn: s.GetConnection(),
		log:  log,
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
		file, err := os.Create(string(fileName))
		if err != nil {
			r.log.Error("Error create file:", "Err:", err)
			return
		}
		index := uint64(0)
		gotBite := int64(0)
		var bufferData []DataFile
		for data := range readCh {
			if data.Part < index {
				continue
			}
			if data.Part == index {
				file.Write(data.Body)
				index++
			} else {
				index = checkBuffer(bufferData, data, file, index)
				bufferData = append(bufferData, data)
			}
			gotBite += int64(len(data.Body))
			if gotBite == size {
				break
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
