package utils

import (
	// "fmt"
	// "sort"

	"github.com/mcuadros/go-octoprint"
	// "github.com/mcuadros/go-octoprint/apis"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type FileResponsesSortedByDate []*octoprint.FileResponse

func (this FileResponsesSortedByDate) Len() int {
	 return len(this)
}

func (this FileResponsesSortedByDate) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this FileResponsesSortedByDate) Less(i, j int) bool {
	return this[j].Date.Time.Before(this[i].Date.Time)
}
