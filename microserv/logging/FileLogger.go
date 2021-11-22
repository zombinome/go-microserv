package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zombinome/go-microserv/microserv"
)

const (
	logLevelCodeError string = "ERR"
	logLevelCodeWarn  string = "WRN"
	logLevelCodeInfo  string = "INF"
	logLevelCodeTrace string = "TRC"
)

var logCodes = [5]string{logLevelCodeTrace, logLevelCodeInfo, logLevelCodeWarn, logLevelCodeError, ""}

type FileLogger struct {
	name        string
	requestId   string
	isScoped    bool
	logger      *log.Logger
	logFile     *os.File
	minLevel    microserv.LogLevel
	opStartTime *timePoint
}

func (logger *FileLogger) logMessage(logLevel microserv.LogLevel, message string) error {
	if logger.minLevel > logLevel {
		return nil
	}

	var code = logCodes[logLevel]
	var logMessage = fmt.Sprintf("[%s][%s][%s]: %s", code, logger.requestId, logger.name, message)
	return logger.logger.Output(2, logMessage)
}

func (logger *FileLogger) logMessageF(logLevel microserv.LogLevel, message string, args []interface{}) error {
	if logger.minLevel > logLevel {
		return nil
	}

	var code = logCodes[logLevel]
	var msgArgs = append([]interface{}{code, logger.requestId, logger.name}, args...)
	var logMessage = fmt.Sprintf("[%s][%s][%s]: "+message, msgArgs...)
	return logger.logger.Output(2, logMessage)
}

func (logger FileLogger) Info(message string) {
	logger.logMessage(microserv.LogLevelInfo, message)
}

func (logger FileLogger) InfoF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelInfo, message, args)
}

func (logger FileLogger) Warn(message string) {
	logger.logMessage(microserv.LogLevelWarn, message)
}

func (logger FileLogger) WarnF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelWarn, message, args)
}

func (logger FileLogger) Error(message string) {
	logger.logMessage(microserv.LogLevelError, message)
}

func (logger FileLogger) ErrorF(message string, args ...interface{}) {
	logger.logMessageF(microserv.LogLevelError, message, args)
}

func (logger FileLogger) BeginMeasure(message string) {
	logger.Info(message)
	logger.opStartTime.point = time.Now()
}
func (logger FileLogger) EndMeasure(message string) {
	var elapsed = time.Since(logger.opStartTime.point)
	logger.InfoF("%s, took %s", message, elapsed)
}

func (logger *FileLogger) Close() error {
	if logger.logFile != nil {
		return logger.logFile.Close()
	}

	return nil
}

func (logger FileLogger) CreateScope(name string, requestId string) microserv.LoggerScope {
	return &FileLogger{
		name:        logger.name + "->" + name,
		requestId:   requestId,
		isScoped:    true,
		logger:      logger.logger,
		logFile:     nil,
		minLevel:    logger.minLevel,
		opStartTime: &timePoint{},
	}
}

func NewFileLogger(name string, pathToLogFile string, minLevel microserv.LogLevel) *microserv.Logger {
	microserv.EnsurePathExists(pathToLogFile)
	fileWriter, err := os.OpenFile(pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("Invalid log path")
	}

	var logger = log.New(fileWriter, "", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	var result microserv.Logger = FileLogger{
		name:        name,
		requestId:   "",
		isScoped:    false,
		logger:      logger,
		logFile:     fileWriter,
		minLevel:    minLevel,
		opStartTime: &timePoint{},
	}

	return &result
}
