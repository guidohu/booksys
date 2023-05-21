package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	Database DBConfig `json:"database"`
}

type DBConfig struct {
	host     string `yaml:"host"`
	port     string `yaml:"port"`
	user     string `yaml:"user"`
	password string `yaml:"password"`
}

func InitViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.booksys-server")
	viper.AddConfigPath(".")

	return nil
}

func ReadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.booksys-server")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using empty config")
			return nil
		} else {
			log.Fatalf("Cannot read configuration file config.yaml")
		}
	}

	// verify that the configuration matches
	// the Configuration struct
	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("Cannot decode configuration into config struct: %v", err)
	}

	return nil
}

func WriteConfig() error {
	return viper.WriteConfig()
}
