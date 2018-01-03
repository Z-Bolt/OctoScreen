package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/ui"
	"gopkg.in/yaml.v1"
)

const (
	EnvStylePath  = "OCTOPRINT_TFT_STYLE_PATH"
	EnvBaseURL    = "OCTOPRINT_HOST"
	EnvAPIKey     = "OCTOPRINT_APIKEY"
	EnvConfigFile = "OCTOPRINT_CONFIG_FILE"

	DefaultBaseURL = "http://localhost"
)

var (
	BaseURL    string
	APIKey     string
	ConfigFile string
)

func init() {
	ui.StylePath = os.Getenv(EnvStylePath)
	APIKey = os.Getenv(EnvAPIKey)

	BaseURL = os.Getenv(EnvBaseURL)
	if BaseURL == "" {
		BaseURL = DefaultBaseURL
	}

	ConfigFile = os.Getenv(EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	if APIKey == "" && ConfigFile != "" {
		APIKey = readAPIKey(ConfigFile)
		ui.Logger.Infof("Found API key at %q file", ConfigFile)
	}
}

func main() {
	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	ui.New(BaseURL, APIKey)
	gtk.Main()
}

var (
	configLocation = ".octoprint/config.yaml"
	homeOctoPi     = "/home/pi/"
)

func readAPIKey(config string) string {
	var cfg struct{ API struct{ Key string } }

	data, err := ioutil.ReadFile(config)
	if err != nil {
		ui.Logger.Fatal(err)
		return ""
	}

	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		ui.Logger.Fatalf("Error decoding YAML config file %q: %s", config, err)
		return ""
	}

	return cfg.API.Key
}

func findConfigFile() string {
	if file := doFindConfigFile(homeOctoPi); file != "" {
		return file
	}

	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return doFindConfigFile(usr.HomeDir)
}

func doFindConfigFile(home string) string {
	path := filepath.Join(home, configLocation)

	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}
