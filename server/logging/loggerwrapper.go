package logging

import (
	"fmt"

	"github.com/SashwatAnagolum/picodb/utils"
)

type LoggerWrapper struct {
	Wrappee   LoggerIF
	NewLogger LoggerIF
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

func (wrapper *LoggerWrapper) GetName() string {
	return wrapper.NewLogger.GetName()
}
