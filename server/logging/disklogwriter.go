package logging

import (
	"fmt"
	"os"
)

type DiskLogWriter struct {
	Filepath string
}

func NewDiskLogWriter(filepath string) DiskLogWriter {
	return DiskLogWriter{Filepath: filepath}
}

func (writer DiskLogWriter) WriteLogs(logs string) {
	file, err := os.OpenFile(writer.Filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(logs)

	if err != nil {
		fmt.Println(err)
	}
}
