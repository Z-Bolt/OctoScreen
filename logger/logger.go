package logger

import (
	"fmt"
	"io"
	standardLog "log"
	"os"
	// "path"
	// "runtime"
	"strings"
	"sync"
	// "time"

	"github.com/sirupsen/logrus"
)


var _mutex *sync.RWMutex = nil
var _indentLevel int = -1
var _indentation string = ""
const INDENTATION_TOKEN = "    "
const INDENTATION_TOKEN_LENGTH = 4

var _logrusLogger *logrus.Logger
var _logrusEntry *logrus.Entry
var _logLevel logrus.Level
var _strLogLevel string


func init() {
	_mutex = &sync.RWMutex {}

	_indentLevel = 0
	_indentation = ""

	_logrusLogger = logrus.New()
	_logrusLogger.AddHook(ContextHook{})

	//
	// TODO: ...(maybe?) it would be nice it this could be made generic,
	// but this is getting set in init().
	var logFilePath = os.Getenv("OCTOSCREEN_LOG_FILE_PATH")
	//

	if logFilePath == "" {
		standardLog.Print("logger.init() - logFilePath is was not defined.  Now using just the standard console output.")
		_logrusLogger.Out = os.Stdout
	} else {
		standardLog.Printf("logger.init() - logFilePath is: %s", logFilePath)
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			standardLog.Printf("logger.init() - OpenFile() succeeded and now setting log.Out to %s", logFilePath)
			_logrusLogger.Out = file

			_logrusLogger.Out = io.MultiWriter(os.Stdout, file)
			logrus.SetOutput(_logrusLogger.Out)
		} else {
			standardLog.Printf("logger.init() - OpenFile() FAILED!  err is: %s", err.Error())
			standardLog.Print("Failed to open the log file, defaulting to use the standard console output.")
			_logrusLogger.Out = os.Stdout
		}
	}

	_logrusEntry = _logrusLogger.WithFields(logrus.Fields{})

	// Start off with the logging level set to debug until we get a chance to read the configuration settings.
	SetLogLevel(logrus.DebugLevel)
}

func readIndentation() string {
	_mutex.RLock()
	var indentation = _indentation
	_mutex.RUnlock()

	return indentation
}

func writeIndentation(indentation string) {
	_mutex.Lock()
	_indentation = indentation
	_mutex.Unlock()
}


func SetLogLevel(newLevel logrus.Level) {
	_logLevel = newLevel
	_strLogLevel = strings.ToLower(_logLevel.String())

	_logrusLogger.SetLevel(_logLevel)
	standardLog.Printf("logger.SetLogLevel() - the log level is now set to: %s", _strLogLevel)
}

func LogLevel() string {
	// Returns a lower case string.
	return _strLogLevel
}


func TraceEnter(functionName string) {
	message := fmt.Sprintf("%sentering %s", readIndentation(), functionName)
	_logrusEntry.Debug(message)
	_indentLevel++
	indentation := readIndentation() + INDENTATION_TOKEN
	writeIndentation(indentation)
}

func TraceLeave(functionName string) {
	_indentLevel--

	indentation := readIndentation()
	indentation = indentation[:(_indentLevel * INDENTATION_TOKEN_LENGTH)]
	writeIndentation(indentation)

	message := fmt.Sprintf("%sleaving %s", readIndentation(), functionName)
	_logrusEntry.Debug(message)
}


func LogError(currentFunctionName, functionCalledName string, err error) {
	if err != nil {
		_logrusEntry.Errorf("%s%s - %s returned an error: %q", readIndentation(), currentFunctionName, functionCalledName, err)
	} else {
		_logrusEntry.Errorf("%s%s - %s returned an error", readIndentation(), currentFunctionName, functionCalledName)
	}
}

func LogFatalError(currentFunctionName, functionCalledName string, err error) {
	if err != nil {
		_logrusEntry.Fatalf("%s%s - %s returned an error: %q", readIndentation(), currentFunctionName, functionCalledName, err)
	} else {
		_logrusEntry.Fatalf("%s%s - %s returned an error", readIndentation(), currentFunctionName, functionCalledName)
	}
}


func Debug(args ...interface{}) {
	_logrusEntry.Debug(readIndentation() + fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	_logrusEntry.Debugf(readIndentation() + format, args...)
}


func Info(args ...interface{}) {
	_logrusEntry.Info(readIndentation() + fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	_logrusEntry.Infof(readIndentation() + format, args...)
}


func Warn(args ...interface{}) {
	_logrusEntry.Warn(readIndentation() + fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	_logrusEntry.Warnf(readIndentation() + format, args...)
}


func Error(args ...interface{}) {
	_logrusEntry.Error(readIndentation() + fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	_logrusEntry.Errorf(readIndentation() + format, args...)
}


func Fatal(args ...interface{}) {
	_logrusEntry.Fatal(readIndentation() + fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	_logrusEntry.Fatalf(readIndentation() + format, args...)
}

func Panic(args ...interface{}) {
	_logrusEntry.Panic(readIndentation() + fmt.Sprint(args...))
}

func Panicf(format string, args ...interface{}) {
	_logrusEntry.Panicf(readIndentation() + format, args...)
}
