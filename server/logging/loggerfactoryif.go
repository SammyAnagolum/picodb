package logging

type LoggerFactoryIF interface {
	NewLogger(loggerName string) LoggerIF
}
