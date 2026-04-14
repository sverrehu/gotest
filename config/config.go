package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	// fields need to have capitalized names to be unmarshalled properly.
	// in the YAML file, the corresponding keys may be lowercase.
	Port int
	Host string
}

func LoadConfig() (Config, error) {
	yamlContents, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(yamlContents, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config)
}
