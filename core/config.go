package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	SessionSecret string `yaml:"session_secret"`
	Addr          string `yaml:"addr"`
	RedisHost     string `yaml:"redis_host"`
	RedisPort     int    `yaml:"redis_port"`
	RedisAuth     string `yaml:"redis_auth"`
	WebsiteName   string `yaml:"website_name"`
	ImageHost     string `yaml:"image_host"`
	DBHost        string `yaml:"db_host"`
	DBPort        int    `yaml:"db_port"`
	DBUser        string `yaml:"db_user"`
	DBPasswd      string `yaml:"db_passwd"`
	DBName        string `yaml:"db_name"`
}

var configuration *config

func InitConfig(filename string) error {
	var err error

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var c config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	configuration = &c
	return err
}

func GetConfig() *config {
	return configuration
}
