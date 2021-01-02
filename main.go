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
	utils.Logger.Debug("OctoScreen - entering main.main()")

	gtk.Init(nil)
	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	utils.DumpEnvironmentVariables()

	if utils.RequiredEnvironmentVariablesAreSet(APIKey) {
		width, height := getSize()
		// width and height come from EnvResolution/OCTOSCREEN_RESOLUTION
		// and aren't required - if not set, ui.New() will use the default
		// values (defined in globalVars.go).
		_ = ui.New(BaseURL, APIKey, width, height)
	} else {
		fatalErrorWindow := ui.CreateFatalErrorWindow("Required environment variable is not set:", utils.NameOfMissingRequiredEnvironmentVariable(APIKey))
		fatalErrorWindow.ShowAll()
	}

	gtk.Main()

	utils.Logger.Debug("OctoScreen - leaving main.main()")
}


func init() {
	utils.Logger.Debug("OctoScreen - entering main.init()")

	ConfigFile = os.Getenv(utils.EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	cfg := readConfig(ConfigFile)
	setApiKey(cfg)

	if !utils.RequiredEnvironmentVariablesAreSet(APIKey) {
		utils.Logger.Error("OctoScreen - main.init() - RequiredEnvironmentVariablesAreSet() returned false")

		utils.Logger.Debug("OctoScreen - leaving main.init()")
		return
	}

	setLogLevel()

	utils.StylePath = os.Getenv(utils.EnvStylePath)
	Resolution = os.Getenv(utils.EnvResolution)
	setBaseUrl(cfg)

	utils.Logger.Debug("OctoScreen - leaving main.init()")
}

func setLogLevel() {
	logLevel := utils.LowerCaseLogLevel()
	switch logLevel {
		case "debug":
			utils.SetLogLevel(logrus.DebugLevel)

		case "info":
			utils.SetLogLevel(logrus.InfoLevel)

		case "warn":
			utils.SetLogLevel(logrus.WarnLevel)

		case "":
			logLevel = "error"
			os.Setenv(utils.EnvLogLevel, "error")
			fallthrough
		case "error":
			utils.SetLogLevel(logrus.ErrorLevel)

		default:
			// unknown log level, so exit
			utils.Logger.Fatalf("main.setLogLevel() - unknown logLevel: %q", logLevel)
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
	} else {
		url := strings.ToLower(BaseURL)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			utils.Logger.Warn("WARNING!  OCTOPRINT_HOST requires the transport protocol ('http://' or 'https://') but is missing.  'http://' is being added to BaseURL.");
			BaseURL = fmt.Sprintf("http://%s", BaseURL)
		}
	}

	utils.Logger.Infof("main.setBaseUrl() - using %q as server address", BaseURL)
}

func setApiKey(cfg *config) {
	utils.Logger.Debug("OctoScreen - entering main.setApiKey()")

	APIKey = os.Getenv(utils.EnvAPIKey)
	if APIKey == "" {
		utils.Logger.Debug("main.setApiKey() - APIKey is empty, now using cfg.API.Key")

		APIKey = cfg.API.Key
	}

	if APIKey == "" {
		utils.Logger.Debug("main.setApiKey() - APIKey is empty!")
	} else {
		obfuscatedApiKey := utils.GetObfuscatedValue(APIKey)
		utils.Logger.Debugf("main.setApiKey() - APIKey is %q", obfuscatedApiKey)
	}

	utils.Logger.Debug("OctoScreen - leaving main.setApiKey()")
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
	utils.Logger.Debug("")
	utils.Logger.Debug("")
	utils.Logger.Debug("entering main.readConfig()")

	cfg := &config{}
	if configFile == "" {
		utils.Logger.Info("main.readConfig() - configFile is empty")

		utils.Logger.Debug("leaving main.readConfig(), returning the default config")
		return cfg
	} else {
		utils.Logger.Infof("Path to OctoPrint's config file: %q", configFile)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		utils.Logger.Fatalf("main.readConfig() - ReadFile() returned an error: %q", err)
	} else {
		utils.Logger.Info("main.readConfig() - ReadFile() succeeded")
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		utils.Logger.Fatalf("main.readConfig() - error decoding YAML config file %q: %s", configFile, err)
	} else {
		utils.Logger.Info("main.readConfig() - YAML config file was decoded")
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}

	utils.Logger.Infof("main.readConfig() - server host is: %q", cfg.Server.Host)

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 5000
	}

	utils.Logger.Infof("main.readConfig() - server port is: %d", cfg.Server.Port)


	utils.Logger.Debug("leaving main.readConfig()")
	utils.Logger.Debug("")
	utils.Logger.Debug("")

	return cfg
}

func findConfigFile() string {
	utils.Logger.Debug("entering main.findConfigFile()")

	if file := doFindConfigFile(homeOctoPi); file != "" {
		utils.Logger.Info("main.findConfigFile() - doFindConfigFile() found a file")

		utils.Logger.Debug("leaving main.findConfigFile(), returning the file")
		return file
	}

	usr, err := user.Current()
	if err != nil {
		utils.LogError("main.findConfigFile()", "Current()", err)

		utils.Logger.Debug("leaving main.findConfigFile(), returning an empty string")
		return ""
	}

	configFile := doFindConfigFile(usr.HomeDir)

	utils.Logger.Debug("leaving main.findConfigFile(), returning configFile")
	return configFile
}

func doFindConfigFile(home string) string {
	utils.Logger.Debug("entering main.doFindConfigFile()")

	path := filepath.Join(home, configLocation)

	if _, err := os.Stat(path); err == nil {
		utils.LogError("main.doFindConfigFile()", "Stat()", err)

		utils.Logger.Debug("leaving main.doFindConfigFile(), returning path")
		return path
	}

	utils.Logger.Debug("leaving main.doFindConfigFile(), returning an empty string")
	return ""
}

func getSize() (width, height int) {
	utils.Logger.Debug("entering main.getSize()")

	if Resolution == "" {
		utils.Logger.Info("main.getSize() - Resolution is empty, returning 0 for width and height, and will default to the default values defined in globalVars.go")

		utils.Logger.Debug("leaving main.getSize()")
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

	utils.Logger.Debug("leaving main.getSize()")
	return
}
