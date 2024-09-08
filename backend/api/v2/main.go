package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/config"
	"server/database"
	"server/handlers"
	"sync"
	"syscall"
	"time"

	"golang.org/x/exp/slog"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Flag definition
// Note: Configuration can be either specified in the following precedence:
// - Flags
// - Environment Variables
// - Config File
var configFile *string = flag.String("config", "", "The configuration file to use.")

var httpPort *int = flag.Int("http_port", 0, "The port the HTTP server listens on.")
var httpSessionInactivityTimeout *uint = flag.Uint("http_session_inactivity_timeout", 0, "Time until a HTTP session with no activity will be cancelled and a user gets logged out.")
var httpSessionTimeout *uint = flag.Uint("http_session_timeout", 0, "Time in seconds until a user is logged out.")
var httpUploadPath *string = flag.String("http_upload_path", "", "Path where content is uploaded to.")

var databaseName *string = flag.String("database_name", "", "The database name.")
var databaseHost *string = flag.String("database_host", "", "The IP/hostname of the host the DB is on.")
var databasePassword *string = flag.String("database_password", "", "The password for the DB user.")
var databasePort *string = flag.String("database_port", "", "The port for the DB connection.")
var databaseProtocol *string = flag.String("database_protocol", "", "The protocol for the DB connection.")
var databaseUser *string = flag.String("database_user", "", "The user for the DB connection.")

var myNautiqueAPIKey *string = flag.String("mynautique_api_key", "", "The API key for the mynautique integration.")

var printConfig *bool = flag.Bool("print_config", false, "Prints the config an exits.")

func getFlags(v *viper.Viper) {
	v.BindPFlag("config", flag.Lookup("config"))
	v.BindPFlag("http.port", flag.Lookup("http_port"))
	v.BindPFlag("http.sessioninactivitytimeout", flag.Lookup("http_session_inactivity_timeout"))
	v.BindPFlag("http.sessiontimeout", flag.Lookup("http_session_timeout"))
	v.BindPFlag("http.uploadpath", flag.Lookup("http_upload_path"))
	v.BindPFlag("database.user", flag.Lookup("database_user"))
	v.BindPFlag("database.password", flag.Lookup("database_password"))
	v.BindPFlag("database.protocol", flag.Lookup("database_protocol"))
	v.BindPFlag("database.host", flag.Lookup("database_host"))
	v.BindPFlag("database.port", flag.Lookup("database_port"))
	v.BindPFlag("database.dbname", flag.Lookup("database_name"))
	v.BindPFlag("mynautique.api.key", flag.Lookup("mynautique_api_key"))
	flag.Parse()
}

// func setDefaults(v *viper.Viper) {
// 	v.SetDefault("config", "")
// 	v.SetDefault("http.port", "80")
// 	v.SetDefault("http.sessioninactivitytimeout", "604800")
// 	v.SetDefault("http.sessiontimeout", "31536000")
// 	v.SetDefault("http.uploadpath", "./uploads")
// 	v.SetDefault("database.user", "")
// 	v.SetDefault("database.password", "")
// 	v.SetDefault("database.protocol", "tcp")
// 	v.SetDefault("database.host", "127.0.0.1")
// 	v.SetDefault("database.port", "3306")
// 	// v.SetDefault("database.dbname", "")
// 	v.SetDefault("mynautique.api.key", "")
// }

// func getEnvironment(v *viper.Viper) {
// 	v.SetEnvPrefix("BOOKSYS")
// 	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
// 	v.AutomaticEnv()
// }

// func printConfigString(v *viper.Viper) {
// 	conf := &config.Configuration{}
// 	err := v.Unmarshal(conf)
// 	if err != nil {
// 		slog.Error("Invalid configuration. Cannot decode configuration.", slog.String("error", err.Error()))
// 		os.Exit(1)
// 	}
// 	c, err := yaml.Marshal(conf)
// 	if err != nil {
// 		slog.Error("Cannot generate configuration", slog.String("error", err.Error()))
// 		os.Exit(1)
// 	}
// 	fmt.Print(string(c))
// }

