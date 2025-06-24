package taskengine

import "log"

type Logger interface {
	Info(msg string)
	Infof(format string, args ...any)

	Warn(msg string)
	Warnf(format string, args ...any)

	Error(msg string)
	Errorf(format string, args ...any)
}

type stdLogger struct{}

func (l *stdLogger) Info(msg string) { log.Println("[INFO] " + msg) }

func (l *stdLogger) Infof(format string, args ...any) { log.Printf("[INFO] "+format, args...) }

func (l *stdLogger) Warn(msg string) { log.Println("[WARN] " + msg) }

func (l *stdLogger) Warnf(format string, args ...any) { log.Printf("[WARN] "+format, args...) }

func (l *stdLogger) Error(msg string) { log.Println("[ERROR] " + msg) }

func (l *stdLogger) Errorf(format string, args ...any) { log.Printf("[ERROR] "+format, args...) }
