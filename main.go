package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/ui"
	"gopkg.in/yaml.v1"
)

const (
	EnvStylePath  = "OCTOPRINT_TFT_STYLE_PATH"
	EnvResolution = "OCTOPRINT_TFT_RESOLUTION"
	EnvBaseURL    = "OCTOPRINT_HOST"
	EnvAPIKey     = "OCTOPRINT_APIKEY"
	EnvConfigFile = "OCTOPRINT_CONFIG_FILE"
)

var (
	BaseURL    string
	APIKey     string
	ConfigFile string
	Resolution string
)

func init() {
	ui.StylePath = os.Getenv(EnvStylePath)
	Resolution = os.Getenv(EnvResolution)

	ConfigFile = os.Getenv(EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	cfg := readConfig(ConfigFile)

	BaseURL = os.Getenv(EnvBaseURL)
	if BaseURL == "" {
		BaseURL = fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
		ui.Logger.Infof("Using %q as server address", BaseURL)

	}

	APIKey = os.Getenv(EnvAPIKey)
	if APIKey == "" {
		APIKey = cfg.API.Key
		if cfg.API.Key != "" {
			ui.Logger.Infof("Found API key at %q file", ConfigFile)
		}
	}
}

func main() {
	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()

	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	width, height := getSize()
	_ = ui.New(BaseURL, APIKey, width, height)

	gtk.Main()
}

var (
	configLocation = ".octoprint/config.yaml"
	homeOctoPi     = "/home/pi/"
)

type config struct {
	// API Settings.
	API struct {
		// Key is the current API key needed for accessing the API.
		Key string
	}
	// Server settings.
	Server struct {
		// Hosts define the host to which to bind the server, defaults to "0.0.0.0".
		Host string
		// Port define the port to which to bind the server, defaults to 5000.
		Port int
	}
}

func readConfig(configFile string) *config {
	cfg := &config{}
	if configFile == "" {
		return cfg
	}

	ui.Logger.Infof("OctoPrint's config file found: %q", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		ui.Logger.Fatal(err)
		return cfg
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		ui.Logger.Fatalf("Error decoding YAML config file %q: %s", configFile, err)
		return cfg
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 5000
	}

	return cfg
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

func getSize() (width, height int) {
	if Resolution == "" {
		return
	}

	parts := strings.SplitN(Resolution, "x", 2)
	if len(parts) != 2 {
		ui.Logger.Fatalf("Malformed %s variable: %q", EnvResolution, Resolution)
		return
	}

	var err error
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		ui.Logger.Fatalf("Malformed %s variable: %q, %s",
			EnvResolution, Resolution, err)
		return
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		ui.Logger.Fatalf("Malformed %s variable: %q, %s",
			EnvResolution, Resolution, err)
		return
	}

	return
}
