package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var temperaturePanel *TemperaturePanel

type TemperaturePanel struct {
	CommonPanel

	tool   *StepButton
	amount *StepButton

	box    *gtk.Box
	labels map[string]*LabelWithImage
}

func NewTemperaturePanel(ui *UI, parent Panel) *TemperaturePanel {
	if temperaturePanel == nil {
		m := &TemperaturePanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string]*LabelWithImage{},
		}

		m.b = NewBackgroundTask(time.Second, m.updateTemperatures)
		m.initialize()

		temperaturePanel = m
	}

	return temperaturePanel
}

func (m *TemperaturePanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createChangeButton("Increase", "increase.svg", 1), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeButton("Decrease", "decrease.svg", -1), 4, 0, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetMarginStart(10)
	m.Grid().Attach(m.box, 2, 0, 2, 1)

	m.Grid().Attach(m.createToolButton(), 1, 1, 1, 1)
	m.amount = MustStepButton("move-step.svg", Step{"5°C", 5.}, Step{"10°C", 10.}, Step{"1°C", 1.})
	m.Grid().Attach(m.amount, 2, 1, 1, 1)

	m.Grid().Attach(MustButtonImage("Profiles", "heat-up.svg", m.profilePanel), 3, 1, 1, 1)
}

func (m *TemperaturePanel) createToolButton() *StepButton {
	m.tool = MustStepButton("")
	m.tool.Callback = func() {
		img := "extruder.svg"
		if m.tool.Value().(string) == "bed" {
			img = "bed.svg"
		}

		m.tool.SetImage(MustImageFromFile(img))
	}

	return m.tool
}

func (m *TemperaturePanel) createChangeButton(label, image string, value float64) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		target := value * m.amount.Value().(float64)
		if err := m.increaseTarget(m.tool.Value().(string), target); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *TemperaturePanel) increaseTarget(tool string, value float64) error {
	target, err := m.getToolTarget(tool)
	if err != nil {
		return err
	}

	target += value
	if target < 0 {
		target = 0
	}

	Logger.Infof("Setting target temperature for %s to %1.f°C.", tool, target)
	return m.setTarget(tool, target)
}

func (m *TemperaturePanel) setTarget(tool string, target float64) error {
	if tool == "bed" {
		cmd := &octoprint.BedTargetRequest{Target: target}
		return cmd.Do(m.UI.Printer)
	}

	cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{tool: target}}
	return cmd.Do(m.UI.Printer)
}

func (m *TemperaturePanel) getToolTarget(tool string) (float64, error) {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd", "state"}}).Do(m.UI.Printer)
	if err != nil {
		return -1, err
	}

	current, ok := s.Temperature.Current[tool]
	if !ok {
		return -1, fmt.Errorf("unable to find tool %q", tool)
	}

	return current.Target, nil
}

func (m *TemperaturePanel) updateTemperatures() {
	s, err := (&octoprint.StateRequest{
		History: true,
		Limit:   1,
		Exclude: []string{"sd", "state"},
	}).Do(m.UI.Printer)

	if err != nil {
		Logger.Error(err)
		return
	}

	m.loadTemperatureState(&s.Temperature)
}

func (m *TemperaturePanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}
}

func (m *TemperaturePanel) addNewTool(tool string) {
	img := "extruder.svg"
	if tool == "bed" {
		img = "bed.svg"
	}

	m.labels[tool] = MustLabelWithImage(img, "")
	m.box.Add(m.labels[tool])
	m.tool.AddStep(Step{strings.Title(tool), tool})
	m.tool.Callback()

	Logger.Infof("New tool detected %s", tool)
}

func (m *TemperaturePanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%s: %.1f°C / %.1f°C", strings.Title(tool), d.Actual, d.Target)
	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *TemperaturePanel) profilePanel() {
	m.UI.Add(NewProfilesPanel(m.UI, m))
}

var profilePanel *ProfilesPanel

type ProfilesPanel struct {
	CommonPanel
}

func NewProfilesPanel(ui *UI, parent Panel) Panel {
	if profilePanel == nil {
		m := &ProfilesPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		profilePanel = m
	}

	return profilePanel
}

func (m *ProfilesPanel) initialize() {
	defer m.Initialize()
	m.loadProfiles()
}

func (m *ProfilesPanel) loadProfiles() {
	s, err := (&octoprint.SettingsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	for _, profile := range s.Temperature.Profiles {
		m.AddButton(m.createProfileButton("filament.svg", profile))
	}

	m.AddButton(m.createProfileButton("cool-down.svg", &octoprint.TemperatureProfile{
		Name:     "Cool Down",
		Bed:      0,
		Extruder: 0,
	}))
}

func (m *ProfilesPanel) createProfileButton(img string, p *octoprint.TemperatureProfile) gtk.IWidget {
	return MustButtonImage(p.Name, img, func() {
		Logger.Warningf("Setting temperature profile %s.", p.Name)
		if err := m.setProfile(p); err != nil {
			Logger.Error(err)
		}
	})
}

func (m *ProfilesPanel) setProfile(p *octoprint.TemperatureProfile) error {
	for tool := range temperaturePanel.labels {
		temp := p.Extruder
		if tool == "bed" {
			temp = p.Bed
		}

		if err := temperaturePanel.setTarget(tool, temp); err != nil {
			return err
		}
	}

	return nil
}
