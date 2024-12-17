package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hmarron/keymap/internal"
)

const DefaultDescription = "MISSING DESCRIPTION"

func main() {
	var (
		ignoreMissing bool
		ignorePlug    bool
		configFile    string
	)
	flag.BoolVar(&ignoreMissing, "ignore-missing", false, "Ignore keybindings with missing descriptions")
	flag.BoolVar(&ignorePlug, "ignore-plug", false, "Ignore <Plug> keybindings")
	flag.StringVar(&configFile, "config-file", "", "Path to write config to")
	flag.Parse()

	// Get the mapping from nvim.
	cmd := exec.Command("nvim", "--headless", "+redir => g:output | silent map | redir END | echo g:output", "+q")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// Parse the output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	keymaps := make(map[string]string)

	// Parse each line
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if isKeyBinding(line) {
			description := DefaultDescription
			if !isKeyBinding(lines[i+1]) {
				description = strings.TrimSpace(lines[i+1])
			}

			if description == DefaultDescription && ignoreMissing {
				continue
			}
			parts := strings.Fields(line)
			mode := parts[0]
			binding := parts[1]

			if strings.HasPrefix(binding, "<Plug>") && ignorePlug {
				continue
			}

			keymaps[fmt.Sprintf("%s %s", mode, binding)] = description
		}
	}

	// Write the out put to the config file or stdout
	if configFile != "" {
		err := internal.WriteConfig(keymaps, configFile, "nvim")
		if err != nil {
			panic(err)
		}
	} else {
		for k, v := range keymaps {
			fmt.Printf("- \"%s: %s\"\n", k, v)
		}
	}

}

func isKeyBinding(str string) bool {
	return strings.HasPrefix(str, "n ") || strings.HasPrefix(str, "o ") || strings.HasPrefix(str, "x ")
}
