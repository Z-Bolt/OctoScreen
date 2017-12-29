package main

import (
	"fmt"
	"os"

	"github.com/mcuadros/go-octoprint"
	"github.com/mgutz/logxi/v1"
)

func main() {
	baseURL, apiKey := os.Args[1], os.Args[2]

	c := octoprint.NewClient(baseURL, apiKey)
	printConnectionState(c)
	printTemperature(c)
}

func printConnectionState(c *octoprint.Client) {
	r := octoprint.ConnectionRequest{}
	s, err := r.Do(c)
	if err != nil {
		log.Error("error requesting connection state: %s", err)
	}

	fmt.Printf("Connection State: %q\n", s.Current.State)
}

func printTemperature(c *octoprint.Client) {
	r := octoprint.StateRequest{}
	s, err := r.Do(c)
	if err != nil {
		log.Error("error requesting state: %s", err)
	}

	fmt.Println("Current Temperatures:")
	for tool, state := range s.Temperature.Current {
		fmt.Printf("- %s: %.1f°C / %.1f°C\n", tool, state.Actual, state.Target)
	}
}
