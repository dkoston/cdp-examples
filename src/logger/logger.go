package logger

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

var logger *logging.Logger

// Get gives back a new logging.Logger instance configured to our specs
// or the existing instance
func Get() *logging.Logger {
	if logger != nil {
		return logger
	}

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{color:reset} %{message}`,
	)

	var log = logging.MustGetLogger("")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	// Get log level from ENV (we only support DEBUG and INFO (default))
	logLevel := os.Getenv("LOG_LEVEL")
	logLevel = strings.Trim(logLevel, "\n ")
	level := logging.INFO

	if logLevel == "debug" || logLevel == "DEBUG" {
		level = logging.DEBUG
	}

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)

	logger = log

	return logger
}
