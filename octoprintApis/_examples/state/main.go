package main

import (
	"fmt"
	"os"

	"github.com/mgutz/logxi/v1"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
)

func main() {
	baseURL, apiKey := os.Args[1], os.Args[2]

	c := octoprintApis.NewClient(baseURL, apiKey)
	printConnectionState(c)
	printTemperature(c)
}

func printConnectionState(c *octoprintApis.Client) {
	r := octoprintApis.ConnectionRequest{}
	s, err := r.Do(c)
	if err != nil {
		log.Error("error requesting connection state: %s", err)
	}

	fmt.Printf("Connection State: %q\n", s.Current.State)
}

func printTemperature(c *octoprintApis.Client) {
	r := octoprintApis.StateRequest{}
	s, err := r.Do(c)
	if err != nil {
		log.Error("error requesting state: %s", err)
	}

	fmt.Println("Current Temperatures:")
	for tool, state := range s.Temperature.Current {
		fmt.Printf("- %s: %.1f°C / %.1f°C\n", tool, state.Actual, state.Target)
	}
}
