package ui

import (
	"sync"

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
