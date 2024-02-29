package writer

import (
	"encoding/binary"
	"log/slog"
	"os"
	"sync"

	"github/adminsemy/UDPSendler/server/internal/logger"
)

type Writer struct {
	fileName string
	size     int64
	file     *os.File
	blocks   map[uint64][]byte
	sync.Mutex
	log   *logger.Logger
	index uint64
}

func New(fileName string, size int64, log *logger.Logger) (*Writer, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Error create file:", "Err:", err)
		return nil, err
	}
	return &Writer{
		fileName: fileName,
		size:     size,
		file:     file,
		blocks:   make(map[uint64][]byte),
		log:      log,
		index:    0,
		Mutex:    sync.Mutex{},
	}, nil
}

func (w *Writer) Close() {
	w.file.Close()
	w.log.Debug("Close file:", slog.String("name", w.fileName))
}

func (w *Writer) Write(buf []byte) {
	part := binary.LittleEndian.Uint64(buf[0:8])
	data := buf[8:]
	w.log.Debug("New data:", slog.Group(
		"part",
		slog.Uint64("part", part),
		slog.String("data", string(data)),
	))
	w.Lock()
	w.blocks[part] = data
	w.Unlock()
	go w.writeData()
}

func (w *Writer) writeData() {
	w.Lock()
	for {
		data, ok := w.blocks[w.index]
		if !ok {
			break
		}
		w.file.Write(w.blocks[w.index])
		w.size -= int64(len(data))
		delete(w.blocks, w.index)
		w.index++
	}
	w.Unlock()
}
