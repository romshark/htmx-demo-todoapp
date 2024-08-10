package config

import (
	"fmt"
	"os"
	"time"

	"github.com/romshark/yamagiconf"
)

type Config struct {
	Host     string `yaml:"host"`
	Simulate struct {
		ResponseDelayMin time.Duration `yaml:"response-delay-min"`
		ResponseDelayMax time.Duration `yaml:"response-delay-max"`
	} `yaml:"simulate"`
}

func MustLoad(filePath string) *Config {
	var c Config
	err := yamagiconf.LoadFile("config.yaml", &c)
	if err != nil {
		fmt.Printf("reading config: %v\n", err)
		os.Exit(1)
	}
	return &c
}
