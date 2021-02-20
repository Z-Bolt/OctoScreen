package octoprintApis

import (
	"log"
)



var EnableApiLogging = false


func LogError(err error, msg string) {
	if EnableApiLogging {
		log.Printf("ERROR!!! %s", err.Error())
		log.Printf("Message: %s", msg)
	}
}


func LogMessage(v ...interface{}) {
	if EnableApiLogging {
		log.Print(v...)
	}
}

func LogMessagef(format string, v ...interface{}) {
	if EnableApiLogging {
		log.Printf(format, v...)
	}
}

func LogMessageln(v ...interface{}) {
	if EnableApiLogging {
		log.Println(v...)
	}
}
