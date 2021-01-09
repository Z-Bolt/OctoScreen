package utils

import (
	// "fmt"
	// "sort"

	"github.com/mcuadros/go-octoprint"
	// "github.com/mcuadros/go-octoprint/apis"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type FileInformationsByDate []*octoprint.FileInformation

func (this FileInformationsByDate) Len() int {
	 return len(this)
}

func (this FileInformationsByDate) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this FileInformationsByDate) Less(i, j int) bool {
	return this[j].Date.Time.Before(this[i].Date.Time)
}
