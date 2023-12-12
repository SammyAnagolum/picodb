package logging

import "fmt"

type TerminalLogWriter struct {
}

func NewTerminalLogWriter() TerminalLogWriter {
	return TerminalLogWriter{}
}

func (writer TerminalLogWriter) WriteLogs(logs string) {
	fmt.Println(logs)
}
