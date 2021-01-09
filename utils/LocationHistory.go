package utils

import (
	// "fmt"
	// "sort"

	"github.com/mcuadros/go-octoprint"
	// "github.com/mcuadros/go-octoprint/apis"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type LocationHistory struct {
	Locations []octoprint.Location
}

func (this *LocationHistory) Length() int {
	return len(this.Locations)
}

func (this *LocationHistory) CurrentLocation() octoprint.Location {
	length := this.Length()
	if length < 1 {
		panic("LocationHistory.current() - locations is empty")
	}

	return this.Locations[length - 1]
}

func (this *LocationHistory) GoForward(folder string) {
	newLocation := string(this.CurrentLocation()) + "/" + folder
	this.Locations = append(this.Locations, octoprint.Location(newLocation))
}

func (this *LocationHistory) GoBack() {
	this.Locations = this.Locations[0 : len(this.Locations) - 1]
}

func (this *LocationHistory) IsRoot() bool {
	if len(this.Locations) > 1 {
		return false
	} else {
		return true
	}
}
