package ui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/sirupsen/logrus"
)

type FilamentPanel struct {
	CommonPanel

	amount  *StepButton
	speed   *StepButton
	tool    *StepButton
	destroy chan bool
}

func NewFilamentPanel(ui *UI) *FilamentPanel {
	m := &FilamentPanel{CommonPanel: NewCommonPanel(ui),
		destroy: make(chan bool, 1),
	}

	m.initialize()
	return m
}

func (m *FilamentPanel) initialize() {
	m.grid.Attach(m.createExtrudeButton("Extrude", "extrude.svg", 1), 1, 0, 1, 1)
	m.grid.Attach(m.createExtrudeButton("Retract", "retract.svg", -1), 4, 0, 1, 1)

	m.tool = MustStepButton("extruct.svg")
	m.grid.Attach(m.tool, 1, 1, 1, 1)

	m.amount = MustStepButton("move-step.svg", Step{"5mm", 5}, Step{"10mm", 10}, Step{"1mm", 1})
	m.grid.Attach(m.amount, 2, 1, 1, 1)

	m.speed = MustStepButton("speed-step.svg", Step{"Normal", 5}, Step{"High", 10}, Step{"Slow", 1})
	m.grid.Attach(m.speed, 3, 1, 1, 1)

	m.grid.Attach(MustButtonImage("Return", "back.svg", m.UI.ShowDefaultPanel), 4, 1, 1, 1)

	go m.update()
}

func (m *FilamentPanel) update() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.updateInfo()
		case <-m.destroy:
			fmt.Println("done")
			return
		}
	}
}

func (m *FilamentPanel) updateInfo() {
	r, err := (&octoprint.ToolStateRequest{}).Do(m.UI.Printer)
	if err != nil {
		logrus.Error("FilamentPanel: %s", err)
	}

	box := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	box.SetVAlign(gtk.ALIGN_CENTER)

	for e, s := range r.Current {
		box.Add(MustLabel("%s: %f", e, s.Actual))
	}

	m.grid.Attach(box, 2, 0, 2, 1)
	m.grid.ShowAll()
}

func (m *FilamentPanel) createExtrudeButton(label, image string, dir int) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		cmd := &octoprint.ToolExtrudeRequest{}
		cmd.Amount = m.amount.Value() * dir

		if err := cmd.Do(m.UI.Printer); err != nil {
			panic(err)
		}
	})
}

func (m *FilamentPanel) Destroy() {
	m.destroy <- true
	m.grid.Destroy()
}
