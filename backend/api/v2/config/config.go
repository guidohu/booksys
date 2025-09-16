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
	"path/filepath"
	"server/database"
	"server/yaml"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	Database   DBConfig   `yaml:"database"`
	Http       HttpConfig `yaml:"http"`
	MyNautique MyNautique `yaml:"mynautique"`
	params     *BasicParams
}

type DBConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type HttpConfig struct {
	Port                     uint   `yaml:"port"`
	SessionInactivityTimeout uint   `yaml:"sessionInactivityTimeout"`
	SessionTimeout           uint   `yaml:"sessionTimeout"`
	UploadPath               string `yaml:"uploadPath"`
}

type MyNautique struct {
	APIKey string `yaml:"api.key"`
}

type BasicParams struct {
	EnvironmentPrefix string
	ConfigFilePath    string
}

// GetKeys returns a sorted list of all the configuration fields that are supported by
// the Configuration struct.
func GetKeys() []string {
	var c *Configuration
	tags := yaml.GetKeys(c, "")
	sort.Strings(tags)
	return tags
}

// GetKeysMap returns a map that contains all the supported keys for the Configuration
// struct.
func GetKeysMap() map[string]bool {
	keys := GetKeys()
	kmap := make(map[string]bool)
	for _, k := range keys {
		kmap[k] = true
	}
	return kmap
}

// IsDBConfigured returns false in case any of the mandatory database settings is not provided
// through the configuration.
func (c *Config) IsDBConfigured() bool {
	if !c.IsSet("database.protocol") {
		slog.Info("Database protocol is not set in configuration")
		return false
	}
	if !c.IsSet("database.user") {
		slog.Info("Databse user is not set in configuration")
		return false
	}
	if !c.IsSet("database.host") {
		slog.Info("Database host is not set in configuration")
		return false
	}
	if !c.IsSet("database.port") {
		slog.Info("Database port is not set in configuration")
		return false
	}
	if !c.IsSet("database.dbname") {
		slog.Info("Database name is not set in configuration")
		return false
	}
	return true
}

type ConfigSource int

const (
	Unknown ConfigSource = iota
	Flags
	Environment
	File
	Database
	Defaults
)

type Config struct {
	mu          sync.Mutex
	properties  map[string]database.Configuration
	flags       *viper.Viper
	environment *viper.Viper
	file        *viper.Viper
	defaults    *viper.Viper
	db          *database.Manager
}

type DBConfigWatcher struct {
	mu             sync.Mutex
	cancelPrevious context.CancelFunc
}

// MandatoryConfigKeys are the minimum set of configuration
// settings that need to be present in the configuration.
var MandatoryConfigKeys = []string{
	"http.port",
	"http.sessioninactivitytimeout",
	"http.sessiontimeout",
	"http.uploadpath",
	"database.user",
	"database.protocol",
	"database.host",
	"database.port",
	"database.dbname",
}

// ConfigDefaults are the defaults for each of the parameters and will be
// applied if none of the higher priority configuration sources define this.
var ConfigDefaults map[string]string = map[string]string{
	"config":                        "",
	"http.port":                     "80",
	"http.sessioninactivitytimeout": "604800",
	"http.sessiontimeout":           "31536000",
	"http.uploadpath":               "./uploads",
	"database.user":                 "",
	"database.password":             "",
	"database.protocol":             "tcp",
	"database.host":                 "127.0.0.1",
	"database.port":                 "3306",
	"database.dbname":               "",
	"mynautique.api.key":            "",
}

// NewConfig creates a new Config instance that is initialized with the
// configuration present in flags, environment var, config file or
// properties from the database.
func NewConfig(flagConfig *viper.Viper) (*Config, error) {
	// We manually handle different config sources in the
	// priority order that we want. spf13/viper does not support
	// custom order and mysql as config source.
	// In decreasing priority:
	// - flags: are static and parsed
	// - environment vars: read once and immutable
	// - file: mutable and read upon change
	// - properties from database: mutable and continuously checked
	// - default: static
	c := &Config{
		properties:  map[string]database.Configuration{},
		flags:       flagConfig,
		environment: viper.New(),
		file:        viper.New(),
		defaults:    viper.New(),
	}

	// read environment variables
	// e.g. the flag http.port becomes BOOKSYS_HTTP_PORT
	c.environment.SetEnvPrefix("BOOKSYS")
	c.environment.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.environment.AutomaticEnv()
	for _, key := range GetKeys() {
		c.environment.BindEnv(key)
	}

	// parse configuration file.
	// Note: We get the configuration from flags, env or default
	err := c.ReadConfigFile()
	if err != nil {
		return c, err
	}

	// fill defaults
	for key, value := range ConfigDefaults {
		c.defaults.SetDefault(key, value)
	}

	// Read database properties (if possible)
	c.ReadConfigProperties()

	// Check that we have the mandatory configuration
	// Note: We do not check for precense if a default
	// has been defined.
	return c, nil
}

