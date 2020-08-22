package ui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
	// "github.com/Z-Bolt/OctoScreen/utils"
)


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








type Notifications struct {
	*gtk.Box
}

func NewNotifications() *Notifications {
	b := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	b.SetVAlign(gtk.ALIGN_START)
	b.SetHAlign(gtk.ALIGN_CENTER)
	b.SetHExpand(true)

	n := &Notifications{Box: b}
	h := NewNotificationsHook(n)
	logrus.AddHook(h)

	return n
}

func (n *Notifications) Show(style, msg string, d time.Duration) {
	defer n.Box.ShowAll()

	l := n.newLabel(style, msg)
	n.Box.Add(l)

	go func() {
		time.Sleep(d)
		glib.IdleAdd(l.Destroy)
	}()
}

func (n *Notifications) newLabel(style, msg string) *gtk.EventBox {
	l := MustLabel("")
	l.SetMarkup(fmt.Sprintf("<b>%s</b>", msg))
	l.SetLineWrap(true)

	ctx, _ := l.GetStyleContext()
	ctx.AddClass("notification")
	ctx.AddClass(style)

	e, _ := gtk.EventBoxNew()
	e.Add(l)
	e.SetEvents(int(gdk.BUTTON_PRESS_MASK))
	e.Connect("button-press-event", e.Destroy)

	return e
}
