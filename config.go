package main

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type FrontMatter map[interface{}]interface{}

type Config struct {
	Out         string `yaml:"out,omitempty"`
	FrontMatter `yaml:"front_matter,omitempty"`
}

func LoadConfig(file io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewEmptyConfig() *Config {
	frontMatter := make(FrontMatter)

	return &Config{
		"",
		frontMatter,
	}
}
