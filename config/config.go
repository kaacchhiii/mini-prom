package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Target struct {
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type StorageConfig struct {
	Retention time.Duration `yaml:"retention"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Targets []Target      `yaml:"targets"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
