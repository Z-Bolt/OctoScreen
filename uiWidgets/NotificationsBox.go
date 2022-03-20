package uiWidgets

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type NotificationsBox struct {
	*gtk.Box
}

func NewNotificationsBox() *NotificationsBox {
	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	base.SetVAlign(gtk.ALIGN_START)
	base.SetHAlign(gtk.ALIGN_CENTER)
	base.SetHExpand(true)

	instance := &NotificationsBox {
		Box: base,
	}
	notificationsHook := NewNotificationsHook(instance)
	logrus.AddHook(notificationsHook)

	return instance
}

func (this *NotificationsBox) Show(style, msg string, duration time.Duration) {
	defer this.Box.ShowAll()

	eventBox := this.newEventBox(style, msg)
	this.Box.Add(eventBox)

	go func() {
		time.Sleep(duration)
		glib.IdleAdd(eventBox.Destroy)
	}()
}

func (this *NotificationsBox) newEventBox(style, msg string) *gtk.EventBox {
	label := utils.MustLabel("")
	label.SetMarkup(fmt.Sprintf("<b>%s</b>", msg))
	label.SetLineWrap(true)

	ctx, _ := label.GetStyleContext()
	ctx.AddClass("notification")
	ctx.AddClass(style)

	eventBox, _ := gtk.EventBoxNew()
	eventBox.Add(label)
	eventBox.SetEvents(int(gdk.BUTTON_PRESS_MASK))
	eventBox.Connect("button-press-event", eventBox.Destroy)

	return eventBox
}
