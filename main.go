package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/ui"
)

const (
	EnvImagesFolder = "OCTOPRINT_TFT_IMAGES"
	EnvBaseURL      = "OCTOPRINT_HOST"
	EnvAPIKey       = "OCTOPRINT_APIKEY"

	DefaultBaseURL = "http://127.0.0.1"
)

var (
	BaseURL string
	APIKey  string
)

func init() {
	ui.ImagesFolder = os.Getenv(EnvImagesFolder)
	BaseURL = os.Getenv(EnvBaseURL)
	APIKey = os.Getenv(EnvAPIKey)

	if BaseURL == "" {
		BaseURL = DefaultBaseURL
	}
}

func main() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	ui := ui.New(BaseURL, APIKey)
	win.Add(ui)
	win.SetTitle("foo")
	win.SetDefaultSize(480, 320)
	win.ShowAll()
	gtk.Main()
}