// findConfigFile tries to find the configuration file in one of the following locations:
// - ~/.booksys
// - ./config.yaml
// - /etc/booksys/config.yaml
func findConfigFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Cannot get user home directory", slog.String("error", err.Error()))
	}
	pathsToSearch := []string{
		filepath.Join(homeDir, ".booksys"),
		filepath.Join(".", "config.yaml"),
		"/etc/booksys/config.yaml",
	}
	// Iterate through the paths and check for existence and readability.
	for _, path := range pathsToSearch {
		info, err := os.Stat(path)
		if err == nil {
			if !info.IsDir() {
				return path
			}
		}
	}

	// No configuration file was found in any of the specified locations.
	return ""
}

// ReadConfigFile reads the configuration file from the provided configuration location
// through the `config` flag or from any of the following locations:
// - ~/.booksys
// - ${PWD}/config.yaml
// - /etc/booksys/config.yaml
func (c *Config) ReadConfigFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	configFile, _ := c.GetString("config")
	if configFile == "" {
		configFile = findConfigFile()
	}
	if configFile != "" {
		slog.Info("Read config from", slog.String("config", configFile))
		_, err := os.Stat(configFile)
		if err != nil {
			slog.Error("File does not exist", slog.String("config", configFile), slog.String("error", err.Error()))
		}
		c.file.SetConfigFile(configFile)
		err = c.file.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				slog.Error("No config file found", slog.String("config", configFile))
				return err
			}
		}
		// Verify that the configuration matches the Configuration struct.
		conf := &Configuration{}
		err = c.file.Unmarshal(conf)
		if err != nil {
			slog.Error("Invalid configuration. Cannot decode configuration.", slog.String("error", err.Error()))
			return err
		}
	} else {
		slog.Info("No config file provided and none found in default directories.")
	}
	return nil
}

// SetDB the database that should be used to read the configuration properties from.
func (c *Config) SetDB(db *database.Manager) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.db = db
}

// SetConfigFileValue this will set the in memory representation of a config value. To persist it to
// disk, you should call WriteConfigFile afterwards.
func (c *Config) SetConfigFileValue(key string, value interface{}) error {
	allowedKeys := GetKeysMap()
	if _, ok := allowedKeys[key]; !ok {
		return fmt.Errorf("cannot set config key/value pair, invalid key %s", key)
	}
	c.file.Set(key, value)
	return nil
}

// SetPropertyValue sets a configuration property in the database.
func (c *Config) SetPropertyValue(key string, value string) error {
	properties := []database.Configuration{
		{
			Property: key,
			Value:    value,
		},
	}
	db, done := c.db.GetHandler()
	if db == nil {
		return fmt.Errorf("no database connection available")
	}
	defer done()
	err := db.UpdateOrInsertPropertyValues(properties)
	return err
}

// WriteConfigFile persists the current in memory state of the configuration file to
// the file on disk.
func (c *Config) WriteConfigFile() error {
	// if file does not exist, we try to create it first.
	location, _ := c.GetString("config")
	if location == "" {
		return fmt.Errorf("no config file path provided to store config")
	}
	if _, err := os.Stat(location); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(location), 0755)
			if err != nil {
				return fmt.Errorf("cannot create directory for config file path: %v", err)
			}
			file, err := os.Create(location)
			if err != nil {
				return err
			}
			defer file.Close()
		} else {
			return fmt.Errorf("cannot access config file path: %v", err)
		}
	}
	c.file.SetConfigFile(location)
	return c.file.WriteConfig()
}

// ReadConfigProperties reads all the properties from the database.
func (c *Config) ReadConfigProperties() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.readConfigProperties()
}

// readConfigProperties reads all the properties from the database and
// updates the internal state of the properties within the Config object.
// The caller MUST hold mu.
func (c *Config) readConfigProperties() error {
	if c.db == nil {
		return fmt.Errorf("no database manager present yet")
	}
	db, done := c.db.GetHandler()
	defer done()
	if db == nil {
		return fmt.Errorf("no database configured")
	}
	err := db.Ping()
	if err != nil {
		return err
	}
	props, err := db.GetAllPropertyValuesMap()
	if err != nil {
		return err
	}
	c.properties = props
	return nil
}

