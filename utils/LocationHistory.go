package utils

import (
	// "fmt"
	// "sort"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
)

type LocationHistory struct {
	Locations []octoprintApis.Location
}

func (this *LocationHistory) Length() int {
	return len(this.Locations)
}

func (this *LocationHistory) CurrentLocation() octoprintApis.Location {
	length := this.Length()
	if length < 1 {
		panic("LocationHistory.current() - locations is empty")
	}

	return this.Locations[length - 1]
}

func (this *LocationHistory) GoForward(folder string) {
	newLocation := string(this.CurrentLocation()) + "/" + folder
	this.Locations = append(this.Locations, octoprintApis.Location(newLocation))
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
