package config

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Application struct {
	Port    string `yaml:"port" env-default:"8080"`
	BaseUrl string `yaml:"baseUrl" env-default:"/"`
}

var Conf Config

type Config struct {
	App Application `yaml:"app"`
}

func LoadConfig(path string) error {

	yamlFile, err := os.Open(path)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(yamlFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &Conf)
	if err != nil {
		return err
	}
	return Conf.validate()
}

func (c *Config) validate() error {

	// Port
	if i, err := strconv.Atoi(c.App.Port); err != nil {
		return fmt.Errorf("invalid format app.port: %s", c.App.Port)
	} else if i < 100 {
		return fmt.Errorf("not allowed value app.port: %s", c.App.Port)
	}

	// BaseUrl
	if _, err := url.Parse(c.App.BaseUrl); err != nil {
		return fmt.Errorf("invalid fortma app.baseUrl: %s", c.App.BaseUrl)
	}

	return nil
}
