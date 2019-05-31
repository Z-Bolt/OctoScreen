package ui

import (
	"os"
	"path"
	"runtime"
	"strings"
	"time"

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
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = path.Base(file)
			entry.Data["func"] = path.Base(name)
			entry.Data["line"] = line
			break
		}
	}
	return nil
}

type NotificationsHook struct {
	n *Notifications
}

func NewNotificationsHook(n *Notifications) *NotificationsHook {
	return &NotificationsHook{n: n}
}

func (h NotificationsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func (h NotificationsHook) Fire(entry *logrus.Entry) error {
	d := 10 * time.Second
	if entry.Level == logrus.WarnLevel {
		d = time.Second
	}

	h.n.Show(entry.Level.String(), entry.Message, d)
	return nil
}

var Logger *logrus.Entry

func init() {
	var LogFile = os.Getenv("OCTOPRINT_TFT_LOG_FILE")

	var log = logrus.New()
	log.AddHook(ContextHook{})
	log.SetLevel(logrus.DebugLevel)

	if LogFile == "" {
		log.Out = os.Stdout
	} else {
		file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}

	Logger = log.WithFields(logrus.Fields{})
}
