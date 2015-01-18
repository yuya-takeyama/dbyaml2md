package main

import (
	"gopkg.in/yaml.v2"
	"github.com/yuya-takeyama/db2yaml/model"
	"io"
	"io/ioutil"
)

type Config struct {
	FrontMatter map[interface{}]interface{} `yaml:"front_matter,omitempty"`
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
	frontMatter := make(map[interface{}]interface{})

	return &Config{frontMatter}
}

func (config *Config) FrontMatterWithTableName(table *model.Table) map[interface{}]interface{} {
	frontMatter := config.FrontMatter
	frontMatter["table"] = table.Name

	return frontMatter
}

func (config *Config) FrontMatterWithTableNameBytes(table *model.Table) ([]byte, error) {
	return yaml.Marshal(config.FrontMatterWithTableName(table))
}
