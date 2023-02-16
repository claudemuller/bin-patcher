package logger

import (
	"log"
	"os"
	"strings"
)

type Log struct {
	logger *log.Logger
	logs   []string
}

func NewLogger() *Log {
	return &Log{
		logger: log.New(os.Stdout, "bin-patcher", 0),
	}
}

func (l *Log) Log(msg string) {
	l.logs = append(l.logs, msg)

	l.logger.Printf("%+v", msg)
}

func (l *Log) GetLogs() string {
	return strings.Join(l.logs, "\n")
}
