package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Database DBConfig   `json:"database"`
	Http     HttpConfig `yaml:"http"`
}

type DBConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
}

type HttpConfig struct {
	SessionInactivityTimeout uint `yaml:"sessionInactivityTimeout"`
	SessionTimeout           uint `yaml:"sessionTimeout"`
}

func InitViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.booksys-server")
	viper.AddConfigPath(".")

	return nil
}

func ReadConfig(path string) error {
	if path != "" {
		_, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("Configuration file does not exist: %s", err)
		}
		viper.SetConfigFile(path)
	} else {
		// look for a config in
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using empty config")
			return nil
		} else {
			return err
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

func IsDBConfigured() bool {
	switch {
	case !viper.IsSet("database.protocol"):
		return false
	case !viper.IsSet("database.password"):
		return false
	case !viper.IsSet("database.user"):
		return false
	case !viper.IsSet("database.host"):
		return false
	case !viper.IsSet("database.port"):
		return false
	case !viper.IsSet("database.dbname"):
		return false
	}

	return true
}
