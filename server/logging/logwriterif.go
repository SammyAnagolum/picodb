package logging

type LogWriterIF interface {
	WriteLogs(logs string)
}
