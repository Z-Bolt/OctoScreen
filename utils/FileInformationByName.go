package utils

import (
	// "fmt"
	// "sort"
	// "strings"

	"github.com/mcuadros/go-octoprint"
	// "github.com/mcuadros/go-octoprint/apis"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type FileInformationsByName []*octoprint.FileInformation

func (this FileInformationsByName) Len() int {
	 return len(this)
}

func (this FileInformationsByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this FileInformationsByName) Less(i, j int) bool {
	return this[j].Name > this[i].Name
}
