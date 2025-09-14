package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"server/config"
	"server/database"
	"server/handlers"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

// Flag definition
// Note: Configurations follow the following precedence from highest to lowest:
// - Flags
// - Environment Variables
// - Config File
var (
	// Place where we load the config from.
	configFile = pflag.String("config", "", "The configuration file to use.")
	// HTTP server settings
	httpPort                     = pflag.Int("http_port", 0, "The port the HTTP server listens on.")
	httpSessionInactivityTimeout = pflag.Uint("http_session_inactivity_timeout", 0, "Time until a HTTP session with no activity will be cancelled and a user gets logged out.")
	httpSessionTimeout           = pflag.Uint("http_session_timeout", 0, "Time in seconds until a user is logged out.")
	httpUploadPath               = pflag.String("http_uploadpath", "", "Path where content is uploaded to.")
	// Database settings
	databaseName     = pflag.String("database_dbname", "", "The database name.")
	databaseHost     = pflag.String("database_host", "", "The IP/hostname of the host the DB is on.")
	databasePassword = pflag.String("database_password", "", "The password for the DB user.")
	databasePort     = pflag.String("database_port", "", "The port for the DB connection.")
	databaseProtocol = pflag.String("database_protocol", "", "The protocol for the DB connection.")
	databaseUser     = pflag.String("database_user", "", "The user for the DB connection.")
	// MyNautique settings
	myNautiqueAPIKey = pflag.String("mynautique_api_key", "", "The API key for the mynautique integration.")
	// Additional control flags
	printConfig = pflag.Bool("print_config", false, "Prints the config an exits.")
)

var (
	mu sync.Mutex
)

func getFlags(v *viper.Viper) {
	v.BindPFlag("config", pflag.Lookup("config"))
	v.BindPFlag("http.port", pflag.Lookup("http_port"))
	v.BindPFlag("http.sessioninactivitytimeout", pflag.Lookup("http_session_inactivity_timeout"))
	v.BindPFlag("http.sessiontimeout", pflag.Lookup("http_session_timeout"))
	v.BindPFlag("http.uploadpath", pflag.Lookup("http_uploadpath"))
	v.BindPFlag("database.user", pflag.Lookup("database_user"))
	v.BindPFlag("database.password", pflag.Lookup("database_password"))
	v.BindPFlag("database.protocol", pflag.Lookup("database_protocol"))
	v.BindPFlag("database.host", pflag.Lookup("database_host"))
	v.BindPFlag("database.port", pflag.Lookup("database_port"))
	v.BindPFlag("database.dbname", pflag.Lookup("database_dbname"))
	v.BindPFlag("mynautique.api.key", pflag.Lookup("mynautique_api_key"))
	pflag.Parse()
}

func getDBSettings(c *config.Config) database.Settings {
	user, _ := c.GetString("database.user")
	password, _ := c.GetString("database.password")
	protocol, _ := c.GetString("database.protocol")
	host, _ := c.GetString("database.host")
	port, _ := c.GetString("database.port")
	dbname, _ := c.GetString("database.dbname")
	return database.Settings{
		User:     user,
		Password: password,
		Protocol: protocol,
		Host:     host,
		Port:     port,
		DBName:   dbname,
	}
}

func maybeConnectedDatabaseManager(c *config.Config) *database.Manager {
	settings := getDBSettings(c)
	dbm := database.NewManager(settings)
	if err := dbm.Connect(nil); err != nil {
		slog.Warn("Could not connect to database", slog.String("error", err.Error()))
	}
	return dbm
}

func setupLogger() {
	lh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(lh)
	slog.SetDefault(logger)
}

func registerUploadsServer(mux *http.ServeMux, c *config.Config) {
	// File server to serve uploaded files
	uploadPath, _ := c.GetString("http.uploadpath")
	if uploadPath == "" {
		slog.Warn("http.uploadpath is not set, not serving uploaded files")
		return
	}
	if !filepath.IsAbs(uploadPath) {
		slog.Warn("http.uploadpath should be an absolute path")
	}
	if strings.Contains(uploadPath, "..") {
		slog.Warn("http.uploadpath cannot not contain '..'")
		os.Exit(1)
	}
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		slog.Warn("http.uploadpath cannot be converted to an absolute path")
		os.Exit(1)
	}
	fs := http.FileServer(http.Dir(absUploadPath))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
	slog.Info("Serving uploaded files from", slog.String("path", absUploadPath))
}

func registerHandlers(mux *http.ServeMux, h *handlers.Handler) {
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
	mux.Handle("/api/v2/boat/fuel-entry/remove", http.HandlerFunc(h.WithAuthentication(h.RemoveFuelEntry)))
	mux.Handle("/api/v2/boat/maintenance-entries/get", http.HandlerFunc(h.WithAuthentication(h.GetMaintenanceEntries)))
	mux.Handle("/api/v2/boat/maintenance-entry/add", http.HandlerFunc(h.WithAuthentication(h.AddMaintenanceEntry)))
	mux.Handle("/api/v2/boat/mynautique/telemetry/get", http.HandlerFunc(h.WithAuthentication(h.GetBoatTelemetry)))

	mux.Handle("/api/v2/configuration/list", http.HandlerFunc(h.WithAuthentication(h.GetPublicConfiguration)))
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

	mux.Handle("/api/v2/admin/configuration/list", http.HandlerFunc(h.WithAuthentication(h.GetConfiguration)))
	mux.Handle("/api/v2/admin/configuration/set", http.HandlerFunc(h.WithAuthentication(h.SetConfiguration)))
	mux.Handle("/api/v2/admin/logs", http.HandlerFunc(h.WithAuthentication(h.GetLogs)))
	mux.Handle("/api/v2/admin/upload/logo", http.HandlerFunc(h.WithAuthentication(h.UploadLogoFile)))
}

func main() {
	setupLogger()

	// Get startup configuration.
	v := viper.New()
	getFlags(v)
	conf, err := config.NewConfig(v)
	if err != nil {
		slog.Error("Invalid config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// If requested, print config and exit.
	if *printConfig {
		fmt.Println(conf.ToStringFull())
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	dbm := maybeConnectedDatabaseManager(conf)
	hp := handlers.HandlerParams{
		Database:      dbm,
		Configuration: conf,
	}
	h := handlers.NewHandler(hp)
	conf.SetDB(dbm)

	// We want to know about config file updates, because the
	// database configuration in there might change.
	chConfigFileUpdate := conf.WatchFile()

	// Watch for changes in the properties.
	// TODO: If we really do not provide a channel, we do not
	// really need to poll the configuration?
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
				db, done := dbm.GetHandler()
				if db == nil || db.Ping() != nil {
					slog.Info("Schedule connection attempt to db.")
					err = dbm.ConnectAndReplace(nil)
					if err != nil {
						slog.Warn("Cannot reconnect database", slog.String("error", err.Error()))
					} else {
						slog.Info("Reconnected database.")
					}
					// TODO update conf with latest DB values
				}
				done()
			case <-chConfigFileUpdate:
				slog.Info("Reconnect database after config change.")
				dbSettings := getDBSettings(conf)
				err = dbm.ConnectAndReplace(&dbSettings)
				if err != nil {
					slog.Warn("Cannot connect with new database settings after config file update", slog.String("error", err.Error()))
				}
				// TODO inform conf to load data from database in case it caches config
			}
		}
	}()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", v.GetUint16("http.port")),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	registerUploadsServer(mux, conf)
	registerHandlers(mux, h)

	// TODO remove
	// jss := http.FileServer(http.Dir("../../../frontend/dist/"))
	// mux.Handle("/", jss)

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