func readConfigFile(v *viper.Viper) {
	// If we do not have a config file we return.
	if v.GetString("config") == "" {
		slog.Info("No config file provided.")
		return
	}
	slog.Info("Read config from", slog.String("config", v.GetString("config")))
	_, err := os.Stat(v.GetString("config"))
	if err != nil {
		slog.Error("File does not exist", slog.String("config", v.GetString("config")), slog.String("error", err.Error()))
	}
	v.SetConfigFile(v.GetString("config"))
	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("No config file found", slog.String("config", v.GetString("config")))
			os.Exit(1)
		}
	}
	// Verify that the configuration matches the Configuration struct.
	conf := &config.Configuration{}
	err = v.Unmarshal(conf)
	if err != nil {
		slog.Error("Invalid configuration. Cannot decode configuration.", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func connectDatabase(c *config.Config) *database.DBMysql {
	user, _ := c.GetString("database.user")
	password, _ := c.GetString("database.password")
	protocol, _ := c.GetString("database.protocol")
	host, _ := c.GetString("database.host")
	port, _ := c.GetString("database.port")
	dbname, _ := c.GetString("database.dbname")
	db := &database.DBMysql{
		User:     user,
		Password: password,
		Protocol: protocol,
		Host:     host,
		Port:     port,
		DBName:   dbname,
	}
	slog.Info("Connecting database client to:", slog.String("address", db.String()))
	if err := db.Connect(); err != nil {
		slog.Warn(fmt.Sprintf("Database is not properly setup or not reachable. Error returned from Connect(): %s", err))
		return nil
	}
	slog.Info("Connected to database", slog.String("name", dbname))
	return db
}

func reloadDBConfigWatcher(ctx context.Context, wg *sync.WaitGroup, w *config.DBConfigWatcher, db *database.DBMysql, ch chan struct{}) {
	// TODO fix racecondition in db
	wg.Add(1)
	go func() {
		slog.Info("Start watching DB config.")
		w.Watch(ctx, db, ch)
		slog.Info("Done watching DB config.")
		wg.Done()
	}()
}

func main() {
	lh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(lh)
	slog.SetDefault(logger)

	// Get startup configuration.
	v := viper.New()
	getFlags(v)
	conf, err := config.NewConfig(v)
	if err != nil {
		slog.Error("Invalid config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	uploadPath, _ := conf.GetString("http.uploadpath")

	// If requested, print config and exit.
	if *printConfig {
		fmt.Println(conf.ToStringFull())
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	db := connectDatabase(conf)
	hp := handlers.HandlerParams{
		Database:      db,
		Configuration: conf,
	}
	h := handlers.NewHandler(hp)

	if db != nil {
		conf.SetDB(db)
	}

	// We want to know about config file updates, because the
	// database configuration might change.
	chConfigFileUpdate := conf.WatchFile()
	// chPropertyChange := make(chan struct{})

	// Watch for changes in the properties.
	wg.Add(1)
	go func() {
		slog.Info("Watch properties routine started.")
		conf.WatchProperties(ctx, nil /*chPropertyChange*/)
		slog.Info("Watch properties routine done.")
		wg.Done()
	}()

	// Auto(re)connect database, e.g. if database is not up
	// when started, or configuration is not present yet.
	wg.Add(1)
	go func() {
		slog.Info("Auto reconnect routine started.")
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				wg.Done()
				slog.Info("Auto reconnect routine done.")
				return
			case <-ticker.C:
				if db == nil || db.Ping() != nil {
					slog.Info("Schedule connection attempt to db.")
					db = connectDatabase(conf)
					h.SetDB(db)
					conf.SetDB(db)
				}
			case <-chConfigFileUpdate:
				slog.Info("Reconnect database after config change.")
				db = connectDatabase(conf)
				h.SetDB(db)
				conf.SetDB(db)
				conf.ReadConfigProperties()
				uploadPath, _ = conf.GetString("http.uploadpath")
			}
		}
	}()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", v.GetUint16("http.port")),
		Handler: mux,
	}

	// File server to serve uploaded files
	fs := http.FileServer(http.Dir(uploadPath))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Register all handlers
	mux.Handle("/api/v2/health/status", http.HandlerFunc(h.HealthStatus))

	mux.Handle("/api/v2/ping", http.HandlerFunc(h.Ping))
	mux.Handle("/api/v2/auth/login", http.HandlerFunc(h.Login))
	mux.Handle("/api/v2/auth/isloggedin", http.HandlerFunc(h.IsLoggedIn))
	mux.Handle("/api/v2/auth/logout", http.HandlerFunc(h.WithAuthentication(h.Logout)))
	mux.Handle("/api/v2/auth/user", http.HandlerFunc(h.WithAuthentication(h.User)))

	mux.Handle("/api/v2/boat/engine-hour/entry/update", http.HandlerFunc(h.WithAuthentication(h.UpdateEngineHoursEntry)))
	mux.Handle("/api/v2/boat/engine-hour/update", http.HandlerFunc(h.WithAuthentication(h.UpdateEngineHours)))
	mux.Handle("/api/v2/boat/engine-hour/latest/get", http.HandlerFunc(h.WithAuthentication(h.GetEngineHourLatest)))
	mux.Handle("/api/v2/boat/engine-hours/list", http.HandlerFunc(h.WithAuthentication(h.GetEngineHoursList)))
	mux.Handle("/api/v2/boat/fuel-entries/get", http.HandlerFunc(h.WithAuthentication(h.GetFuelEntries)))
	mux.Handle("/api/v2/boat/fuel-entry/add", http.HandlerFunc(h.WithAuthentication(h.AddFuelEntry)))
	mux.Handle("/api/v2/boat/fuel-entry/edit", http.HandlerFunc(h.WithAuthentication(h.ChangeFuelEntry)))
	mux.Handle("/api/v2/boat/maintenance-entries/get", http.HandlerFunc(h.WithAuthentication(h.GetMaintenanceEntries)))
	mux.Handle("/api/v2/boat/maintenance-entry/add", http.HandlerFunc(h.WithAuthentication(h.AddMaintenanceEntry)))
	mux.Handle("/api/v2/boat/mynautique/telemetry/get", http.HandlerFunc(h.WithAuthentication(h.GetBoatTelemetry)))

	mux.Handle("/api/v2/configuration/list", http.HandlerFunc(h.WithAuthentication(h.GetConfiguration)))
	mux.Handle("/api/v2/configuration/logo", http.HandlerFunc(h.GetLogoPath))
	mux.Handle("/api/v2/configuration/recaptcha-key", http.HandlerFunc(h.GetRecaptchaKey))

	// mux.Handle("/api/v2/database/config", http.HandlerFunc(h.WithAuthentication(h.GetDBConfig)))
	mux.Handle("/api/v2/database/setup", http.HandlerFunc(h.SetupDBConfig))

	mux.Handle("/api/v2/accounting/expense_types/list", http.HandlerFunc(h.WithAuthentication(h.GetExpenseTypes)))
	mux.Handle("/api/v2/accounting/expense/add", http.HandlerFunc(h.WithAuthentication(h.AddExpense)))
	mux.Handle("/api/v2/accounting/income_types/list", http.HandlerFunc(h.WithAuthentication(h.GetIncomeTypes)))
	mux.Handle("/api/v2/accounting/income/add", http.HandlerFunc(h.WithAuthentication(h.AddIncome)))
	mux.Handle("/api/v2/accounting/statistics/get", http.HandlerFunc(h.WithAuthentication(h.GetAccountingStatistics)))
	mux.Handle("/api/v2/accounting/transactions/delete", http.HandlerFunc(h.WithAuthentication(h.DeleteTransaction)))
	mux.Handle("/api/v2/accounting/transactions/get", http.HandlerFunc(h.WithAuthentication(h.GetAccountingTransactions)))
	mux.Handle("/api/v2/accounting/years/list", http.HandlerFunc(h.WithAuthentication(h.GetAccountingYears)))

	mux.Handle("/api/v2/booking/day/list", http.HandlerFunc(h.WithAuthentication(h.GetBookingDay)))
	mux.Handle("/api/v2/booking/series/list", http.HandlerFunc(h.WithAuthentication(h.GetBookingSeries)))

	mux.Handle("/api/v2/heat/change", http.HandlerFunc(h.WithAuthentication(h.ChangeHeat)))
	mux.Handle("/api/v2/heat/delete", http.HandlerFunc(h.WithAuthentication(h.DeleteHeat)))
	mux.Handle("/api/v2/heats/create", http.HandlerFunc(h.WithAuthentication(h.AddHeats)))

	mux.Handle("/api/v2/session/get", http.HandlerFunc(h.WithAuthentication(h.GetSession)))
	mux.Handle("/api/v2/session/create", http.HandlerFunc(h.WithAuthentication(h.CreateSession)))
	mux.Handle("/api/v2/session/delete", http.HandlerFunc(h.WithAuthentication(h.DeleteSession)))
	mux.Handle("/api/v2/session/edit", http.HandlerFunc(h.WithAuthentication(h.EditSession)))
	mux.Handle("/api/v2/session/heats/get", http.HandlerFunc(h.WithAuthentication(h.GetSessionHeats)))
	mux.Handle("/api/v2/session/metadata/get", http.HandlerFunc(h.WithAuthentication(h.GetSessionMetadata)))
	mux.Handle("/api/v2/session/user/add", http.HandlerFunc(h.WithAuthentication(h.AddUserToSession)))
	mux.Handle("/api/v2/session/user/remove", http.HandlerFunc(h.WithAuthentication(h.RemoveUserFromSession)))

	mux.Handle("/api/v2/mynautique/credentials/setup", http.HandlerFunc(h.SetupMyNautiqueCredentials))
	mux.Handle("/api/v2/user/signup", http.HandlerFunc(h.SignUp))
	mux.Handle("/api/v2/user/create-admin", http.HandlerFunc(h.MakeAdmin))
	mux.Handle("/api/v2/user/delete", http.HandlerFunc(h.WithAuthentication(h.DeleteUser)))
	mux.Handle("/api/v2/user/list-detailed", http.HandlerFunc(h.WithAuthentication(h.GetAllUsersDetailed)))
	mux.Handle("/api/v2/user/list-short", http.HandlerFunc(h.WithAuthentication(h.GetAllUsersShort)))
	mux.Handle("/api/v2/user/lock/set", http.HandlerFunc(h.WithAuthentication(h.SetUserLock)))
	mux.Handle("/api/v2/user/group/create", http.HandlerFunc(h.WithAuthentication(h.CreateUserGroup)))
	mux.Handle("/api/v2/user/group/edit", http.HandlerFunc(h.WithAuthentication(h.ChangeUserGroup)))
	mux.Handle("/api/v2/user/group/delete", http.HandlerFunc(h.WithAuthentication(h.DeleteUserGroup)))
	mux.Handle("/api/v2/user/group/set", http.HandlerFunc(h.WithAuthentication(h.SetUserGroup)))
	mux.Handle("/api/v2/user/groups/get", http.HandlerFunc(h.WithAuthentication(h.GetUserGroups)))
	mux.Handle("/api/v2/user/my/balance", http.HandlerFunc(h.WithAuthentication(h.GetMyBalance)))
	mux.Handle("/api/v2/user/my/heats", http.HandlerFunc(h.WithAuthentication(h.GetMyHeats)))
	mux.Handle("/api/v2/user/my/heats/statistics", http.HandlerFunc(h.WithAuthentication(h.GetMyHeatStats)))
	mux.Handle("/api/v2/user/my/password/update", http.HandlerFunc(h.WithAuthentication(h.UpdateMyPassword)))
	mux.Handle("/api/v2/user/my/session/delete", http.HandlerFunc(h.WithAuthentication(h.RemoveMyUserFromSession)))
	mux.Handle("/api/v2/user/my/sessions", http.HandlerFunc(h.WithAuthentication(h.GetMySessions)))
	mux.Handle("/api/v2/user/my/update", http.HandlerFunc(h.WithAuthentication(h.UpdateMyUser)))
	mux.Handle("/api/v2/user/password/reset-by-token", http.HandlerFunc(h.SetPasswordWithToken))
	mux.Handle("/api/v2/user/password/token-request", http.HandlerFunc(h.GetPasswordResetToken))
	mux.Handle("/api/v2/user/roles/get", http.HandlerFunc(h.WithAuthentication(h.GetUserRoles)))

	mux.Handle("/api/v2/admin/configuration/set", http.HandlerFunc(h.WithAuthentication(h.SetConfiguration)))
	mux.Handle("/api/v2/admin/logs", http.HandlerFunc(h.WithAuthentication(h.GetLogs)))
	mux.Handle("/api/v2/admin/upload/logo", http.HandlerFunc(h.WithAuthentication(h.UploadLogoFile)))

	// TODO remove
	jss := http.FileServer(http.Dir("../../../frontend/dist/"))
	mux.Handle("/", jss)

	wg.Add(1)
	go func() {
		slog.Info(fmt.Sprintf("Server listening on address %s\n", server.Addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failure", slog.String("error", err.Error()))
			os.Exit(1)
		}
		slog.Info("Stopped serving new connections.")
		wg.Done()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Wait for last requests to get served before shutting down.
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown failure", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("HTTP server stopped.")

	// Cancel context and wait for all routines to end.
	cancel()
	wg.Wait()
	slog.Info("Shutdown complete.")
}
