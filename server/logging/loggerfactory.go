package logging

import (
	"errors"
)

type LoggerFactory struct {
	loggers map[string]func() LoggerIF
}

func NewLoggerFactory() LoggerFactory {
	loggerMap := make(map[string]func() LoggerIF)

	loggerMap["requestcount"] = NewRequestCountLogger
	loggerMap["requesttimes"] = NewRequestTimeLogger

	return LoggerFactory{loggers: loggerMap}
}

func (loggerFactory LoggerFactory) NewLogger(loggerName string) (LoggerIF, error) {
	loggerConstructor, exists := loggerFactory.loggers[loggerName]

	if !exists {
		return nil, errors.New("specified logger type does not exist")
	}

	return loggerConstructor(), nil
}
