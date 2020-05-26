package agollo

import (
	"log"
	"os"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func newLogger() Logger {
	return &logger{
		log: log.New(os.Stdout, "[agollo] ", log.LstdFlags),
	}
}

type logger struct {
	log *log.Logger
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Printf("[INFO] "+format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Printf("[ERROR] "+format, args...)
}
