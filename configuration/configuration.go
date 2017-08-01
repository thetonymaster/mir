package configuration

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type conf struct {
	Configuration Configuration `json:"configuration"`
}

type Configuration struct {
	Containers ContainerConfig `json:"containers"`
	Tests      Test            `json:"tests"`
	Provider   string          `json:"provider"`
}

type ContainerConfig struct {
	Limit  int    `json:"limit"`
	Memory string `json:"memory"`
}

type TestConfiguration struct {
	Repo   string `json:"repo"`
	Target string `json:"target"`
	Path   string `json:"path"`
}

type Options struct {
}

type Test map[string]TestConfiguration

func Read(path string) (*Configuration, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg conf

	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Configuration, nil
}
