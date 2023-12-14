package logging

import (
	"fmt"
	"os"
)

type DiskLogWriter struct {
	Filepath  string
	numWrites int
}

func NewDiskLogWriter() LogWriterIF {
	return DiskLogWriter{Filepath: "./saved_logs", numWrites: 0}
}

func (writer DiskLogWriter) WriteLogs(logs string) {
	file, err := os.OpenFile(
		fmt.Sprintf("%s_%d.txt", writer.Filepath, writer.numWrites),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(logs)

	if err != nil {
		fmt.Println(err)
	}

	writer.numWrites += 1
}
