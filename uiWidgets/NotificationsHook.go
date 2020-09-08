package uiWidgets

import (
	// "fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type NotificationsHook struct {
	notificationsBox *NotificationsBox
}

func NewNotificationsHook(notificationsBox *NotificationsBox) *NotificationsHook {
	return &NotificationsHook {
		notificationsBox: notificationsBox,
	}
}

func (this NotificationsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func (this NotificationsHook) Fire(entry *logrus.Entry) error {
	d := 10 * time.Second
	if entry.Level == logrus.WarnLevel {
		d = time.Second
	}

	this.notificationsBox.Show(entry.Level.String(), entry.Message, d)

	return nil
}
