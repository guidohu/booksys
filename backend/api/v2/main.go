package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
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
	httpWebSetup                 = pflag.Bool("http_websetup", false, "If enabled the application can be setup via the /setup URI.")
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

func getFlags(v *viper.Viper) {
	v.BindPFlag("config", pflag.Lookup("config"))
	v.BindPFlag("http.port", pflag.Lookup("http_port"))
	v.BindPFlag("http.sessioninactivitytimeout", pflag.Lookup("http_session_inactivity_timeout"))
	v.BindPFlag("http.sessiontimeout", pflag.Lookup("http_session_timeout"))
	v.BindPFlag("http.uploadpath", pflag.Lookup("http_uploadpath"))
	v.BindPFlag("http.websetup", pflag.Lookup("http_websetup"))
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
	} else {
		slog.Info("Connected to database")
		c.SetDB(dbm)
		c.ReadConfigProperties()
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
	// Handlers without any authentication handling or middle ware.
	//-------------------------------------------------------------------
	mux.Handle("/api/v2/database/setup", h.WithFlagGuarded(h.SetupDBConfig))
	mux.Handle("/api/v2/health/status", http.HandlerFunc(h.HealthStatus))
	mux.Handle("/api/v2/ping", http.HandlerFunc(h.Ping))

	// Handlers that unauthenticated users can call
	//-------------------------------------------------------------------
	mux.Handle("/api/v2/auth/login", h.WithNoAuthentication(handlers.WithRequestBody(h.Login)))
	mux.Handle("/api/v2/auth/isloggedin", h.WithNoAuthentication(h.IsLoggedIn))
	mux.Handle("/api/v2/configuration/logo", h.WithNoAuthentication(h.GetLogoPath))
	mux.Handle("/api/v2/configuration/recaptcha-key", h.WithNoAuthentication(h.GetRecaptchaKey))
	mux.Handle("/api/v2/user/signup", h.WithNoAuthentication(handlers.WithRequestBody(h.SignUp, handlers.SignUpRequestValidationErrors)))
	mux.Handle("/api/v2/user/create-admin", h.WithFlagGuarded(h.WithNoAuthentication(handlers.WithRequestBody(h.MakeAdmin))))
	mux.Handle("/api/v2/user/password/reset-by-token", h.WithNoAuthentication(handlers.WithRequestBody(h.SetPasswordWithToken, handlers.SetPasswordWithTokenValidationErrors)))
	mux.Handle("/api/v2/user/password/token-request", h.WithNoAuthentication(handlers.WithRequestBody(h.GetPasswordResetToken, handlers.GetPasswordResetTokenValidationErrors)))

	// Handlers that users need to be authenticated and the role does not matter
	// or the handler handles the exact role requirements itself.
	//-------------------------------------------------------------------
	mux.Handle("/api/v2/auth/logout", h.WithAnyAuthentication(h.Logout))
	mux.Handle("/api/v2/auth/user", h.WithAnyAuthentication(h.User))
	mux.Handle("/api/v2/booking/day/list", h.WithAnyAuthentication(handlers.WithRequestBody(h.GetBookingDay)))
	mux.Handle("/api/v2/configuration/list", h.WithAnyAuthentication(h.GetPublicConfiguration))

	mux.Handle("/api/v2/user/my/balance", h.WithAnyAuthentication(h.GetMyBalance))
	mux.Handle("/api/v2/user/my/heats", h.WithAnyAuthentication(h.GetMyHeats))
	mux.Handle("/api/v2/user/my/heats/statistics", h.WithAnyAuthentication(h.GetMyHeatStats))
	mux.Handle("/api/v2/user/my/password/update", h.WithAnyAuthentication(h.UpdateMyPassword))
	mux.Handle("/api/v2/user/my/session/delete", h.WithAnyAuthentication(handlers.WithRequestBody(h.RemoveMyUserFromSession)))
	mux.Handle("/api/v2/user/my/sessions", h.WithAnyAuthentication(h.GetMySessions))
	mux.Handle("/api/v2/user/my/update", h.WithAnyAuthentication(h.UpdateMyUser))

	// Handlers that only admins have access to
	//-------------------------------------------------------------------
	mux.Handle("/api/v2/boat/engine-hour/entry/update", h.WithAdminAuthentication(h.UpdateEngineHoursEntry))
	mux.Handle("/api/v2/boat/engine-hour/update", h.WithAdminAuthentication(h.UpdateEngineHours))
	mux.Handle("/api/v2/boat/engine-hour/latest/get", h.WithAdminAuthentication(h.GetEngineHourLatest))
	mux.Handle("/api/v2/boat/engine-hours/list", h.WithAdminAuthentication(h.GetEngineHoursList))
	mux.Handle("/api/v2/boat/fuel-entries/get", h.WithAdminAuthentication(h.GetFuelEntries))
	mux.Handle("/api/v2/boat/fuel-entry/add", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddFuelEntry, handlers.FuelEntryValidationErrors)))
	mux.Handle("/api/v2/boat/fuel-entry/edit", h.WithAdminAuthentication(handlers.WithRequestBody(h.ChangeFuelEntry, handlers.FuelEntryValidationErrors)))
	mux.Handle("/api/v2/boat/fuel-entry/remove", h.WithAdminAuthentication(handlers.WithRequestBody(h.RemoveFuelEntry, handlers.FuelEntryValidationErrors)))
	mux.Handle("/api/v2/boat/maintenance-entries/get", h.WithAdminAuthentication(h.GetMaintenanceEntries))
	mux.Handle("/api/v2/boat/maintenance-entry/add", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddMaintenanceEntry, handlers.AddMaintenanceEntryValidationErrors)))
	mux.Handle("/api/v2/boat/mynautique/telemetry/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetBoatTelemetry)))
	// mux.Handle("/api/v2/database/config", http.HandlerFunc(h.WithAuthentication(h.GetDBConfig)))

	mux.Handle("/api/v2/accounting/expense_types/list", h.WithAdminAuthentication(h.GetExpenseTypes))
	mux.Handle("/api/v2/accounting/expense/add", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddExpense, handlers.AddTransactionValidationErrors)))
	mux.Handle("/api/v2/accounting/income_types/list", h.WithAdminAuthentication(h.GetIncomeTypes))
	mux.Handle("/api/v2/accounting/income/add", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddIncome, handlers.AddTransactionValidationErrors)))
	mux.Handle("/api/v2/accounting/statistics/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetAccountingStatistics, handlers.GetAccountingStatisticsValidationErrors)))
	mux.Handle("/api/v2/accounting/transactions/delete", h.WithAdminAuthentication(handlers.WithRequestBody(h.DeleteTransaction, handlers.DeleteTransactionValidationErrors)))
	mux.Handle("/api/v2/accounting/transactions/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetAccountingTransactions, handlers.GetAccountingTransactionsValidationErrors)))
	mux.Handle("/api/v2/accounting/years/list", h.WithAdminAuthentication(h.GetAccountingYears))

	mux.Handle("/api/v2/booking/series/list", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetBookingSeries)))

	mux.Handle("/api/v2/heat/change", h.WithAdminAuthentication(handlers.WithRequestBody(h.ChangeHeat)))
	mux.Handle("/api/v2/heat/delete", h.WithAdminAuthentication(handlers.WithRequestBody(h.DeleteHeat)))
	mux.Handle("/api/v2/heats/create", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddHeats)))

	mux.Handle("/api/v2/session/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetSession)))
	mux.Handle("/api/v2/session/create", h.WithAdminAuthentication(handlers.WithRequestBody(h.CreateSession, handlers.SessionValidationErrors)))
	mux.Handle("/api/v2/session/delete", h.WithAdminAuthentication(handlers.WithRequestBody(h.DeleteSession, handlers.DeleteSessionValidationErrors)))
	mux.Handle("/api/v2/session/edit", h.WithAdminAuthentication(handlers.WithRequestBody(h.EditSession, handlers.SessionValidationErrors)))
	mux.Handle("/api/v2/session/heats/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetSessionHeats)))
	mux.Handle("/api/v2/session/metadata/get", h.WithAdminAuthentication(handlers.WithRequestBody(h.GetSessionMetadata)))
	mux.Handle("/api/v2/session/user/add", h.WithAdminAuthentication(handlers.WithRequestBody(h.AddUserToSession)))
	mux.Handle("/api/v2/session/user/remove", h.WithAdminAuthentication(handlers.WithRequestBody(h.RemoveUserFromSession)))

	mux.Handle("/api/v2/mynautique/credentials/setup", h.WithAdminAuthentication(handlers.WithRequestBody(h.SetupMyNautiqueCredentials)))

	mux.Handle("/api/v2/user/delete", h.WithAdminAuthentication(handlers.WithRequestBody(h.DeleteUser, handlers.DeleteUserValidationErrors)))
	mux.Handle("/api/v2/user/list-detailed", h.WithAdminAuthentication(h.GetAllUsersDetailed))
	mux.Handle("/api/v2/user/list-short", h.WithAdminAuthentication(h.GetAllUsersShort))
	mux.Handle("/api/v2/user/lock/set", h.WithAdminAuthentication(handlers.WithRequestBody(h.SetUserLock, handlers.SetUserLockValidationErrors)))
	mux.Handle("/api/v2/user/group/create", h.WithAdminAuthentication(handlers.WithRequestBody(h.CreateUserGroup, handlers.UserGroupsValidationErrors)))
	mux.Handle("/api/v2/user/group/edit", h.WithAdminAuthentication(handlers.WithRequestBody(h.ChangeUserGroup, handlers.UserGroupsValidationErrors)))
	mux.Handle("/api/v2/user/group/delete", h.WithAdminAuthentication(handlers.WithRequestBody(h.DeleteUserGroup, handlers.UserGroupsValidationErrors)))
	mux.Handle("/api/v2/user/group/set", h.WithAdminAuthentication(handlers.WithRequestBody(h.SetUserGroup, handlers.SetUserGroupsValidationErrors)))
	mux.Handle("/api/v2/user/groups/get", h.WithAdminAuthentication(h.GetUserGroups))
	mux.Handle("/api/v2/user/roles/get", h.WithAdminAuthentication(h.GetUserRoles))

	mux.Handle("/api/v2/admin/configuration/list", h.WithAdminAuthentication(h.GetConfiguration))
	mux.Handle("/api/v2/admin/configuration/set", h.WithAdminAuthentication(handlers.WithRequestBody(h.SetConfiguration, handlers.ConfigurationMessageValidationErrors)))
	mux.Handle("/api/v2/admin/logs", h.WithAdminAuthentication(h.GetLogs))
	mux.Handle("/api/v2/admin/upload/logo", h.WithAdminAuthentication(h.UploadLogoFile))

	// Register pprof handlers manually
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register expvar handler manually
	mux.Handle("/debug/vars", expvar.Handler())
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
				func() {
					db, done := dbm.GetHandler()
					defer done()
					reconnect := false
					if db == nil || db.Ping() != nil {
						slog.Warn("Connection to database lost. Try to reconnect.")
						reconnect = true
					} else {
						initDone, err := db.IsInitialized()
						if err != nil || !initDone {
							slog.Warn("Uninitialized database. Try to reconnect and initialize.")
							reconnect = true
						}
					}
					if !reconnect {
						return
					}
					slog.Info("Initiate connection attempt to db.")
					err = dbm.ConnectAndReplace(nil)
					if err != nil {
						slog.Warn("Cannot reconnect database", slog.String("error", err.Error()))
						return
					}
					slog.Info("Reconnected database.")
					conf.ReadConfigProperties()
				}()
			case <-chConfigFileUpdate:
				func() {
					slog.Info("Reconnect database after config change.")
					dbSettings := getDBSettings(conf)
					err = dbm.ConnectAndReplace(&dbSettings)
					if err != nil {
						slog.Warn("Cannot connect with new database settings after config file update", slog.String("error", err.Error()))
						return
					}
					conf.ReadConfigProperties()
				}()
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
		slog.Info(fmt.Sprintf("Server listening on address %s", server.Addr))
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
