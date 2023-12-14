package logging

type LogWriterFactoryIF interface {
	NewLogWriter(writeName string) (LogWriterIF, error)
}
