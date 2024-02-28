package main

import (
	"log/slog"

	"github/adminsemy/UDPSendler/server/internal/logger"
)

func main() {
	logger := logger.NewLogger()
	logger.Debug("Test")
	logger.SetLevel(slog.LevelInfo)
	logger.Debug("Test2")
	logger.Info("Test3", slog.String("key", "value"))
	slog.Info("Test4")
}
