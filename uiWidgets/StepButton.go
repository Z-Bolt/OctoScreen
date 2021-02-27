package uiWidgets

import (
	"fmt"
	// "strings"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type Step struct {
	Label			string
	ImageFileName	string
	Image			gtk.IWidget
	Value			interface{}
}

type StepButton struct {
	*gtk.Button
	sync.RWMutex

	Steps			[]Step
	CurrentStep		int
}

func CreateStepButton(colorVariation int, steps ...Step) (*StepButton, error) {
	if len(steps) < 1 {
		logger.Error("PANIC!!! - CreateStepButton() - len(steps) < 1")
		panic("StepButton.CreateStepButton() - steps is empty")
	}

	base := utils.MustButtonImage(steps[0].Label, steps[0].ImageFileName, nil)
	ctx, _ := base.GetStyleContext()
	colorClass := fmt.Sprintf("color-dash-%d", colorVariation)
	ctx.AddClass(colorClass)

	instance := &StepButton{
		Button:			base,
		Steps:			steps,
		CurrentStep:	-1,
	}

	if len(steps) > 0 {
		for i := 0; i < len(steps); i++ {
			instance.Steps[i].Image = utils.MustImageFromFile(instance.Steps[i].ImageFileName)
		}

		instance.CurrentStep = 0
	}

	_, err := instance.Button.Connect("clicked", instance.handleClick)

	return instance, err
}

func (this *StepButton) Value() interface{} {
	this.RLock()
	defer this.RUnlock()

	return this.Steps[this.CurrentStep].Value
}

func (this *StepButton) AddStep(step Step) {
	this.Lock()
	defer this.Unlock()

	if this.Steps == nil || len(this.Steps) == 0 {
		this.Steps = make([]Step, 0)
	}

	this.Steps = append(this.Steps, step)
	index := len(this.Steps) - 1
	this.Steps[index].Image = utils.MustImageFromFile(this.Steps[index].ImageFileName)
}

func (this *StepButton) handleClick() {
	this.RLock()
	defer this.RUnlock()

	if len(this.Steps) < 1 {
		return
	}

	this.CurrentStep++
	if this.CurrentStep >= len(this.Steps) {
		this.CurrentStep = 0
	}

	this.Button.SetLabel(this.Steps[this.CurrentStep].Label)
	this.Button.SetImage(this.Steps[this.CurrentStep].Image)
}
