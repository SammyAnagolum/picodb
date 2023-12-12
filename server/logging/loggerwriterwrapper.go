package logging

import (
	"fmt"

	"github.com/SashwatAnagolum/picodb/utils"
)

type LoggerWriterWrapper struct {
	Logger             LoggerIF
	Writer             LogWriterIF
	LogWriteFrequency  int
	numUnwrittenWrites int
}

func NewLoggerWriterWrapper(logger LoggerIF,
	writer LogWriterIF, writeFreq int) *LoggerWriterWrapper {
	return &LoggerWriterWrapper{
		Logger:             logger,
		Writer:             writer,
		LogWriteFrequency:  writeFreq,
		numUnwrittenWrites: 0}
}

func (wrapper *LoggerWriterWrapper) Log(request *utils.PicoDBRequest) {
	wrapper.numUnwrittenWrites += 1

	wrapper.Logger.Log(request)

	if wrapper.numUnwrittenWrites == wrapper.LogWriteFrequency {
		logs, err := wrapper.Logger.GetLogs()

		if err != nil {
			fmt.Println(err)
			return
		}

		wrapper.Writer.WriteLogs(logs)
		wrapper.numUnwrittenWrites = 0
	}
}

func (wrapper *LoggerWriterWrapper) GetLogs() (string, error) {
	return wrapper.Logger.GetLogs()
}

func (wrapper *LoggerWriterWrapper) ResetLogs() {
	wrapper.Logger.ResetLogs()
}

func (wrapper *LoggerWriterWrapper) GetName() string {
	return wrapper.Logger.GetName()
}
