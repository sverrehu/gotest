package main

import (
	_ "embed"
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

//go:embed config.default.yaml
var defaultConfig []byte

type Credentials struct {
	UserName string `yaml:"userName,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

type Config struct {
	WebServer struct {
		Port int `yaml:"port"`
	} `yaml:"webServer"`
	Credentials map[string]Credentials `yaml:"credentials"`
}

func LoadConfig() (Config, error) {
	config := Config{}
	err := loadConfigFileInto(&config, "config.default.yaml")
	if err != nil {
		return Config{}, err
	}
	err = loadConfigFileInto(&config, "config.yaml")
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func loadConfigFileInto(config *Config, filename string) error {
	yamlContents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return loadConfigBytesInto(config, yamlContents)
}

func loadConfigBytesInto(config *Config, bytes []byte) error {
	err := yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config)
}
