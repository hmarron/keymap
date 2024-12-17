package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"

	"github.com/hmarron/keymap/internal"
)

func main() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config-file", "", "Path to write config to")
	flag.Parse()

	mainKeymaps := loadBindings("mode.main.binding")
	serviceKeymaps := loadBindings("mode.service.binding")

	// Write the output to the config file or stdout
	if configFile != "" {
		err := internal.WriteConfig(mainKeymaps, configFile, "aerospace-main")
		if err != nil {
			panic(err)
		}
		err = internal.WriteConfig(serviceKeymaps, configFile, "aerospace-service")
		if err != nil {
			panic(err)
		}

	} else {
		for k, v := range mainKeymaps {
			fmt.Printf("- \"%s: %s\"\n", k, v)
		}
		for k, v := range serviceKeymaps {
			fmt.Printf("- \"%s: %s\"\n", k, v)
		}

	}

}

func loadBindings(key string) map[string]string {
	// Get the mapping from aerospace.
	cmd := exec.Command("aerospace", "config", "--get", key, "--json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// Parse the output
	var result map[string]string
	err = json.Unmarshal(output, &result)
	if err != nil {
		panic(err)
	}

	return result
}
