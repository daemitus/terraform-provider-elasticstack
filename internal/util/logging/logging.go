package logging

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
)

const (
	envLog     = "TF_LOG"      // Set to True
	envLogFile = "TF_LOG_PATH" // Set to a file
)

var validLevels = []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}

// LogLevel returns the current log level string based the environment vars.
func LogLevel() string {
	envLevel := os.Getenv(envLog)
	if envLevel == "" {
		return ""
	}

	logLevel := "TRACE"
	if isValidLogLevel(envLevel) {
		logLevel = strings.ToUpper(envLevel)
	} else {
		log.Printf("[WARN] Invalid log level: %q. Defaulting to level: TRACE. Valid levels are: %+v", envLevel, validLevels)
	}

	return logLevel
}

// IsDebugOrHigher returns whether or not the current log level is debug or trace.
func IsDebugOrHigher() bool {
	level := string(LogLevel())
	return level == "DEBUG" || level == "TRACE"
}

func isValidLogLevel(level string) bool {
	for _, l := range validLevels {
		if strings.ToUpper(level) == string(l) {
			return true
		}
	}

	return false
}
