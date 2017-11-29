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
}

func NewCommonPanel(ui *UI) CommonPanel {
	return CommonPanel{grid: MustGrid(), UI: ui}
}

func (p *CommonPanel) Grid() *gtk.Grid {
	return p.grid
}

func (p *CommonPanel) Destroy() {
	p.grid.Destroy()
}

type BackgroundTask struct {
	Stop, Resume, Close chan bool

	d    time.Duration
	task func()
}

func NewBackgroundTask(d time.Duration, task func()) *BackgroundTask {
	return &BackgroundTask{
		task: task,
		d:    d,

		Stop:   make(chan bool, 1),
		Resume: make(chan bool, 1),
		Close:  make(chan bool, 1),
	}
}

func (t *BackgroundTask) Start() {
	go t.loop()
	t.Resume <- true
}

func (t *BackgroundTask) loop() {
	for <-t.Resume {
		t.execute()

		ticker := time.NewTicker(t.d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.execute()
			case <-t.Stop:
				fmt.Println("stop")
				break
			case <-t.Close:
				fmt.Println("close")
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
	Value int
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

func (b *StepButton) Value() int {
	b.RLock()
	defer b.RUnlock()

	return b.Steps[b.Current].Value
}

func (b *StepButton) AddStep(s Step) {
	b.Lock()
	defer b.Unlock()

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
