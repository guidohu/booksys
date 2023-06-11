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

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var port int
	var configFile string
	flag.IntVar(&port, "port", 80, "the port to listen on")
	flag.StringVar(&configFile, "config_file", "config.yaml", "the configuration file to be used")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.BindPFlag("port", pflag.Lookup("port"))
	viper.BindPFlag("configfile", pflag.Lookup(("config_file")))
	viper.Debug()

	flag.Parse()

	// Read the configuration file
	viper.Set("config_file", configFile)
	slog.Info("Config file set to", slog.String("config_file", configFile))
	viper.Set("port", port)
	err := config.ReadConfig(configFile)
	if err != nil {
		slog.Error("Cannot read config:", slog.String("configFile", configFile), slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Server port found", slog.Int("port", port))

	// Connect to the database if configured.
	if !viper.IsSet("database.dbname") {
		slog.Warn("Database configuration is not present in configuration file", slog.String("configFile", configFile))
	}
	db := &database.DBMysql{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Protocol: viper.GetString("database.protocol"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		DBName:   viper.GetString("database.dbname"),
	}
	err = db.Connect()
	if err != nil {
		slog.Warn(fmt.Sprintf("Database is not properly setup or not reachable. Error returned from Connet(): %s", err))
	}

	// initialize handlers
	h := handlers.NewHandler(db)
	a := handlers.NewIdenticationMiddleware(db)

	http.Handle("/api/v2/health/status", http.HandlerFunc(h.HealthStatus))

	http.Handle("/api/v2/ping", http.HandlerFunc(h.Ping))
	http.Handle("/api/v2/auth/login", http.HandlerFunc(h.Login))
	http.Handle("/api/v2/auth/isloggedin", http.HandlerFunc(h.IsLoggedIn))
	http.Handle("/api/v2/auth/logout", http.HandlerFunc(a.Authenticated(h.Logout)))
	http.Handle("/api/v2/auth/user", http.HandlerFunc(a.Authenticated(h.User)))

	slog.Info(fmt.Sprintf("Server listening on %d\n", port))

	slog.Warn(http.ListenAndServe(fmt.Sprintf(":%d", port), nil).Error())
}
