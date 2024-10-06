package config

import (
	"fmt"
	"os"

	"github.com/romshark/yamagiconf"
)

type Config struct {
	Host string `yaml:"host"`
}

func MustLoad(filePath string) *Config {
	var c Config
	err := yamagiconf.LoadFile(filePath, &c)
	if err != nil {
		fmt.Printf("reading config: %v\n", err)
		os.Exit(1)
	}
	return &c
}
