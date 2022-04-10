package utils

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/Z-Bolt/OctoScreen/logger"
)


type BackgroundTask struct {
	sync.Mutex

	duration			time.Duration
	task				func()
	close				chan bool
	isRunning			bool
}

func CreateBackgroundTask(
	duration			time.Duration,
	task				func(),
) *BackgroundTask {
	thisInstance := &BackgroundTask{
		task:			task,
		duration: 		duration,
		close: 			make(chan bool, 1),
		isRunning:		false,
	}

	return thisInstance
}

func (this *BackgroundTask) Start() {
	this.Lock()
	defer this.Unlock()

	logger.Info("New background task started")
	go this.loop()

	this.isRunning = true
}

func (this *BackgroundTask) Close() {
	this.Lock()
	defer this.Unlock()
	if !this.isRunning {
		return
	}

	this.close <- true
	this.isRunning = false
}

func (this *BackgroundTask) loop() {
	this.execute()

	ticker := time.NewTicker(this.duration)
	defer ticker.Stop()
	for {
		select {
			case <-ticker.C:
				this.execute()

			case <-this.close:
				logger.Info("Background task closed")
				return
		}
	}
}

func (this *BackgroundTask) execute() {
	_, err := glib.IdleAdd(this.task)
	if err != nil {
		logger.LogFatalError("BackgroundTask.execute()", "IdleAdd()", err)
	}
}


func GetExperimentalFrequency(
	defaultTimeout				int,
	experimentalConfigName		string,
) time.Duration {
	duration := time.Second * time.Duration(defaultTimeout)

	// Experimental, set the timeout based on config setting, but only if the config is pressent.
	updateFrequency := os.Getenv(experimentalConfigName)
	if updateFrequency != "" {
		logger.Infof(
			"BackgroundTask.GetExperimentalFrequency() - '%s' is present, frequency is %s",
			experimentalConfigName,
			updateFrequency,
		)
		val, err := strconv.Atoi(updateFrequency)
		if err == nil {
			duration = time.Second * time.Duration(val)
		} else {
			logger.LogError("BackgroundTask.GetExperimentalFrequency()", "strconv.Atoi()", err)
		}
	}
	
	return duration
}