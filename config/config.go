package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres Database
	Mysql    Database
}

type Database struct {
	Host     string
	Port     int
	UserName string
	Password string
	Database string
}

var config *Config

func init() {
	load()
}

func GetConfig() Config {
	if config == nil {
		load()
	}
	return *config
}
func load() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Config can't read")
		return
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Panic("Config can't load")
		return
	}
	config = &conf
	log.Println(config)
}
