package logging

import "github.com/SashwatAnagolum/picodb/utils"

type NullLogger struct {
	AbsLogger
}

func NewNullLogger() LoggerIF {
	return &NullLogger{AbsLogger{Name: "NullLogger"}}
}

func (logger *NullLogger) Log(request *utils.PicoDBRequest) {
}

func (logger *NullLogger) GetLogs() (string, error) {
	return "", nil
}

func (logger *NullLogger) ResetLogs() {
}
