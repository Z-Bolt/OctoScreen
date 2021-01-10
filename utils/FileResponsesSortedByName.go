package utils

import (
	// "fmt"
	// "sort"
	// "strings"

	"github.com/mcuadros/go-octoprint"
	// "github.com/mcuadros/go-octoprint/apis"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type FileResponsesSortedByName []*octoprint.FileResponse

func (this FileResponsesSortedByName) Len() int {
	 return len(this)
}

func (this FileResponsesSortedByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this FileResponsesSortedByName) Less(i, j int) bool {
	return this[j].Name > this[i].Name
}
