package lsp

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLog(logLevel int, logFile string) {
	if logFile == "" {
		// set logger to zap no op logger
		logConfig := zap.NewDevelopmentConfig()
		logger, _ = logConfig.Build()
		return
	}
	var err error
	logConfig := zap.NewDevelopmentConfig()
	logConfig.OutputPaths = []string{logFile}
	logConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(logLevel))
	logger, err = logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
}
