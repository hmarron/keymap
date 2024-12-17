package internal

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func WriteConfig(bindings map[string]string, path string, frameName string) error {
	// Convert the bindings to a slice of strings.
	bindingsSlice := make([]string, 0, len(bindings))
	for k, v := range bindings {
		bindingsSlice = append(bindingsSlice, fmt.Sprintf("%s: %s", k, v))
	}

	// Read the YAML file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into a map
	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err

	}
	// Get the contents of existing frame mapping if there is one
	// Access and modify specific fields
	if frames, ok := config["frames"].(map[interface{}]interface{}); ok {
		frames[frameName] = bindingsSlice
	}
	// Marshal the map back to YAML
	newData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write the modified YAML data back to the file
	err = ioutil.WriteFile(path, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}
