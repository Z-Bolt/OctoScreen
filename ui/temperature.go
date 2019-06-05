package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var temperaturePanelInstance *temperaturePanel

type temperaturePanel struct {
	CommonPanel

	tool   *StepButton
	amount *StepButton

	box    *gtk.Box
	labels map[string]*LabelWithImage
}

func TemperaturePanel(ui *UI, parent Panel) Panel {
	if temperaturePanelInstance == nil {
		m := &temperaturePanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string]*LabelWithImage{},
		}

		m.b = NewBackgroundTask(time.Second, m.updateTemperatures)
		m.initialize()

		temperaturePanelInstance = m
	}

	return temperaturePanelInstance
}

func (m *temperaturePanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createChangeButton("Increase", "increase.svg", 1), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeButton("Decrease", "decrease.svg", -1), 4, 0, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetMarginStart(10)
	m.Grid().Attach(m.box, 2, 0, 2, 1)

	m.Grid().Attach(m.createToolButton(), 1, 1, 1, 1)
	m.amount = MustStepButton("move-step.svg", Step{"10°C", 10.}, Step{"5°C", 5.}, Step{"1°C", 1.})
	m.Grid().Attach(m.amount, 2, 1, 1, 1)

	m.Grid().Attach(MustButtonImage("More", "heat-up.svg", m.profilePanel), 3, 1, 1, 1)
}

func (m *temperaturePanel) createToolButton() *StepButton {
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

func (m *temperaturePanel) createChangeButton(label, image string, value float64) gtk.IWidget {

	return MustPressedButton(label, image, func() {
		target := value * m.amount.Value().(float64)
		if err := m.increaseTarget(m.tool.Value().(string), target); err != nil {
			Logger.Error(err)
			return
		}
	}, 100)
}

func (m *temperaturePanel) increaseTarget(tool string, value float64) error {
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

func (m *temperaturePanel) setTarget(tool string, target float64) error {
	if tool == "bed" {
		cmd := &octoprint.BedTargetRequest{Target: target}
		return cmd.Do(m.UI.Printer)
	}

	cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{tool: target}}
	return cmd.Do(m.UI.Printer)
}

func (m *temperaturePanel) getToolTarget(tool string) (float64, error) {
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

func (m *temperaturePanel) updateTemperatures() {
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

func (m *temperaturePanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}
}

func (m *temperaturePanel) addNewTool(tool string) {
	img := "extruder.svg"
	if tool == "bed" {
		img = "bed.svg"
	}

	m.labels[tool] = MustLabelWithImage(img, "")
	m.box.Add(m.labels[tool])
	m.tool.AddStep(Step{strings.Title(tool), tool})
	m.tool.Callback()

	Logger.Infof("Tool detected %s", tool)
}

func (m *temperaturePanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%s: %.1f°C / %.1f°C", strings.Title(tool), d.Actual, d.Target)
	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *temperaturePanel) profilePanel() {
	m.UI.Add(ProfilesPanel(m.UI, m))
}

var profilePanelInstance *profilesPanel

type profilesPanel struct {
	CommonPanel
}

func ProfilesPanel(ui *UI, parent Panel) Panel {
	if profilePanelInstance == nil {
		m := &profilesPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		profilePanelInstance = m
	}

	return profilePanelInstance
}

func (m *profilesPanel) initialize() {
	defer m.Initialize()
	m.loadProfiles()
}

func (m *profilesPanel) loadProfiles() {
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

func (m *profilesPanel) createProfileButton(img string, p *octoprint.TemperatureProfile) gtk.IWidget {
	return MustButtonImage(p.Name, img, func() {
		Logger.Warningf("Setting temperature profile %s.", p.Name)
		if err := m.setProfile(p); err != nil {
			Logger.Error(err)
		}
		m.UI.GoHistory()
	})
}

func (m *profilesPanel) setProfile(p *octoprint.TemperatureProfile) error {
	for tool := range temperaturePanelInstance.labels {
		temp := p.Extruder
		if tool == "bed" {
			temp = p.Bed
		}

		if err := temperaturePanelInstance.setTarget(tool, temp); err != nil {
			return err
		}
	}

	return nil
}
