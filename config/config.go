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

var cfg Config
var cfgLoaded bool = false

func Cfg() Config {
	if !cfgLoaded {
		panic("config file not loaded")
	}
	return cfg
}

func LoadConfig(fileName string) error {
	if cfgLoaded {
		panic("config file already loaded")
	}
	cfg = Config{}
	err := loadBytesInto(&cfg, defaultConfig)
	if err != nil {
		return err
	}
	if fileName != "" {
		err = loadFileInto(&cfg, fileName)
		if err != nil {
			return err
		}
	}
	cfgLoaded = true
	return nil
}

func loadFileInto(config *Config, filename string) error {
	yamlContents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return loadBytesInto(config, yamlContents)
}

func loadBytesInto(config *Config, bytes []byte) error {
	return yaml.Unmarshal(bytes, config)
}

func main() {
	err := LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", Cfg())
}
