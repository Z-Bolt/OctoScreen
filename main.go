package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	standardLog "log"
	"os"
	"os/user"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"regexp"

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
	initSucceeded bool
)


func init() {
	logger.Debug("-")
	logger.Debug("-")
	logger.TraceEnter("OctoScreen - main.init()")
	initSucceeded = false

	ConfigFile = os.Getenv(utils.EnvConfigFile)
	if ConfigFile == "" {
		ConfigFile = findConfigFile()
	}

	cfg := readConfig(ConfigFile)
	if cfg == nil {
		initSucceeded = false
		return
	}

	setApiKey(cfg)
	setLogLevel()
	utils.StylePath = os.Getenv(utils.EnvStylePath)
	Resolution = os.Getenv(utils.EnvResolution)
	setBaseUrl(cfg)

	initSucceeded = true
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

	if initSucceeded != true {
		// readConfig() logs any errors it encounters.  Don't display
		// the error here, because the error could be long, and we don't
		// want to display a long message and have the screen resize.
		fatalErrorWindow := ui.CreateFatalErrorWindow(
			"Initialization failed:",
			"readConfig() failed, see log for errors",
		)
		fatalErrorWindow.ShowAll()
	} else {
		if utils.RequiredEnvironmentVariablesAreSet(APIKey) {
			width, height, err := getSize()
			if err == nil {
				// width and height come from EnvResolution/OCTOSCREEN_RESOLUTION
				// and aren't required - if not set, ui.New() will use the default
				// values (defined in globalVars.go).
				_ = ui.New(BaseURL, APIKey, width, height)
			} else {
				// But if there is an error while parsing OCTOSCREEN_RESOLUTION,
				// then display the error.
				fatalErrorWindow := ui.CreateFatalErrorWindow(
					"getSize() failed",
					err.Error(),
				)
				fatalErrorWindow.ShowAll()
			}
		} else {
			fatalErrorWindow := ui.CreateFatalErrorWindow(
				"Required environment variable is not set:",
				utils.NameOfMissingRequiredEnvironmentVariable(APIKey),
			)
			fatalErrorWindow.ShowAll()
		}
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
			// unknown log level
			logLevel = "error"
			logger.Errorf("main.setLogLevel() - unknown logLevel: %q, defaulting to error", logLevel)
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

	if configFile == "" {
		logger.Info("main.readConfig() - configFile is empty")
		logger.TraceLeave("main.readConfig()")
		return nil
	}

	logger.Infof("Path to OctoPrint's config file: %q", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Errorf("main.readConfig() - ReadFile() returned an error: %q", err)
		logger.TraceLeave("main.readConfig()")
		return nil
	}

	cfg := &config{}
	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		logger.Errorf("main.readConfig() - error decoding YAML config file %q: %s", configFile, err)
		logger.TraceLeave("main.readConfig()")
		return nil
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


func getSize() (width int, height int, err error) {
	logger.TraceEnter("main.getSize()")
	
	if Resolution == "" {
		logger.Info("main.getSize() - Resolution is empty, returning 0 for width and height, and will default to the default values defined in globalVars.go")
		logger.TraceLeave("main.getSize()")
		return
	} else if strings.ToLower(Resolution) == "auto" {
		
		logger.Info("Automatically detecting resolution with 'xrandr'.")
		
		xrandr, err := exec.LookPath( "xrandr" )
		
		if err != nil {
			logger.Error("Unable to determine 'xrandr' executable path.")
      err = errors.New("Unable to determine 'xrandr' executable path.")
			logger.TraceLeave("main.getSize()")
			return
		}
		
		cmd := exec.Command(
			xrandr,
			"-d", ":0.0",
			"--prop",
		);
		
		output, err := cmd.Output()
		
		if err != nil {
			logger.Errorf("When determining resolution, 'xrandr' returned with an error. Output: %s", output)
      err = errors.New(fmt.Sprintf("When determining resolution, 'xrandr' returned with an error. Output: %s", output))
			logger.TraceLeave("main.getSize()")
			return
		}
		
		/*
		There is no real error handeling here, as if `xrandr` executes without
		an error, we are gauranteed an output of this format, so the regex
		won't fail.
		*/
		
		re := regexp.MustCompile(`current ([0-9]+) x ([0-9]+)`)
		matches := re.FindStringSubmatch(string(output))
		
		width,  _ = strconv.Atoi(matches[1])
		height, _ = strconv.Atoi(matches[2])
		
	} else {

		parts := strings.SplitN(Resolution, "x", 2)
		if len(parts) != 2 {
			logger.Error("main.getSize() - SplitN() - len(parts) != 2")
			err = errors.New(fmt.Sprintf("%s is malformed\nvalue: %q", utils.EnvResolution, Resolution))
		}
	
		var err error
		width, err = strconv.Atoi(parts[0])
		if err != nil {
			logger.LogError("main.getSize()", "Atoi(parts[0])", err)
			err = errors.New(fmt.Sprintf("%s is malformed\nAtoi(0) failed\nvalue: %q", utils.EnvResolution, Resolution))
		}
	
		height, err = strconv.Atoi(parts[1])
		if err != nil {
			logger.LogError("main.getSize()", "Atoi(parts[1])", err)
			err = errors.New(fmt.Sprintf("%s is malformed\nAtoi(1) failed\nvalue: %q", utils.EnvResolution, Resolution))
		}
  }
	
  logger.Info("Determined a screen resolution of: %d x %d", width, height)
  
  logger.TraceLeave("main.getSize()")
  return
}
