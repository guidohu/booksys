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

// Returns a sorted list of all the configuration fields that are supported.
func GetKeys() []string {
	var c *Configuration
	tags := yaml.GetKeys(c, "")
	sort.Strings(tags)
	return tags
}

func GetKeysMap() map[string]bool {
	keys := GetKeys()
	kmap := make(map[string]bool)
	for _, k := range keys {
		kmap[k] = true
	}
	return kmap
}

// func NewBasicConfiguration(p *BasicParams) (*Configuration, error) {
// 	// v := viper.New()
// 	if p.ConfigFilePath != "" {
// 		_, err := os.Stat(p.ConfigFilePath)
// 		if err != nil {
// 			return nil, fmt.Errorf("Configuration file does not exist: %s", err)
// 		}
// 		viper.SetConfigFile(p.ConfigFilePath)
// 	} else {
// 		// look for a config in
// 		viper.SetConfigName("config")
// 		viper.SetConfigType("yaml")
// 		viper.AddConfigPath("$HOME/")
// 		viper.AddConfigPath(".")
// 	}
// 	return &Configuration{}, nil
// }

// func InitViper() error {
// 	// vipe
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath(".")
// 	viper.AddConfigPath("$HOME/.booksys-server")

// 	return nil
// }

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

// func WriteConfig() error {
// 	return viper.WriteConfig()
// }

func (c *Config) IsDBConfigured() bool {
	protocol, _ := c.GetString("database.protocol")
	if protocol == "" {
		slog.Info("Database protocol is not set in configuration")
		return false
	}
	user, _ := c.GetString("database.user")
	if user == "" {
		slog.Info("Databse user is not set in configuration")
		return false
	}
	host, _ := c.GetString("database.host")
	if host == "" {
		slog.Info("Database host is not set in configuration")
		return false
	}
	port, _ := c.GetString("database.port")
	if port == "" {
		slog.Info("Database port is not set in configuration")
		return false
	}
	dbname, _ := c.GetString("database.dbname")
	if dbname == "" {
		slog.Info("Database name is not set in configuration")
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
	db          *database.DBMysql
}

type DBConfigWatcher struct {
	mu             sync.Mutex
	cancelPrevious context.CancelFunc
}

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

var ConfigDefaults map[string]string = map[string]string{
	"config":                        "./config.yaml",
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

func (c *Config) ReadConfigFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	configFile, _ := c.GetString("config")
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
	}
	return nil
}

// Set the database that should be used for the configuration.
func (c *Config) SetDB(db *database.DBMysql) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.db = db
}

func (c *Config) SetConfigFileValue(key string, value interface{}) error {
	allowedKeys := GetKeysMap()
	if _, ok := allowedKeys[key]; !ok {
		return fmt.Errorf("cannot set config key/value pair, invalid key %s", key)
	}
	c.file.Set(key, value)
	return nil
}

func (c *Config) SetPropertyValue(key string, value string) error {
	properties := []database.Configuration{
		{
			Property: key,
			Value:    value,
		},
	}
	err := c.db.UpdateOrInsertPropertyValues(properties)
	return err
}

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

// func (c *Config) connectDB() error {
// 	user, _ := c.GetString("database.user")
// 	password, _ := c.GetString("database.password")
// 	protocol, _ := c.GetString("database.protocol")
// 	host, _ := c.GetString("database.host")
// 	port, _ := c.GetString("database.port")
// 	dbname, _ := c.GetString("database.dbname")
// 	db := &database.DBMysql{
// 		User:     user,
// 		Password: password,
// 		Protocol: protocol,
// 		Host:     host,
// 		Port:     port,
// 		DBName:   dbname,
// 	}
// 	slog.Info("Connecting database client to:", slog.String("address", db.String()))
// 	if err := db.Connect(); err != nil {
// 		slog.Warn(fmt.Sprintf("Database is not properly setup or not reachable. Error returned from Connect(): %s", err))
// 		return err
// 	}
// 	slog.Info("Connected to database", slog.String("name", dbname))

// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	c.db = db
// 	return nil
// }

func (c *Config) ReadConfigProperties() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.readConfigProperties()
}

// readConfigProperties reads all the properties from the database. The
// caller MUST hold mu.
func (c *Config) readConfigProperties() error {
	if c.db == nil {
		return fmt.Errorf("no database configured")
	}
	err := c.db.Ping()
	if err != nil {
		return err
	}
	props, err := c.db.GetAllPropertyValuesMap()
	if err != nil {
		return err
	}
	c.properties = props
	return nil
}

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
	s.WriteString(fmt.Sprintf("\nConfig File '%s':\n", configFile))
	s.WriteString("-----------\n")
	printerFunc(c.file)

	// print database properties
	s.WriteString("\nDatabase properties:\n")
	s.WriteString("-----------\n")
	dbKeys := []string{}
	for key, _ := range database.AllowedProperties {
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

func (c *Config) GetInt64(key string) int64 {
	value, _ := c.GetString(key)
	i, _ := strconv.ParseInt(value, 10, 64)
	return i
}

func (c *Config) GetBool(key string) bool {
	value, _ := c.GetString(key)
	return value == "true"
}

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
// return.
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

// Watch the properties in the database. This function runs until the context is cancelled.
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
			c.mu.Lock()
			if c.db == nil {
				slog.Warn("WatchProperties: no database connection available.")
				c.mu.Unlock()
				continue
			}
			if c.db.Ping() != nil {
				slog.Warn("WatchProperties: database unresponsive.")
				c.mu.Unlock()
				continue
			}
			version, err := c.db.GetConfigurationVersion()
			if err != nil {
				slog.Warn("WatchProperties: cannot retrieve configuration version from database", slog.String("error", err.Error()))
				// TODO this could be because there is no db connection
				continue
			}

			if lastVersion.Version != version.Version || len(c.properties) == 0 {
				if version.Timestamp == nil {
					version.Timestamp = &time.Time{}
				}
				slog.Info("WatchProperties: new config version detected", slog.Uint64("version", uint64(version.Version)), slog.String("date", version.Timestamp.String()))
				err = c.readConfigProperties()
				if err != nil {
					c.mu.Unlock()
					slog.Warn("WatchProperties: cannot read config properties", slog.String("error", err.Error()))
					continue
				}
				lastVersion = version
				slog.Info("DEBUG: version assigned")
				select {
				case notifyCh <- struct{}{}:
					// write was successful
				default:
					// channel is not ready to receive update
					break
				}
				slog.Info("DEBUG: config properties have been read")
			} else {
				slog.Info("WatchProperties: no changes")
			}
			c.mu.Unlock()
		}
	}
}

func getEnvironment(v *viper.Viper) {
	v.SetEnvPrefix("BOOKSYS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func NewDBConfigWatcher() *DBConfigWatcher {
	return &DBConfigWatcher{}
}

func (w *DBConfigWatcher) Watch(ctx context.Context, db *database.DBMysql, notifyCh chan struct{}) {
	// cancel previous watchers
	w.mu.Lock()
	slog.Info("Cancel any previous DBConfigWatcher.")
	if w.cancelPrevious != nil {
		w.cancelPrevious()
	}

	// start new watcher
	slog.Info("Start new DBConfigWatcher.")
	wCtx, cancelPrevious := context.WithCancel(ctx)
	w.cancelPrevious = cancelPrevious
	w.mu.Unlock()
	w.watchInternal(wCtx, db, notifyCh)
}

func (w *DBConfigWatcher) watchInternal(ctx context.Context, db *database.DBMysql, notifyCh chan struct{}) {
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
