package logging

import "github.com/SashwatAnagolum/picodb/utils"

type RequestCountLogger struct {
	AbsLogger
	counts map[string]int
}

func NewRequestCountLogger() LoggerIF {
	return &RequestCountLogger{
		AbsLogger: AbsLogger{Name: "RequestCountLogger"},
		counts:    make(map[string]int)}
}

func (logger *RequestCountLogger) Log(request *utils.PicoDBRequest) {
	_, keyExists := logger.counts[request.Key]

	if !keyExists {
		logger.counts[request.Key] = 0
	}

	logger.counts[request.Key] += 1
}

func (logger *RequestCountLogger) GetLogs() (string, error) {
	return logger.MarshalData(logger.counts)
}

func (logger *RequestCountLogger) ResetLogs() {
	logger.counts = make(map[string]int)
}
