package main

import (
	"fmt"
	"io/ioutil"
	standardLog "log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/ui"
	"github.com/Z-Bolt/OctoScreen/utils"

	"gopkg.in/yaml.v1"
)


var (
	BaseURL    string
	APIKey     string
	ConfigFile string
	Resolution string
)


func init() {
	logger.Debug("-")
	logger.Debug("-")
	// logger.Debug("OctoScreen - entering main.init()")
	logger.TraceEnter("OctoScreen - main.init()")

	ConfigFile = os.Getenv(utils.EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	cfg := readConfig(ConfigFile)
	setApiKey(cfg)

	if !utils.RequiredEnvironmentVariablesAreSet(APIKey) {
		logger.Error("OctoScreen - main.init() - RequiredEnvironmentVariablesAreSet() returned false")
		// logger.Debug("OctoScreen - leaving main.init()")
		logger.TraceLeave("OctoScreen - main.init()")
		return
	}

	setLogLevel()

	utils.StylePath = os.Getenv(utils.EnvStylePath)
	Resolution = os.Getenv(utils.EnvResolution)
	setBaseUrl(cfg)

	// logger.Debug("OctoScreen - leaving main.init()")
	logger.TraceLeave("OctoScreen - main.init()")
	logger.Debug("-")
	logger.Debug("-")
}


func main() {
	defer func() {
		standardLog.Println("main's defer() was called, now calling recover()")
		rec := recover();
		if rec != nil {
			standardLog.Println("main's defer() - recover:", rec)
		} else {
			standardLog.Println("main's defer() - recover was nil")
		}

		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		/*
		programCounter, fileName, lineNumber, infoWasRecovered := runtime.Caller(2)
		standardLog.Println("main's defer() - programCounter:", programCounter)
		standardLog.Println("main's defer() - fileName:", fileName)
		standardLog.Println("main's defer() - lineNumber:", lineNumber)
		standardLog.Println("main's defer() - infoWasRecovered:", infoWasRecovered)
		*/

		pc := make([]uintptr, 20)
		numberOfPcEntries := runtime.Callers(0, pc)
		if numberOfPcEntries > 10 {
			numberOfPcEntries = 10
		}

		for i := 1; i < numberOfPcEntries; i++ {
			/*
			standardLog.Printf("main's defer() - [%d]", i)
			standardLog.Printf("main's defer() - [%d]", numberOfPcEntries)

			programCounter, fileName, lineNumber, infoWasRecovered := runtime.Caller(i)
			standardLog.Printf("main's defer() - programCounter[%d]: %v", i, programCounter)
			standardLog.Printf("main's defer() - fileName[%d]: %v", i, fileName)
			standardLog.Printf("main's defer() - lineNumber[%d]: %v", i, lineNumber)
			standardLog.Printf("main's defer() - infoWasRecovered[%d]: %v", i, infoWasRecovered)
			standardLog.Println("")
			*/

			_, fileName, lineNumber, infoWasRecovered := runtime.Caller(i)
			if infoWasRecovered {
				standardLog.Printf("main's defer() - [%d] %s, line %d", i, fileName, lineNumber)
			}
		}

		standardLog.Println("main's defer() was called, now exiting func()")
	}()


	logger.Debug("+")
	logger.Debug("+")
	logger.TraceEnter("OctoScreen - main.main()")

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

	logger.TraceLeave("OctoScreen - main.main()")
	logger.Debug("+")
	logger.Debug("+")
}


func setLogLevel() {
	logLevel := strings.ToLower(os.Getenv(utils.EnvLogLevel))

	switch logLevel {
		case "debug":
			logger.SetLogLevel(logrus.DebugLevel)

		case "info":
			logger.SetLogLevel(logrus.InfoLevel)

		case "warn":
			logger.SetLogLevel(logrus.WarnLevel)

		case "":
			logLevel = "error"
			os.Setenv(utils.EnvLogLevel, "error")
			fallthrough
		case "error":
			logger.SetLogLevel(logrus.ErrorLevel)

		default:
			// unknown log level, so exit
			logger.Fatalf("main.setLogLevel() - unknown logLevel: %q", logLevel)
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
			logger.Warn("WARNING!  OCTOPRINT_HOST requires the transport protocol ('http://' or 'https://') but is missing.  'http://' is being added to BaseURL.");
			BaseURL = fmt.Sprintf("http://%s", BaseURL)
		}
	}

	logger.Infof("main.setBaseUrl() - using %q as server address", BaseURL)
}

