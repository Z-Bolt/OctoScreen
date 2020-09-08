package main

import (
	"fmt"
	"io/ioutil"
	standardLog "log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Z-Bolt/OctoScreen/ui"
	"github.com/Z-Bolt/OctoScreen/utils"
	"github.com/gotk3/gotk3/gtk"
	"gopkg.in/yaml.v1"

	"github.com/sirupsen/logrus"
)

var (
	BaseURL    string
	APIKey     string
	ConfigFile string
	Resolution string
)

func main() {
	utils.Logger.Info("entering main.main()")

	gtk.Init(nil)
	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	utils.DumpEnvironmentVariables()

	if utils.RequiredEnvironmentVariablesAreSet() {
		width, height := getSize()
		_ = ui.New(BaseURL, APIKey, width, height)
	} else {
		fatalErrorWindow := ui.CreateFatalErrorWindow("Required environment variable is not set:", utils.NameOfMissingRequiredEnvironmentVariable())
		fatalErrorWindow.ShowAll()
	}

	gtk.Main()

	utils.Logger.Info("leaving main.main()")
}


func init() {
	utils.Logger.Info("entering main.init()")

	if !utils.RequiredEnvironmentVariablesAreSet() {
		utils.Logger.Error("main.init() - RequiredEnvironmentVariablesAreSet() returned false")
		utils.Logger.Error("leaving main.init()")
		return
	}

	setLogLevel()

	utils.StylePath = os.Getenv(utils.EnvStylePath)
	Resolution = os.Getenv(utils.EnvResolution)
	ConfigFile = os.Getenv(utils.EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	cfg := readConfig(ConfigFile)
	setBaseUrl(cfg)
	setApiKey(cfg)

	utils.Logger.Info("leaving main.init()")
}

func setLogLevel() {
	standardLog.Print("entering main.setLogLevel()")

	logLevel := os.Getenv(utils.EnvLogLevel)
	switch strings.ToLower(logLevel) {
		case "debug":
			utils.SetLogLevel(logrus.DebugLevel)

		case "info":
			utils.SetLogLevel(logrus.InfoLevel)

		case "":
			fallthrough
		case "warn":
			utils.SetLogLevel(logrus.WarnLevel)

		case "error":
			utils.SetLogLevel(logrus.ErrorLevel)

		default:
			utils.Logger.Fatalf("main.init() - unknown logLevel: %q", logLevel)
	}

	standardLog.Printf("main.SetLogLevel() - logLevel is now set to: %q", logLevel)
}

func setBaseUrl(cfg *config) {
	BaseURL = os.Getenv(utils.EnvBaseURL)
	if BaseURL == "" {
		if cfg.Server.Host != "" {
			BaseURL = fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
		} else {
			BaseURL = "http://0.0.0.0:5000"
		}
	}

	utils.Logger.Infof("main.setBaseUrl() - using %q as server address", BaseURL)
}

func setApiKey(cfg *config) {
	APIKey = os.Getenv(utils.EnvAPIKey)
	if APIKey == "" {
		APIKey = cfg.API.Key
		if cfg.API.Key != "" {
			utils.Logger.Infof("main.setApiKey() - found API key in file %q", ConfigFile)
		}
	}
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
	utils.Logger.Info("entering main.readConfig()")

	cfg := &config{}
	if configFile == "" {
		utils.Logger.Info("main.readConfig() - configFile is empty")
		utils.Logger.Info("leaving main.readConfig(), returning the default config")
		return cfg
	}

	utils.Logger.Infof("OctoPrint's config file was found: %q", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		utils.Logger.Fatalf("main.readConfig() - ReadFile() returned an error: %q", err)
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		utils.Logger.Fatalf("main.readConfig() - error decoding YAML config file %q: %s", configFile, err)
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 5000
	}

	utils.Logger.Info("leaving main.readConfig()")
	return cfg
}

func findConfigFile() string {
	utils.Logger.Info("entering main.findConfigFile()")

	if file := doFindConfigFile(homeOctoPi); file != "" {
		utils.Logger.Info("main.findConfigFile() - doFindConfigFile() found a file")
		utils.Logger.Info("leaving main.findConfigFile(), returning the file")
		return file
	}

	usr, err := user.Current()
	if err != nil {
		utils.LogError("main.findConfigFile()", "Current()", err)
		utils.Logger.Error("leaving main.findConfigFile(), returning an empty string")
		return ""
	}

	configFile := doFindConfigFile(usr.HomeDir)

	utils.Logger.Info("leaving main.findConfigFile(), returning configFile")
	return configFile
}

func doFindConfigFile(home string) string {
	utils.Logger.Info("entering main.doFindConfigFile()")

	path := filepath.Join(home, configLocation)

	if _, err := os.Stat(path); err == nil {
		utils.LogError("main.doFindConfigFile()", "Stat()", err)
		utils.Logger.Error("leaving main.doFindConfigFile(), returning path")
		return path
	}

	utils.Logger.Info("leaving main.doFindConfigFile(), returning an empty string")
	return ""
}

func getSize() (width, height int) {
	utils.Logger.Info("entering main.getSize()")

	if Resolution == "" {
		utils.Logger.Error("main.getSize() - Resolution is empty")
		utils.Logger.Error("leaving main.getSize()")
		return
	}

	parts := strings.SplitN(Resolution, "x", 2)
	if len(parts) != 2 {
		utils.Logger.Error("main.getSize() - SplitN() - len(parts) != 2")
		utils.Logger.Fatalf("main.getSize() - malformed %s variable: %q", utils.EnvResolution, Resolution)
	}

	var err error
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		utils.LogError("main.getSize()", "Atoi(parts[0])", err)
		utils.Logger.Fatalf("main.getSize() - malformed %s variable: %q, %s", utils.EnvResolution, Resolution, err)
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		utils.LogError("main.getSize()", "Atoi(parts[1])", err)
		utils.Logger.Fatalf("main.getSize() - malformed %s variable: %q, %s", utils.EnvResolution, Resolution, err)
	}

	utils.Logger.Info("leaving main.getSize()")
	return
}
