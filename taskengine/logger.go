package taskengine

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string)
	Infof(format string, args ...any)

	Warn(msg string)
	Warnf(format string, args ...any)

	Error(msg string)
	Errorf(format string, args ...any)
}

type stdLogger struct {
	prefixInfo  string
	prefixWarn  string
	prefixError string
}

func (l *stdLogger) Info(msg string) { log.Println(l.prefixInfo + msg) }

func (l *stdLogger) Infof(format string, args ...any) { log.Printf(l.prefixInfo+format, args...) }

func (l *stdLogger) Warn(msg string) { log.Println(l.prefixWarn + msg) }

func (l *stdLogger) Warnf(format string, args ...any) { log.Printf(l.prefixWarn+format, args...) }

func (l *stdLogger) Error(msg string) { log.Println(l.prefixError + msg) }

func (l *stdLogger) Errorf(format string, args ...any) { log.Printf(l.prefixError+format, args...) }

func newLogger(module string) Logger {
	if module == "" {
		return &stdLogger{
			prefixInfo:  "[INFO] ",
			prefixWarn:  "[WARN] ",
			prefixError: "[ERROR] ",
		}
	}
	return &stdLogger{
		prefixInfo:  fmt.Sprintf("[%s][INFO] ", module),
		prefixWarn:  fmt.Sprintf("[%s][WARN] ", module),
		prefixError: fmt.Sprintf("[%s][ERROR] ", module),
	}
}
