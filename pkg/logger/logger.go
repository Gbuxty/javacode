package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed start new looger:%w", err)
	}
	return &Logger{Logger: logger}, nil
}
