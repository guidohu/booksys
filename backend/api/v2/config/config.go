// This package handles configuration. It follows the following principle.
//
// Configuration required globally is stored in a configuration file and can be provided
// with environment variables or flags alternatively.
//
// Priority is: flag, environment variables and then config file.
//
// Configuration that could be different for different tenants will be stored in a database
// to be more flexible in case we implement multi tenancy at some point.
//
// Priority is always the configuration stored in the database.
package config

import (
	"context"
	"fmt"
	"os"
	"server/database"
	"time"

	"golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

type Configuration struct {
	Database DBConfig   `yaml:"database"`
	Http     HttpConfig `yaml:"http"`
	params   *BasicParams
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
	Port                     uint   `yaml:"port"`
	SessionInactivityTimeout uint   `yaml:"sessionInactivityTimeout"`
	SessionTimeout           uint   `yaml:"sessionTimeout"`
	UploadPath               string `yaml:"uploadPath"`
}

type BasicParams struct {
	EnvironmentPrefix string
	ConfigFilePath    string
}

func NewBasicConfiguration(p *BasicParams) (*Configuration, error) {
	// v := viper.New()
	if p.ConfigFilePath != "" {
		_, err := os.Stat(p.ConfigFilePath)
		if err != nil {
			return nil, fmt.Errorf("Configuration file does not exist: %s", err)
		}
		viper.SetConfigFile(p.ConfigFilePath)
	} else {
		// look for a config in
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(".")
	}
	return &Configuration{}, nil
}

func InitViper() error {
	// vipe
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.booksys-server")

	return nil
}

// func ReadConfig(v *viper.Viper) error {
// 	if v.GetString() != "" {
// 		_, err := os.Stat(path)
// 		if err != nil {
// 			return fmt.Errorf("Configuration file does not exist: %s", err)
// 		}
// 		viper.SetConfigFile(path)
// 	} else {
// 		// look for a config in
// 		viper.SetConfigName("config")
// 		viper.SetConfigType("yaml")
// 		viper.AddConfigPath("$HOME/")
// 		viper.AddConfigPath(".")
// 	}

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			log.Println("No config file found, using empty config")
// 			return nil
// 		} else {
// 			return err
// 		}
// 	}

// 	// verify that the configuration matches
// 	// the Configuration struct
// 	conf := &Configuration{}
// 	err = viper.Unmarshal(conf)
// 	if err != nil {
// 		log.Fatalf("Cannot decode configuration into config struct: %v", err)
// 	}
// 	return nil
// }

func WriteConfig() error {
	return viper.WriteConfig()
}

func IsDBConfigured(v *viper.Viper) bool {
	switch {
	case !v.IsSet("database.protocol"):
		return false
	case !v.IsSet("database.password"):
		return false
	case !v.IsSet("database.user"):
		return false
	case !v.IsSet("database.host"):
		return false
	case !v.IsSet("database.port"):
		return false
	case !v.IsSet("database.dbname"):
		return false
	}

	return true
}

func LoadDBConfig(v *viper.Viper, db *database.DBMysql) error {
	p, err := db.GetAllPropertyValues()
	if err != nil {
		slog.Error("Cannot retrieve configuration from database", slog.String("error", err.Error()))
		return err
	}
	for _, property := range p {
		v.Set(property.Property, property.Value)
	}
	return nil
}

func WatchDBConfig(ctx context.Context, db *database.DBMysql, notifyCh chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	lastVersion := database.ConfigurationVersion{}
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()

			return
		case <-ticker.C:
			if db == nil {
				slog.Warn("Watch DB Config: No database connection available.")
				continue
			}
			version, err := db.GetConfigurationVersion()
			if err != nil {
				slog.Warn("Watch DB Config: Cannot retrieve configuration version from database", slog.String("error", err.Error()))
			}
			if lastVersion.Version != version.Version {
				if version.Timestamp == nil {
					version.Timestamp = &time.Time{}
				}
				slog.Info("Watch DB Config: New config version detected", slog.Uint64("version", uint64(version.Version)), slog.String("date", version.Timestamp.String()))
				lastVersion = version
				notifyCh <- struct{}{}
			} else {
				slog.Info("Watch DB Config: No config change.")
			}
		}
	}
}
