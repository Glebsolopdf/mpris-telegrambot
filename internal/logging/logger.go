package logging

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	debug  bool
	logger *log.Logger
}

func New(level string, writers ...io.Writer) Logger {
	writer := io.Writer(os.Stderr)
	if len(writers) > 0 && writers[0] != nil {
		writer = writers[0]
	}
	return Logger{
		debug:  level == "debug",
		logger: log.New(writer, "", log.LstdFlags),
	}
}

func (l Logger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l Logger) Debugf(format string, v ...any) {
	if l.debug {
		l.logger.Printf(format, v...)
	}
}
