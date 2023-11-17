package config

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

var App Application

type Application struct {
	Port    string `yaml:"port" env-default:"8080"`
	BaseUrl string `yaml:"baseUrl" env-default:"/"`
	DbUrl   string `yaml:"db_url"`
}

type cfg struct {
	App Application `yaml:"app"`
}

func LoadConfig(path string) error {
	var conf cfg

	yamlFile, err := os.Open(path)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(yamlFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return err
	}

	App = conf.App

	return conf.validate()
}

func (c cfg) validate() error {

	// Port
	if i, err := strconv.Atoi(c.App.Port); err != nil {
		return fmt.Errorf("invalid format app.port: %s", c.App.Port)
	} else if i < 100 {
		return fmt.Errorf("not allowed value app.port: %s", c.App.Port)
	}

	// BaseUrl
	if _, err := url.Parse(c.App.BaseUrl); err != nil {
		return fmt.Errorf("invalid format app.baseUrl: %s", c.App.BaseUrl)
	}

	if _, err := url.Parse(c.App.DbUrl); err != nil {
		return fmt.Errorf("invalid format app.DbUrl: %s", c.App.DbUrl)
	}

	return nil
}
