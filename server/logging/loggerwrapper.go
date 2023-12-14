package logging

import (
	"fmt"

	"github.com/SashwatAnagolum/picodb/utils"
)

type LoggerWrapper struct {
	Wrappee   LoggerIF
	NewLogger LoggerIF
}

func NewLoggerWrapper(wrappee LoggerIF, newLogger LoggerIF) *LoggerWrapper {
	return &LoggerWrapper{
		Wrappee: wrappee, NewLogger: newLogger}
}

func (wrapper *LoggerWrapper) Log(request *utils.PicoDBRequest) {
	wrapper.NewLogger.Log(request)
	wrapper.Wrappee.Log(request)
}

func (wrapper *LoggerWrapper) GetLogs() (string, error) {
	newLogs, err := wrapper.NewLogger.GetLogs()
	oldLogs, err2 := wrapper.Wrappee.GetLogs()

	if err != nil {
		return "", err
	}

	if err2 != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s", newLogs, oldLogs), nil
}

func (wrapper *LoggerWrapper) ResetLogs() {
	wrapper.NewLogger.ResetLogs()
	wrapper.Wrappee.ResetLogs()
}
