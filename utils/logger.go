package utils

import (
	"io"
	standardLog "log"
	"os"
	"path"
	"runtime"
	"strings"
	// "time"

	"github.com/sirupsen/logrus"
)

type ContextHook struct{}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(strings.ToLower(name), "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = path.Base(file)
			entry.Data["func"] = path.Base(name)
			entry.Data["line"] = line
			break
		}
	}
	return nil
}


var log *logrus.Logger
var Logger *logrus.Entry

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
	standardLog.Printf("logger.SetLogLevel() - the log level is now set to: %q", level)
}

func LowerCaseLogLevel() string {
	logLevel := os.Getenv(EnvLogLevel)
	return strings.ToLower(logLevel)
}


func init() {
	log = logrus.New()
	log.AddHook(ContextHook{})

	// Start off with the logging level set to debug until we get a chance to read the configuration settings.
	SetLogLevel(logrus.DebugLevel)

	var logFilePath = os.Getenv("OCTOSCREEN_LOG_FILE_PATH")

	if logFilePath == "" {
		log.Infof("logger.init() - logFilePath is was not defined.  Now using just the standard console output.")
		log.Out = os.Stdout
	} else {
		log.Infof("logger.init() - logFilePath is: %q", logFilePath)
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Infof("logger.init() - OpenFile() succeeded and now setting log.Out to %s", logFilePath)
			log.Out = file

			log.Out = io.MultiWriter(os.Stdout, file)
			logrus.SetOutput(log.Out)
		} else {
			log.Errorf("logger.init() - OpenFile() FAILED!  err is: %q", err)
			log.Error("Failed to log to file, using default stderr")
		}
	}

	Logger = log.WithFields(logrus.Fields{})
}

func LogError(currentFunctionName, functionNameCalled string, err error) {
	if err != nil {
		Logger.Errorf("%s - %s returned an error: %q", currentFunctionName, functionNameCalled, err)
	} else {
		Logger.Errorf("%s - %s returned an error", currentFunctionName, functionNameCalled)
	}
}

func LogFatalError(currentFunctionName, functionNameCalled string, err error) {
	if err != nil {
		Logger.Fatalf("%s - %s returned an error: %q", currentFunctionName, functionNameCalled, err)
	} else {
		Logger.Fatalf("%s - %s returned an error", currentFunctionName, functionNameCalled)
	}
}
