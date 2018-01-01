package main

import (
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/ui"
)

const (
	EnvStylePath = "OCTOPRINT_TFT_STYLE_PATH"
	EnvBaseURL   = "OCTOPRINT_HOST"
	EnvAPIKey    = "OCTOPRINT_APIKEY"

	DefaultBaseURL = "http://127.0.0.1"
)

var (
	BaseURL string
	APIKey  string
)

func init() {
	ui.StylePath = os.Getenv(EnvStylePath)
	BaseURL = os.Getenv(EnvBaseURL)
	APIKey = os.Getenv(EnvAPIKey)
	if BaseURL == "" {
		BaseURL = DefaultBaseURL
	}
}

func main() {
	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	ui.New(BaseURL, APIKey)
	gtk.Main()
}
