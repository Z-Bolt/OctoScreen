package ui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Panel interface {
	Grid() *gtk.Grid
	Destroy()
}

type CommonPanel struct {
	UI   *UI
	grid *gtk.Grid
	b    *BackgroundTask
}

func NewCommonPanel(ui *UI) CommonPanel {
	return CommonPanel{
		UI:   ui,
		grid: MustGrid(),
	}
}

func (p *CommonPanel) Show() {
	if p.b != nil {
		p.b.Start()
	}
}

func (p *CommonPanel) Grid() *gtk.Grid {
	return p.grid
}

func (p *CommonPanel) Destroy() {
	if p.b != nil {
		p.b.Close()
	}

	p.grid.Destroy()
}

type BackgroundTask struct {
	stop, resume, close chan bool

	d    time.Duration
	task func()
}

func NewBackgroundTask(d time.Duration, task func()) *BackgroundTask {
	return &BackgroundTask{
		task: task,
		d:    d,

		stop:   make(chan bool, 1),
		resume: make(chan bool, 1),
		close:  make(chan bool, 1),
	}
}

func (t *BackgroundTask) Start() {
	Logger.Debug("New background task started")
	go t.loop()
	t.resume <- true
}

func (t *BackgroundTask) Stop() {
	t.stop <- true
}

func (t *BackgroundTask) Resume() {
	t.resume <- true
}

func (t *BackgroundTask) Close() {
	t.close <- true
}

func (t *BackgroundTask) loop() {
	for <-t.resume {
		t.execute()

		ticker := time.NewTicker(t.d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.execute()
			case <-t.stop:
				fmt.Println("stop")
				break
			case <-t.close:
				Logger.Debug("Background task closed")
				return
			}
		}
	}
}

func (t *BackgroundTask) execute() {
	_, err := glib.IdleAdd(t.task)
	if err != nil {
		log.Fatal("IdleAdd() failed:", err)
	}
}

type StepButton struct {
	Current  int
	Steps    []Step
	Callback func()

	*gtk.Button
	sync.RWMutex
}

type Step struct {
	Label string
	Value interface{}
}

func MustStepButton(image string, s ...Step) *StepButton {
	var l string
	if len(s) != 0 {
		l = s[0].Label
	}

	b := &StepButton{
		Button: MustButtonImage(l, image, nil),
		Steps:  s,
	}

	b.Connect("clicked", b.clicked)
	return b
}

func (b *StepButton) Label() string {
	b.RLock()
	defer b.RUnlock()

	return b.Steps[b.Current].Label
}

func (b *StepButton) Value() interface{} {
	b.RLock()
	defer b.RUnlock()

	return b.Steps[b.Current].Value
}

func (b *StepButton) AddStep(s Step) {
	b.Lock()
	defer b.Unlock()

	if len(b.Steps) == 0 {
		b.SetLabel(s.Label)
	}

	b.Steps = append(b.Steps, s)
}

func (b *StepButton) clicked() {
	b.RLock()
	defer b.RUnlock()

	b.Current++
	if b.Current >= len(b.Steps) {
		b.Current = 0
	}

	b.SetLabel(b.Steps[b.Current].Label)

	if b.Callback != nil {
		b.Callback()
	}
}
