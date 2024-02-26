package main

import (
	"log/slog"

	"github/adminsemy/UDPSendler/server/internal/logger"
)

func main() {
	logger := logger.NewLogger()
	logger.Debug("Test")
	logger.SetLevel(slog.LevelInfo)
	logger.AddSource(false)
	logger.Debug("Test2")
	logger.Info("Test2")
}
