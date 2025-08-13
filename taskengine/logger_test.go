package taskengine

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDefaultLoggerFactory(t *testing.T) {
	tests := []struct {
		name          string
		module        string
		expectedInfo  string
		expectedWarn  string
		expectedError string
	}{
		{
			name:          "empty module",
			module:        "",
			expectedInfo:  "[INFO] ",
			expectedWarn:  "[WARN] ",
			expectedError: "[ERROR] ",
		},
		{
			name:          "with module",
			module:        "test-module",
			expectedInfo:  "[test-module][INFO] ",
			expectedWarn:  "[test-module][WARN] ",
			expectedError: "[test-module][ERROR] ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := DefaultLoggerFactory(tc.module)
			stdLogger, ok := logger.(*stdLogger)
			if !ok {
				t.Errorf("expected *stdLogger, got %T", logger)
			}

			if got := stdLogger.prefixInfo; got != tc.expectedInfo {
				t.Errorf("expected info prefix %s, got %s", tc.expectedInfo, got)
			}

			if got := stdLogger.prefixWarn; got != tc.expectedWarn {
				t.Errorf("expected warn prefix %s, got %s", tc.expectedWarn, got)
			}

			if got := stdLogger.prefixError; got != tc.expectedError {
				t.Errorf("expected error prefix %s, got %s", tc.expectedError, got)
			}
		})
	}
}

func TestStdLoggerPrefixes(t *testing.T) {
	logger := &stdLogger{
		prefixInfo:  "[TEST][INFO] ",
		prefixWarn:  "[TEST][WARN] ",
		prefixError: "[TEST][ERROR] ",
	}

	// Test that the logger implements the Logger interface
	var _ Logger = logger

	// Test prefixes are set correctly
	if want := "[TEST][INFO] "; logger.prefixInfo != want {
		t.Errorf("expected info prefix %s, got %s", want, logger.prefixInfo)
	}

	if want := "[TEST][WARN] "; logger.prefixWarn != want {
		t.Errorf("expected warn prefix %s, got %s", want, logger.prefixWarn)
	}

	if want := "[TEST][ERROR] "; logger.prefixError != want {
		t.Errorf("expected error prefix %s, got %s", want, logger.prefixError)
	}
}

func TestLoggerFactoryType(t *testing.T) {
	// Test that LoggerFactory is a function type that works as expected
	var factory LoggerFactory = DefaultLoggerFactory

	logger := factory("test")
	if logger == nil {
		t.Error("expected logger to be created, got nil")
	}

	// Test that the returned logger implements the Logger interface
	var _ Logger = logger
}

func TestMultipleLoggersFromFactory(t *testing.T) {
	factory := DefaultLoggerFactory

	logger1 := factory("module1")
	logger2 := factory("module2")

	// Should be different instances
	if logger1 == logger2 {
		t.Error("expected different logger instances, got same")
	}

	// Should have different prefixes
	std1, ok1 := logger1.(*stdLogger)
	std2, ok2 := logger2.(*stdLogger)

	if !ok1 || !ok2 {
		t.Fatal("expected both loggers to be *stdLogger")
	}

	if std1.prefixInfo == std2.prefixInfo {
		t.Error("expected different info prefixes for different modules")
	}

	if !strings.Contains(std1.prefixInfo, "module1") {
		t.Errorf("expected logger1 prefix to contain 'module1', got %s", std1.prefixInfo)
	}

	if !strings.Contains(std2.prefixInfo, "module2") {
		t.Errorf("expected logger2 prefix to contain 'module2', got %s", std2.prefixInfo)
	}
}

func TestStdLoggerMethods(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	logger := &stdLogger{
		prefixInfo:  "[TEST][INFO] ",
		prefixWarn:  "[TEST][WARN] ",
		prefixError: "[TEST][ERROR] ",
	}

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			name:     "Info",
			logFunc:  func() { logger.Info("test message") },
			expected: "[TEST][INFO] test message",
		},
		{
			name:     "Infof",
			logFunc:  func() { logger.Infof("test %s", "formatted") },
			expected: "[TEST][INFO] test formatted",
		},
		{
			name:     "Warn",
			logFunc:  func() { logger.Warn("warning message") },
			expected: "[TEST][WARN] warning message",
		},
		{
			name:     "Warnf",
			logFunc:  func() { logger.Warnf("warning %d", 123) },
			expected: "[TEST][WARN] warning 123",
		},
		{
			name:     "Error",
			logFunc:  func() { logger.Error("error message") },
			expected: "[TEST][ERROR] error message",
		},
		{
			name:     "Errorf",
			logFunc:  func() { logger.Errorf("error %v", "test") },
			expected: "[TEST][ERROR] error test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.logFunc()

			output := strings.TrimSpace(buf.String())

			if !strings.Contains(output, tc.expected) {
				t.Errorf("expected output to contain %q, got %q", tc.expected, output)
			}
		})
	}
}