func setApiKey(cfg *config) {
	logger.TraceEnter("main.setApiKey()")

	APIKey = os.Getenv(utils.EnvAPIKey)
	if APIKey == "" {
		logger.Debug("main.setApiKey() - APIKey is empty, now using cfg.API.Key")
		APIKey = cfg.API.Key
	}

	if APIKey == "" {
		logger.Debug("main.setApiKey() - APIKey is empty!")
	} else {
		obfuscatedApiKey := utils.GetObfuscatedValue(APIKey)
		logger.Debugf("main.setApiKey() - APIKey is %q", obfuscatedApiKey)
	}

	logger.TraceLeave("main.setApiKey()")
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
	logger.TraceEnter("main.readConfig()")

	cfg := &config{}
	if configFile == "" {
		logger.Info("main.readConfig() - configFile is empty")
		logger.TraceLeave("main.readConfig(), returning the default config")
		return cfg
	} else {
		logger.Infof("Path to OctoPrint's config file: %q", configFile)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Fatalf("main.readConfig() - ReadFile() returned an error: %q", err)
	} else {
		logger.Info("main.readConfig() - ReadFile() succeeded")
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		logger.Fatalf("main.readConfig() - error decoding YAML config file %q: %s", configFile, err)
	} else {
		logger.Info("main.readConfig() - YAML config file was decoded")
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}

	logger.Infof("main.readConfig() - server host is: %q", cfg.Server.Host)

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 5000
	}

	logger.Infof("main.readConfig() - server port is: %d", cfg.Server.Port)

	logger.TraceLeave("main.readConfig()")
	return cfg
}

func findConfigFile() string {
	logger.TraceEnter("main.findConfigFile()")

	if file := doFindConfigFile(homeOctoPi); file != "" {
		logger.Info("main.findConfigFile() - doFindConfigFile() found a file")
		logger.TraceLeave("main.findConfigFile(), returning the file")
		return file
	}

	usr, err := user.Current()
	if err != nil {
		logger.LogError("main.findConfigFile()", "Current()", err)
		logger.TraceLeave("main.findConfigFile(), returning an empty string")
		return ""
	}

	configFile := doFindConfigFile(usr.HomeDir)

	logger.TraceLeave("main.findConfigFile(), returning configFile")
	return configFile
}

func doFindConfigFile(home string) string {
	logger.TraceEnter("main.doFindConfigFile()")

	path := filepath.Join(home, configLocation)
	if _, err := os.Stat(path); err == nil {
		logger.LogError("main.doFindConfigFile()", "Stat()", err)
		logger.TraceLeave("main.doFindConfigFile(), returning path")
		return path
	}

	logger.TraceLeave("main.doFindConfigFile(), returning an empty string")
	return ""
}







func getSize() (width, height int) {
	logger.TraceEnter("main.getSize()")

	if Resolution == "" {
		logger.Info("main.getSize() - Resolution is empty, returning 0 for width and height, and will default to the default values defined in globalVars.go")
		logger.TraceLeave("main.getSize()")
		return
	}

	parts := strings.SplitN(Resolution, "x", 2)
	if len(parts) != 2 {
		logger.Error("main.getSize() - SplitN() - len(parts) != 2")
		logger.Fatalf("main.getSize() - malformed %s variable: %q", utils.EnvResolution, Resolution)
	}

	var err error
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		logger.LogError("main.getSize()", "Atoi(parts[0])", err)
		logger.Fatalf("main.getSize() - malformed %s variable: %q, %s", utils.EnvResolution, Resolution, err)
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		logger.LogError("main.getSize()", "Atoi(parts[1])", err)
		logger.Fatalf("main.getSize() - malformed %s variable: %q, %s", utils.EnvResolution, Resolution, err)
	}

	logger.TraceLeave("main.getSize()")
	return
}
