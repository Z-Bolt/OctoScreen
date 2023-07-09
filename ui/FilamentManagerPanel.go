package ui

import (
	"fmt"
	// "strings"

	// "github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


const MAX_EXTRUDER_COUNT = 5

type filamentManagerPanel struct {
	CommonPanel

	selectExtruderStepButton		*uiWidgets.SelectToolStepButton
	scrollableListBox				*uiWidgets.ScrollableListBox
	filamentManagerListBoxRows		[]*uiWidgets.FilamentManagerListBoxRow

	filamentManagerSelections		[]*dataModels.FilamentManagerSelection
	filamentManagerSpools			[]*dataModels.FilamentManagerSpool
	spoolSelectionIds				[MAX_EXTRUDER_COUNT]int // Support up to 5 extruders.
}

var filamentManagerPanelInstance *filamentManagerPanel

func GetFilamentManagerPanelInstance(
	ui							*UI,
	filamentManagerSelections	[]*dataModels.FilamentManagerSelection,
	filamentManagerSpools		[]*dataModels.FilamentManagerSpool,
) *filamentManagerPanel {
	if filamentManagerPanelInstance == nil {
		instance := &filamentManagerPanel {
			CommonPanel: CreateCommonPanel("FilamentManagerPanel", ui),
		}

		instance.initializeData(filamentManagerSelections, filamentManagerSpools)
		instance.initializeUi()
		filamentManagerPanelInstance = instance
	}

	return filamentManagerPanelInstance
}

func (this *filamentManagerPanel) initializeData(
	filamentManagerSelections	[]*dataModels.FilamentManagerSelection,
	filamentManagerSpools		[]*dataModels.FilamentManagerSpool,
) {
	logger.TraceEnter("FilamentManagerPanel.initializeData()")

	this.filamentManagerSelections = filamentManagerSelections
	this.filamentManagerSpools = filamentManagerSpools

	if this.filamentManagerSelections == nil || this.filamentManagerSpools == nil {
		// This should never be the case... not sure how to recover from this.
		logger.Error("FilamentManagerPanel.initializeData() - filamentManagerSelections and/or filamentManagerSpools is invalid")
		logger.TraceLeave("FilamentManagerPanel.initializeData()")
		return;
	}

	for i := 0; i < MAX_EXTRUDER_COUNT; i++ {
		this.spoolSelectionIds[i] = -1
	}

	for i := 0; i < len(this.filamentManagerSelections); i++ {
		selection := this.filamentManagerSelections[i]
		toolIndex := selection.Tool
		if toolIndex < 0 || toolIndex >= MAX_EXTRUDER_COUNT {
			logger.Errorf("FilamentManagerPanel.initializeData() - toolIndex [%d] is invalid", i)
			continue
		}

		this.spoolSelectionIds[toolIndex] = selection.Spool.Id
	}

	logger.TraceLeave("FilamentManagerPanel.initializeData()")
}

func (this *filamentManagerPanel) initializeUi() {
	logger.TraceEnter("FilamentManagerPanel.initializeUi()")

	defer this.Initialize()

	if this.filamentManagerSelections == nil || this.filamentManagerSpools == nil {
		// This should never be the case... not sure how to recover from this.
		logger.Error("FilamentManagerPanel.initializeUi() - filamentManagerSelections and/or filamentManagerSpools is invalid")
		logger.TraceLeave("FilamentManagerPanel.initializeUi()")
		return;
	}

	this.selectExtruderStepButton = uiWidgets.CreateSelectExtruderStepButton(
		this.UI.Client,
		false,
		3,
		this.handleExtruderStepClick,
	)

	extruderCount := utils.GetExtruderCount(this.UI.Client)
	if extruderCount >= 2 {
		this.Grid().Attach(this.selectExtruderStepButton, 0, 0, 1, 1)
	}

	this.createListBox(extruderCount)
	this.createListBoxRows(extruderCount)

	logger.TraceLeave("FilamentManagerPanel.initializeUi()")
}

func (this *filamentManagerPanel) createListBox(extruderCount int) {
	this.scrollableListBox = uiWidgets.CreateScrollableListBox()
	if extruderCount >= 2 {
		this.Grid().Attach(this.scrollableListBox, 1, 0, 3, 2)
	} else {
		this.Grid().Attach(this.scrollableListBox, 0, 0, 4, 2)
	}
}

func (this *filamentManagerPanel) createListBoxRows(extruderCount int) {
	for i := 0; i < len(this.filamentManagerSpools); i++ {
		spool := this.filamentManagerSpools[i]

		// When initializing, use the first extruder.
		spoolIsSelected := (this.spoolSelectionIds[0] == spool.Id)

		filamentManagerListBoxRow := uiWidgets.CreateFilamentManagerListBoxRow(
			extruderCount,
			spool,
			spoolIsSelected,
			i,
			0,
			this.handleRowClick,
		)
		this.filamentManagerListBoxRows = append(this.filamentManagerListBoxRows, filamentManagerListBoxRow)
		this.scrollableListBox.Add(filamentManagerListBoxRow)
	}
}

func (this *filamentManagerPanel) handleExtruderStepClick() {
	logger.TraceEnter("FilamentManagerPanel.handleExtruderStepClick()")

	this.logFilamentManagerSelections()
	this.logFilamentManagerSpools()

	currentStepIndex := this.selectExtruderStepButton.CurrentStepIndex
	// ...currentStepIndex is also the current tool index.
	logger.Debugf("FilamentManagerPanel.handleExtruderStepClick() - currentStepIndex: %d", currentStepIndex)

	currentSelectedSpoolId := this.spoolSelectionIds[currentStepIndex]
	spoolName := "none"
	if currentSelectedSpoolId >= 0 {
		filamentManagerSpool := this.findFilamentManagerSpoolFromSpoolId(currentSelectedSpoolId)
		spoolName = filamentManagerSpool.Name
	}
	logger.Debugf("FilamentManagerPanel.handleExtruderStepClick() - currentSelectedSpoolId: %d (%s)", currentSelectedSpoolId, spoolName)

	listItemIndex := this.findListItemRowIndexFromSpoolId(currentSelectedSpoolId)
	// It's OK if listItemIndex is -1...
	// ...when it's -1 that means there are no selected spools, and the list becomes unchecked.
	logger.Debugf("FilamentManagerPanel.handleExtruderStepClick() - listItemIndex: %d", listItemIndex)

	this.updateRadioButtons(listItemIndex)

	logger.TraceLeave("FilamentManagerPanel.handleExtruderStepClick()")
}

func (this *filamentManagerPanel) handleRowClick(button *gtk.Button, rowIndex int) {
	logger.TraceEnter("FilamentManagerPanel.handleRowClick()")

	filamentManagerListBoxRow := this.filamentManagerListBoxRows[rowIndex]
	svgImageRadioButton := filamentManagerListBoxRow.SvgImageRadioButton
	isSelected := svgImageRadioButton.IsSelected

	if isSelected != true {
		rowIndex := filamentManagerListBoxRow.RowIndex
		logger.Debugf("FilamentManagerPanel.handleRowClick() - rowIndex: %d", rowIndex)

		currentStepIndex := this.selectExtruderStepButton.CurrentStepIndex
		logger.Debugf("FilamentManagerPanel.handleRowClick() - currentStepIndex: %d", currentStepIndex)

		this.updateRadioButtons(rowIndex)
		this.updateData(currentStepIndex, rowIndex)
		this.sendUpdateToOctoPrint(currentStepIndex)

		this.logFilamentManagerSelections()
		this.logFilamentManagerSpools()
	}

	logger.TraceLeave("FilamentManagerPanel.handleRowClick()")
}

func (this *filamentManagerPanel) updateRadioButtons(rowIndex int) {
	logger.TraceEnter("FilamentManagerPanel.updateRadioButtons()")

	logger.Debugf("FilamentManagerPanel.updateRadioButtons() rowIndex is: %d", rowIndex)

	for i := 0; i < len(this.filamentManagerListBoxRows); i++ {
		filamentManagerListBoxRow := this.filamentManagerListBoxRows[i]
		svgImageRadioButton := filamentManagerListBoxRow.SvgImageRadioButton
		isSelected := svgImageRadioButton.IsSelected

		if i != rowIndex {
			if isSelected == true {
				logger.Debugf("FilamentManagerPanel.updateRadioButtons() unselecting %d", i)
				svgImageRadioButton.Unselect()
			}
		} else {
			if isSelected == false {
				logger.Debugf("FilamentManagerPanel.updateRadioButtons() now selecting %d", i)
				svgImageRadioButton.Select()
			}
		}
	}

	logger.TraceLeave("FilamentManagerPanel.updateRadioButtons()")
}

func (this *filamentManagerPanel) updateData(currentStepIndex int, rowIndex int) {
	logger.TraceEnter("FilamentManagerPanel.updateData()")

	// rowIndex is the newly selected spool ID.
	logger.Debugf("FilamentManagerPanel.updateData() - rowIndex: %d", rowIndex)

	// currentStepIndex is AKA the current tool index.
	toolId := currentStepIndex
	logger.Debugf("FilamentManagerPanel.updateData() - toolId: %d", toolId)

	this.logFilamentManagerSelections()
	this.logFilamentManagerSpools()

	spool := this.findFilamentManagerSpoolFromListItemRowIndex(rowIndex)
	if spool == nil {
		msg := fmt.Sprintf("Unable to change filament: could not find FilamentManagerSpool from rowIndex of %d", rowIndex)
		logger.Errorf(msg)
		utils.ErrorMessageDialogBox(this.UI.window, msg)
		logger.TraceLeave("FilamentManagerPanel.updateData()")
		return
	}

	newFilamentManagerSelection := this.findFilamentManagerSelectionFromToolId(toolId)
	if newFilamentManagerSelection == nil {
		// If nil is returned, that means this was a FilamentManagerSelection that hasn't been set yet.

		logger.Debugf("FilamentManagerPanel.updateData() - dumping logFilamentManagerSelections()")
		this.logFilamentManagerSelections()

		filamentManagerSelection := new(dataModels.FilamentManagerSelection)
		filamentManagerSelection.Tool = toolId
		filamentManagerSelection.Spool.Id = -1
		this.filamentManagerSelections = append(this.filamentManagerSelections, filamentManagerSelection)

		logger.Debugf("FilamentManagerPanel.updateData() - dumping logFilamentManagerSelections() again")
		this.logFilamentManagerSelections()

		newFilamentManagerSelection = filamentManagerSelection
	}

	newFilamentManagerSelection.Spool = *spool
	this.spoolSelectionIds[currentStepIndex] = spool.Id

	logger.TraceLeave("FilamentManagerPanel.updateData()")
}

func (this *filamentManagerPanel) sendUpdateToOctoPrint(currentStepIndex int) {
	logger.TraceEnter("FilamentManagerPanel.sendUpdateToOctoPrint()")

	this.logFilamentManagerSelections()
	this.logFilamentManagerSpools()

	toolId := currentStepIndex
	filamentManagerSelection := this.findFilamentManagerSelectionFromToolId(toolId)
	if filamentManagerSelection == nil {
		msg := fmt.Sprintf("Unable to change filament: could not find FilamentManagerSelection from toolId of %d", toolId)
		logger.Errorf(msg)
		utils.ErrorMessageDialogBox(this.UI.window, msg)
		logger.TraceLeave("FilamentManagerPanel.sendUpdateToOctoPrint()")
		return
	}

	spoolId := filamentManagerSelection.Spool.Id

	// TODO: verify we aren't currently printing
	request := &octoprintApis.FilamentManagerSetSelectionRequest {
		Tool:		toolId,
		Spool:		spoolId,
	}

	response, err := request.Do(this.UI.Client)
	if err != nil {
		logger.LogError("FilamentManagerPanel.sendUpdateToOctoPrint()", "Do(FilamentManagerSetSelectionRequest)", err)
		msg := fmt.Sprintf("Unable to change filament: %v", err)
		logger.Errorf(msg)
		utils.ErrorMessageDialogBox(this.UI.window, msg)
		logger.TraceLeave("FilamentManagerPanel.sendUpdateToOctoPrint()")
		return
	}

	logger.Debugf("FilamentManagerPanel.sendUpdateToOctoPrint() - response: %v", response)

	// NOTE: The UI in OctoScreen does not update the change in OctoPrint's browser window,
	// and one needs to manually refresh the web page in order to see the change.

	logger.TraceLeave("FilamentManagerPanel.sendUpdateToOctoPrint()")
}


// There are many "IDs" - here are some routines to help keep everything straight.
func (this *filamentManagerPanel) findListItemRowIndexFromSpoolId(spoolId int) int {
	maxListBoxRows := len(this.filamentManagerListBoxRows)
	if spoolId < 0 || spoolId >= maxListBoxRows {
		return -1
	}

	for i := 0; i < maxListBoxRows; i++ {
		filamentManagerListBoxRow := this.filamentManagerListBoxRows[i]
		filamentManagerSpool := filamentManagerListBoxRow.FilamentManagerSpool
		if filamentManagerSpool.Id == spoolId {
			return i
		}
	}

	return -1
}

func (this *filamentManagerPanel) findFilamentManagerSelectionFromToolId(toolId int) *dataModels.FilamentManagerSelection {
	extruderCount := utils.GetExtruderCount(this.UI.Client)
	if toolId < 0 || toolId >= extruderCount {
		return nil
	}

	maxSelections := len(this.filamentManagerSelections)
	for i := 0; i < maxSelections; i++ {
		filamentManagerSelection := this.filamentManagerSelections[i]

		// filamentManagerSelection.Tool is an int (the index) and not "tool0"
		if filamentManagerSelection.Tool == toolId {
			return filamentManagerSelection
		}
	}

	return nil
}

func (this *filamentManagerPanel) findFilamentManagerSpoolFromSpoolId(spoolId int) *dataModels.FilamentManagerSpool {
	maxSpools := len(this.filamentManagerSpools)
	if spoolId < 0 || spoolId >= maxSpools {
		return nil
	}

	for i := 0; i < maxSpools; i++ {
		filamentManagerSpool := this.filamentManagerSpools[i]

		// filamentManagerSelection.Tool is an int (the index) and not "tool0"
		if filamentManagerSpool.Id == spoolId {
			return filamentManagerSpool
		}
	}

	return nil
}

func (this *filamentManagerPanel) findFilamentManagerSpoolFromListItemRowIndex(rowIndex int) *dataModels.FilamentManagerSpool {
	maxSpools := len(this.filamentManagerSpools)
	if rowIndex < 0 || rowIndex >= maxSpools {
		return nil
	}

	return this.filamentManagerSpools[rowIndex]
}


// Some helper routines for debugging.
func (this *filamentManagerPanel) logFilamentManagerSelections() {
	logger.TraceEnter("FilamentManagerPanel.logFilamentManagerSelections()")

	for i := 0; i < len(this.filamentManagerSelections); i++ {
		filamentManagerSelection := this.filamentManagerSelections[i]
		logger.Infof("logFilamentManagerSelections() - filamentManagerSelection[%d]: %v", i, filamentManagerSelection)
	}

	logger.TraceLeave("FilamentManagerPanel.logFilamentManagerSelections()")
}

func (this *filamentManagerPanel) logFilamentManagerSpools() {
	logger.TraceEnter("FilamentManagerPanel.logFilamentManagerSpools()")

	for i := 0; i < len(this.filamentManagerSpools); i++ {
		filamentManagerSpool := this.filamentManagerSpools[i]
		logger.Infof("logFilamentManagerSpools() - filamentManagerSpool[%d]: %v", i, filamentManagerSpool)
	}

	logger.TraceLeave("FilamentManagerPanel.logFilamentManagerSpools()")
}

func (this *filamentManagerPanel) logSpoolSelectionIds() {
	for i := 0; i < MAX_EXTRUDER_COUNT; i++ {
		logger.Debugf("this.spoolSelectionIds[%d]: %d", i, this.spoolSelectionIds[i])
	}
}
