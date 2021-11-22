package microserv

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelOff
)

type Logger interface {
	Info(message string)

	InfoF(message string, args ...interface{})

	Warn(message string)

	WarnF(message string, args ...interface{})

	Error(message string)

	ErrorF(message string, args ...interface{})

	CreateScope(name string, requestId string) LoggerScope
}

type LoggerScope interface {
	Logger

	BeginMeasure(message string)
	EndMeasure(message string)
}
