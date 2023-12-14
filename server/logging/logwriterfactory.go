package logging

import (
	"errors"
)

type LogWriterFactory struct {
	logWriters map[string]func() LogWriterIF
}

func NewLogWriterFactory() LogWriterFactoryIF {
	logWriterMap := make(map[string]func() LogWriterIF)

	logWriterMap["disk"] = NewDiskLogWriter
	logWriterMap["terminal"] = NewTerminalLogWriter

	return LogWriterFactory{logWriters: logWriterMap}
}

func (logWriterFactory LogWriterFactory) NewLogWriter(logWriterName string) (LogWriterIF, error) {
	logWriterConstructor, exists := logWriterFactory.logWriters[logWriterName]

	if !exists {
		return nil, errors.New("specified LogWriter type does not exist")
	}

	return logWriterConstructor(), nil
}
