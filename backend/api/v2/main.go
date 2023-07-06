package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/database"
	"server/handlers"

	"golang.org/x/exp/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func parseFlags() {
	var port int
	var configFile string
	flag.IntVar(&port, "port", 80, "the port to listen on")
	flag.StringVar(&configFile, "configfile", "config.yaml", "the configuration file to be used")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.BindPFlag("port", pflag.Lookup("port"))
	viper.BindPFlag("configfile", pflag.Lookup(("configfile")))
	viper.Debug()

	flag.Parse()

	viper.Set("configfile", configFile)
	slog.Info("Config file set to", slog.String("configfile", configFile))
	viper.Set("port", port)
}

func readConfigFile() {
	configFile := viper.GetString("configfile")
	err := config.ReadConfig(configFile)
	if err != nil {
		slog.Error("Cannot read config:", slog.String("configFile", configFile), slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func connectDatabase() *database.DBMysql {
	if !viper.IsSet("database.dbname") {
		slog.Warn("Database settings are not present in configuration")
		return nil
	}

	db := &database.DBMysql{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Protocol: viper.GetString("database.protocol"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		DBName:   viper.GetString("database.dbname"),
	}
	if err := db.Connect(); err != nil {
		slog.Warn(fmt.Sprintf("Database is not properly setup or not reachable. Error returned from Connet(): %s", err))
		return nil
	}
	slog.Info("Connected to database", slog.String("name", viper.GetString("database.dbname")))
	return db
}

func main() {
	parseFlags()
	readConfigFile()

	// Connect to the database if configured.
	db := connectDatabase()
	defer db.Disconnect()
	h := handlers.NewHandler(db)

	// Watch config chages and create a new DB connection
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Configuration file changed", slog.String("file", e.Name))
		// Reconnect database upon config change
		newDB := connectDatabase()
		if newDB != nil {
			h.SetDB(newDB)
			fmt.Println("New db connected")
		}
		readConfigFile()
	})
	viper.WatchConfig()

	// File server to serve uploaded files
	fs := http.FileServer(http.Dir(viper.GetString("upload.path")))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Register all handlers
	http.Handle("/api/v2/health/status", http.HandlerFunc(h.HealthStatus))

	http.Handle("/api/v2/ping", http.HandlerFunc(h.Ping))
	http.Handle("/api/v2/auth/login", http.HandlerFunc(h.Login))
	http.Handle("/api/v2/auth/isloggedin", http.HandlerFunc(h.IsLoggedIn))
	http.Handle("/api/v2/auth/logout", http.HandlerFunc(h.WithAuthentication(h.Logout)))
	http.Handle("/api/v2/auth/user", http.HandlerFunc(h.WithAuthentication(h.User)))

	http.Handle("/api/v2/configuration/list", http.HandlerFunc(h.WithAuthentication(h.GetConfiguration)))
	http.Handle("/api/v2/configuration/logo", http.HandlerFunc(h.GetLogoPath))

	http.Handle("/api/v2/database/config", http.HandlerFunc(h.WithAuthentication(h.GetDBConfig)))
	http.Handle("/api/v2/database/setup", http.HandlerFunc(h.SetupDBConfig))

	http.Handle("/api/v2/booking/day/list", http.HandlerFunc(h.WithAuthentication(h.GetBookingDay)))
	// http.Handle("/api/v2/booking/month/list", http.HandlerFunc(h.WithAuthentication(h.GetBookingDay)))

	http.Handle("/api/v2/mynautique/credentials/setup", http.HandlerFunc(h.SetupMyNautiqueCredentials))

	http.Handle("/api/v2/user/signup", http.HandlerFunc(h.SignUp))
	http.Handle("/api/v2/user/create-admin", http.HandlerFunc(h.MakeAdmin))
	http.Handle("/api/v2/user/my/balance", http.HandlerFunc(h.WithAuthentication(h.GetMyBalance)))
	http.Handle("/api/v2/user/my/heats", http.HandlerFunc(h.WithAuthentication(h.GetMyHeats)))
	http.Handle("/api/v2/user/my/heats/statistics", http.HandlerFunc(h.WithAuthentication(h.GetMyHeatStats)))
	http.Handle("/api/v2/user/my/password/update", http.HandlerFunc(h.WithAuthentication(h.UpdateMyPassword)))
	http.Handle("/api/v2/user/my/sessions", http.HandlerFunc(h.WithAuthentication(h.GetMySessions)))
	http.Handle("/api/v2/user/my/update", http.HandlerFunc(h.WithAuthentication(h.UpdateMyUser)))

	http.Handle("/api/v2/admin/configuration/set", http.HandlerFunc(h.WithAuthentication(h.SetConfiguration)))
	http.Handle("/api/v2/admin/logs", http.HandlerFunc(h.WithAuthentication(h.GetLogs)))
	http.Handle("/api/v2/admin/upload/logo", http.HandlerFunc(h.WithAuthentication(h.UploadLogoFile)))
	// http.Handle("/api/v2/admin/configuration/list", http.HandlerFunc(h.WithAuthentication(h.GetLogs)))

	// TODO remove
	jss := http.FileServer(http.Dir("../../../frontend/dist/"))
	http.Handle("/", jss)

	slog.Info(fmt.Sprintf("Server listening on port %d\n", viper.GetUint16("port")))

	slog.Warn(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetUint16("port")), nil).Error())
}
