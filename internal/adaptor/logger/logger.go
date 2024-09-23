package logger

import (
	"github.com/hashicorp/go-hclog"
)

type HclogAdapter struct {
	logger hclog.Logger
}

func NewHclogAdapter() *HclogAdapter {
	return &HclogAdapter{
		logger: hclog.New(&hclog.LoggerOptions{
			Name:  "sarva",
			Level: hclog.LevelFromString("DEBUG"),
		}),
	}
}

func (l *HclogAdapter) Log(level string, message string) {
	switch level {
	case "DEBUG":
		l.logger.Debug(message)
	case "INFO":
		l.logger.Info(message)
	case "ERROR":
		l.logger.Error(message)
	default:
		l.logger.Info(message)
	}
}
