package logging

import (
	"fmt"
	"time"

	"github.com/zombinome/go-microserv/microserv"
)

type ConsoleLogger struct {
	name        string
	requestId   string
	isScoped    bool
	minLevel    microserv.LogLevel
	opStartTime *timePoint
}

func (logger *ConsoleLogger) logMessage(logLevel microserv.LogLevel, message string) error {
	if logger.minLevel > logLevel {
		return nil
	}

	var code = logCodes[logLevel]
	_, err := fmt.Printf("[%s][%s][%s]: %s\n", code, logger.requestId, logger.name, message)
	return err
}

func (logger *ConsoleLogger) logMessageF(logLevel microserv.LogLevel, message string, args []interface{}) error {
	if logger.minLevel > logLevel {
		return nil
	}

	var code = logCodes[logLevel]
	var msgArgs = append([]interface{}{code, logger.requestId, logger.name}, args...)
	_, err := fmt.Printf("[%s][%s][%s]: "+message+"\n", msgArgs...)
	return err
}

func (logger ConsoleLogger) Info(message string) {
	logger.logMessage(microserv.LogLevelInfo, message)
}

func (logger ConsoleLogger) InfoF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelInfo, message, args)
}

func (logger ConsoleLogger) Warn(message string) {
	logger.logMessage(microserv.LogLevelWarn, message)
}

func (logger ConsoleLogger) WarnF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelWarn, message, args)
}

func (logger ConsoleLogger) Error(message string) {
	logger.logMessage(microserv.LogLevelError, message)
}

func (logger ConsoleLogger) ErrorF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelError, message, args)
}

func (logger ConsoleLogger) BeginMeasure(message string) {
	logger.Info(message)
	logger.opStartTime.point = time.Now()
}
func (logger ConsoleLogger) EndMeasure(message string) {
	var elapsed = time.Since(logger.opStartTime.point)
	logger.InfoF("%s, took %s", message, elapsed)
}

func (logger ConsoleLogger) CreateScope(name string, requestId string) microserv.LoggerScope {
	return &ConsoleLogger{
		name:        logger.name + "->" + name,
		requestId:   requestId,
		isScoped:    true,
		minLevel:    logger.minLevel,
		opStartTime: &timePoint{},
	}
}

func NewConsoleLogger(name string, minLevel microserv.LogLevel) microserv.Logger {
	var logger microserv.Logger = ConsoleLogger{
		name:        name,
		requestId:   "",
		isScoped:    false,
		minLevel:    minLevel,
		opStartTime: &timePoint{},
	}

	return logger
}
