package gospace

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// environment variable to control the verbosity level
	LOGGING_ENV = "GOSPACE_VERBOSE"
)

const (
	// logging is disabled
	LOG_OFF LogLevel = iota
	// logging channel for unrecoverable errors
	LOG_FATAL = iota
	// logging channel for recoverable errors
	LOG_ERROR = iota
	// logging channel for potential problems
	LOG_WARN = iota
	// logging channel for general information
	LOG_INFO = iota
	// logging channel for detailed information
	LOG_DEBUG = iota
	// logging channel for verbose information about the current state
	LOG_TRACE = iota
)

var (
	// the current verbosity level
	LOG_LEVEL LogLevel = LOG_OFF
)

// verbosity level
type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LOG_TRACE:
		return "TRACE"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	default:
		return "OFF"
	}
}

// increase the value of the level by the given range.
// the result is clamped between the known level values.
func (l *LogLevel) Increase(diff int) {
	*l = ParseLogLevel(int(*l) + diff)
}

func init() {
	env := os.Getenv(LOGGING_ENV)

	if 0 < len(env) {
		if level, err := strconv.Atoi(env); nil == err {
			LOG_LEVEL = ParseLogLevel(level)
		} else {
			LOG_LEVEL = ParseLogName(env)
		}

		D("verbosity set to", LOG_LEVEL, "via environment")
	}
}

// convert the numeric value into a log level instance.
// invalid values are returned as LOG_OFF
func ParseLogLevel(value int) LogLevel {
	cast := LogLevel(value)
	// cheap sanity check
	if LOG_TRACE > cast || LOG_FATAL < cast {
		return cast
	}

	return LOG_OFF
}

// convert the string representation into an instance of log level.
// the conversion is case insensitive. the default result is LOG_OFF
func ParseLogName(value string) LogLevel {
	switch strings.ToLower(value) {
	case "trace":
		fallthrough
	case "spam":
		return LOG_TRACE

	case "insight":
		fallthrough
	case "debug":
		return LOG_DEBUG

	case "info":
		fallthrough
	case "information":
		fallthrough
	case "informative":
		return LOG_INFO

	case "warn":
		fallthrough
	case "warning":
		return LOG_WARN

	case "err":
		fallthrough
	case "error":
		return LOG_ERROR

	case "fatal":
		fallthrough
	case "panic":
		return LOG_FATAL

	default:
		return LOG_OFF
	}
}

// log a message with verbosity TRACE on stdout.
// if the current log level is lower, nothing is written.
func T(message ...interface{}) {
	if LOG_TRACE <= LOG_LEVEL {
		fmt.Println(message...)
	}
}

// log a message with verbosity DEBUG on stdout.
// if the current log level is lower, nothing is written.
func D(message ...interface{}) {
	if LOG_DEBUG <= LOG_LEVEL {
		fmt.Println(message...)
	}
}

// log a message with verbosity INFO on stdout.
// if the current log level is lower, nothing is written.
func I(message ...interface{}) {
	if LOG_INFO <= LOG_LEVEL {
		fmt.Println(message...)
	}
}

// log a message with verbosity WARN on stdout.
// if the current log level is lower, nothing is written.
func W(message ...interface{}) {
	if LOG_WARN <= LOG_LEVEL {
		fmt.Println(message...)
	}
}

// log a message with verbosity ERROR on stdout.
// if the current log level is lower, nothing is written.
func E(message ...interface{}) {
	if LOG_ERROR <= LOG_LEVEL {
		fmt.Println(message...)
	}
}

// log a message with verbosity FATAL on stdout.
// if the current log level is lower, nothing is written.
func F(message ...interface{}) {
	if LOG_FATAL <= LOG_LEVEL {
		fmt.Println(message...)
	}
}