// ToStringFull returns a string representation of the entire configuration.
func (c *Config) ToStringFull() string {
	s := strings.Builder{}

	selectedKeys := map[string]bool{}
	selectedIdentifier := "[selected]"

	// Get all the config keys according to
	// the Configuration struct.
	keys := GetKeys()

	printerFunc := func(v *viper.Viper) {
		for _, key := range keys {
			identifier := ""
			value := ""
			_, ok := selectedKeys[key]
			if !ok && v.IsSet(key) {
				identifier = selectedIdentifier
				selectedKeys[key] = true
			}
			if v.IsSet(key) {
				value = v.GetString(key)
			} else {
				value = "<not set>"
			}
			s.WriteString(fmt.Sprintf("%-30s:\t%-20s%-10s\n", key, value, identifier))
		}
	}

	// print flags
	s.WriteString("Args:\n")
	s.WriteString("-----------\n")
	printerFunc(c.flags)

	// print environment
	s.WriteString("\nEnvironment:\n")
	s.WriteString("-----------\n")
	printerFunc(c.environment)

	// print file content
	configFile, _ := c.GetString("config")
	if configFile == "" {
		s.WriteString(fmt.Sprintf("\nConfig File '%s':\n", "<not set>"))
	} else {
		s.WriteString(fmt.Sprintf("\nConfig File '%s':\n", configFile))
	}
	s.WriteString("-----------\n")
	printerFunc(c.file)

	// print database properties
	s.WriteString("\nDatabase properties:\n")
	s.WriteString("-----------\n")
	dbKeys := []string{}
	for key := range database.AllowedProperties {
		dbKeys = append(dbKeys, key)
	}
	sort.Strings(dbKeys)
	for _, key := range dbKeys {
		v := c.properties[key]
		value := v.Value
		if value == "" {
			value = "<not set>"
		}
		s.WriteString(fmt.Sprintf("%-30s:\t%s\n", key, value))
	}

	// print defaults
	s.WriteString("\nDefaults:\n")
	s.WriteString("-----------\n")
	printerFunc(c.defaults)

	return s.String()
}

// GetString returns a given configuration value as string respecting
// the priority of the configuration sources.
func (c *Config) GetString(key string) (string, ConfigSource) {
	propValue, propOk := c.properties[key]
	switch {
	case c.flags.IsSet(key):
		return c.flags.GetString(key), Flags
	case c.environment.IsSet(key):
		return c.environment.GetString(key), Environment
	case c.file.IsSet(key):
		return c.file.GetString(key), File
	case propOk:
		return propValue.Value, Database
	case c.defaults.IsSet(key):
		return c.defaults.GetString(key), File
	}
	return "", Unknown
}

// GetInt64 returns a given configuration value as integer respecting
// the priority of the configuration sources.
func (c *Config) GetInt64(key string) int64 {
	value, _ := c.GetString(key)
	i, _ := strconv.ParseInt(value, 10, 64)
	return i
}

// GetBool returns a given configuration value as string respecting
// the priority of the configuration sources.
func (c *Config) GetBool(key string) bool {
	value, _ := c.GetString(key)
	return value == "true"
}

// IsSet returns whether a given configuration value is set.
func (c *Config) IsSet(key string) bool {
	_, propOk := c.properties[key]
	switch {
	case c.flags.IsSet(key):
		return true
	case c.environment.IsSet(key):
		return true
	case c.file.IsSet(key):
		return true
	case propOk:
		return true
	case c.defaults.IsSet(key):
		return true
	}
	return false
}

// WatchFile reloads the configuration and notifies whenever
// the configuration file has been updated. This function will
// return a channel to listen on for changes.
func (c *Config) WatchFile() chan struct{} {
	notifyCh := make(chan struct{})
	c.file.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Configuration file changed", slog.String("file", e.Name))
		err := c.ReadConfigFile()
		if err != nil {
			slog.Warn("Configuration cannot be read", slog.String("error", err.Error()))
		}
		notifyCh <- struct{}{}
	})
	return notifyCh
}

// Watch the properties in the database and update the configuration state in case
// the database content gets updated. This function runs until the context is cancelled.
// If it cannot connect to the database it keeps trying.
func (c *Config) WatchProperties(ctx context.Context, notifyCh chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	var lastVersion database.ConfigurationVersion
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			slog.Info("WatchProperties: attempt to read properties from database.")
			func() {
				c.mu.Lock()
				defer c.mu.Unlock()
				dbm := c.db
				if dbm == nil {
					slog.Warn("WatchProperties: no database manager available.")
					return
				}
				db, done := dbm.GetHandler()
				defer done()
				if db == nil {
					slog.Warn("WatchProperties: no database available.")
					return
				}
				version, err := db.GetConfigurationVersion()
				if err != nil {
					slog.Warn("WatchProperties: cannot retrieve configuration version from database", slog.String("error", err.Error()))
					return
				}
				if lastVersion.Version != version.Version || len(c.properties) == 0 {
					if version.Timestamp == nil {
						version.Timestamp = &time.Time{}
					}
					slog.Info("WatchProperties: new config version detected", slog.Uint64("version", uint64(version.Version)), slog.String("date", version.Timestamp.String()))
					err = c.readConfigProperties()
					if err != nil {
						slog.Warn("WatchProperties: cannot read config properties", slog.String("error", err.Error()))
						return
					}
					lastVersion = version
					// Notify the potential listeners in a non-blocking way.
					select {
					case notifyCh <- struct{}{}:
						// write was successful
					default:
						// channel is not ready to receive update
					}
				} else {
					slog.Info("WatchProperties: no changes")
				}
			}()

		}
	}
}
