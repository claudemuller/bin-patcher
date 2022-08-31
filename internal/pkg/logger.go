package pkg

import (
	"log"
	"os"
	"strings"
)

type Log struct {
	logger *log.Logger
	logs   []string
}

func newLogger() *Log {
	return &Log{
		logger: log.New(os.Stdout, "bin-patcher", 0),
	}
}

func (l *Log) log(msg string) {
	l.logs = append(l.logs, msg)

	l.logger.Printf("%+v", msg)
}

func (l *Log) getLogs() string {
	return strings.Join(l.logs, "\n")
}
