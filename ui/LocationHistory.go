package ui

import (
	// "fmt"
	// "sort"

	"github.com/mcuadros/go-octoprint"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type byDate []*octoprint.FileInformation

func (s byDate) Len() int           { return len(s) }
func (s byDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byDate) Less(i, j int) bool { return s[j].Date.Time.Before(s[i].Date.Time) }

type locationHistory struct {
	locations []octoprint.Location
}

func (l *locationHistory) current() octoprint.Location {
	return l.locations[len(l.locations) - 1]
}

func (l *locationHistory) goForward(folder string) {
	newLocation := string(l.current()) + "/" + folder
	l.locations = append(l.locations, octoprint.Location(newLocation))
}

func (l *locationHistory) goBack() {
	l.locations = l.locations[0 : len(l.locations) - 1]
}

func (l *locationHistory) isRoot() bool {
	if len(l.locations) > 1 {
		return false
	} else {
		return true
	}
}
