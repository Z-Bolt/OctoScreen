package main

import (
	// "errors"
	"fmt"
	// "io/ioutil"
	standardLog "log"
	// "os"
	// "os/user"
	// "path/filepath"
	"runtime"
	// "strconv"
	// "strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/ui"
	"github.com/Z-Bolt/OctoScreen/utils"

	// "gopkg.in/yaml.v1"
)


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

	startSystemDHeartbeat()

	initializeGtk()

	octoScreenConfig := utils.GetOctoScreenConfigInstance()

	if octoScreenConfig.RequiredConfigsAreSet() != true {
		message := fmt.Sprintf("Required setting is not set: %s", octoScreenConfig.MissingRequiredConfigName())
		panic(message)
	}

	setLogLevel(octoScreenConfig.LogLevel)

	utils.DumpSystemInformation()
	utils.DumpEnvironmentVariables()
	octoScreenConfig.DumpConfigs()

	setCursor(octoScreenConfig.DisplayCursor)

	_ = ui.CreateUi()

	gtk.Main()

	logger.TraceLeave("OctoScreen - main.main()")
	logger.Debug("+")
	logger.Debug("+")
}

func startSystemDHeartbeat() {
	systemDHeartbeat := utils.GetSystemDHeartbeatInstance()
	systemDHeartbeat.Start()
}

func initializeGtk() {
	gtk.Init(nil)
	gtkSettings, _ := gtk.SettingsGetDefault()
	gtkSettings.SetProperty("gtk-application-prefer-dark-theme", true)
}

func setLogLevel(logLevel string) {
	switch logLevel {
		case "debug":
			logger.SetLogLevel(logrus.DebugLevel)

		case "info":
			logger.SetLogLevel(logrus.InfoLevel)

		case "warn":
			logger.SetLogLevel(logrus.WarnLevel)

		case "error":
			logger.SetLogLevel(logrus.ErrorLevel)

		default:
			// unknown log level
			logger.Errorf("main.setLogLevel() - unknown logLevel: %q, defaulting to error", logLevel)
			logLevel = "error"
	}

	standardLog.Printf("main.SetLogLevel() - logLevel is now set to: %q", logLevel)
}

func setCursor(displayCursor bool) {
	// For reference, see "How to turn on a pointer"
	// 	https://github.com/Z-Bolt/OctoScreen/issues/285
	// and "No mouse pointer when running xinit"
	// 	https://www.raspberrypi.org/forums/viewtopic.php?t=139546

	if displayCursor != true {
		return
	}

	window, err := getRootWindow()
	if err != nil {
		return
	}

	cursor, err := getDefaultCursor()
	if err != nil {
		return
	}

	window.SetCursor(cursor)
}

func getRootWindow() (*gdk.Window, error) {
	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return nil, err
	}

	window, err := screen.GetRootWindow()

	return window, err
}

func getDefaultCursor() (*gdk.Cursor, error) {
	display, err := gdk.DisplayGetDefault()
	if err != nil {
		return nil, err
	}

	// Examples of the different cursors can be found at
	// https://developer.gnome.org/gdk3/stable/gdk3-Cursors.html#gdk-cursor-new-from-name
	cursor, err := gdk.CursorNewFromName(display, "default")

	return cursor, err
}
