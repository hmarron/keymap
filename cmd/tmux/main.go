package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hmarron/keymap/internal"
)

func main() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config-file", "", "Path to write config to")
	flag.Parse()

	// Get the mapping from tmux.
	cmd := exec.Command("tmux", "list-keys", "-N")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// Parse the output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	keymaps := make(map[string]string)

	// Parse each line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		keymaps[fmt.Sprintf("%s %s", parts[0], parts[1])] = strings.Join(parts[2:], " ")
	}

	// Write the output to the config file or stdout
	if configFile != "" {
		err := internal.WriteConfig(keymaps, configFile, "tmux")
		if err != nil {
			panic(err)
		}
	} else {
		for k, v := range keymaps {
			fmt.Printf("- \"%s: %s\"\n", k, v)
		}
	}

}
