package writer

import (
	"log/slog"
	"os"

	"github/adminsemy/UDPSendler/server/internal/logger"
)

type Writer struct {
	fileName string
	size     int64
	file     *os.File
	blocks   map[uint64][]byte
	log      *logger.Logger
	index    uint64
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
	}, nil
}

func (w *Writer) Close() {
	w.file.Close()
	w.log.Debug("Close file:", slog.String("name", w.fileName))
}
