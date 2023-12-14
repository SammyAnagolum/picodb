package logging

import (
	"time"

	"github.com/SashwatAnagolum/picodb/utils"
)

type RequestTimeLogger struct {
	AbsLogger
	times []time.Time
	Name  string
}

func NewRequestTimeLogger() LoggerIF {
	return &RequestTimeLogger{
		AbsLogger: AbsLogger{Name: "RequestTimeLogger"},
		times:     make([]time.Time, 0, 1024)}
}

func (logger *RequestTimeLogger) Log(request *utils.PicoDBRequest) {
	logger.times = append(logger.times, time.Now())
}

func (logger *RequestTimeLogger) GetLogs() (string, error) {
	return logger.MarshalData(logger.times)
}

func (logger *RequestTimeLogger) ResetLogs() {
	logger.times = make([]time.Time, 0, 1024)
}
