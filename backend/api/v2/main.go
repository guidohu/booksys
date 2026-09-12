package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"server/config"
	"server/database"
	"server/handlers"
	"server/version"
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
	httpSecureCookie             = pflag.Bool("http_secure_cookie", true, "If enabled the session cookie is only sent over HTTPS. Only disable this for local development over plain HTTP.")
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
	// Debug settings
	debugPort = pflag.Int("debug_port", 0, "The port the debug listener (pprof, expvar) serves on. 0 disables it.")
	// Display settings
	environment = pflag.String("environment", "", "The environment this app is running as.")
	// Additional control flags
	printConfig       = pflag.Bool("print_config", false, "Prints the config with the secrets redacted and exits.")
	printConfigUnsafe = pflag.Bool("print_config_unsafe", false, "Prints the config with the secrets in plaintext and exits. Only for debugging a credential itself, the output must not be shared.")
	printVersion      = pflag.Bool("version", false, "Prints the version an exits.")
)

// getFlags binds the command line flags to their configuration keys and parses
// them. A binding failure means a flag name and its key have drifted apart,
// which is a programming error, so it is reported rather than ignored.
func getFlags(v *viper.Viper) error {
	bindings := map[string]string{
		"config":                        "config",
		"debug.port":                    "debug_port",
		"http.port":                     "http_port",
		"http.securecookie":             "http_secure_cookie",
		"http.sessioninactivitytimeout": "http_session_inactivity_timeout",
		"http.sessiontimeout":           "http_session_timeout",
		"http.uploadpath":               "http_uploadpath",
		"http.websetup":                 "http_websetup",
		"database.user":                 "database_user",
		"database.password":             "database_password",
		"database.protocol":             "database_protocol",
		"database.host":                 "database_host",
		"database.port":                 "database_port",
		"database.dbname":               "database_dbname",
		"mynautique.api.key":            "mynautique_api_key",
		"environment":                   "environment",
	}

	for viperKey, flagName := range bindings {
		flag := pflag.Lookup(flagName)
		if flag == nil {
			return fmt.Errorf("flag %q not found", flagName)
		}
		if err := v.BindPFlag(viperKey, flag); err != nil {
			return fmt.Errorf("cannot bind flag %q to %q: %w", flagName, viperKey, err)
		}
	}

	pflag.Parse()
	return nil
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
		slog.Warn("Could not connect to database", slog.Any("error", err))
	} else {
		slog.Info("Connected to database")
		c.SetDB(dbm)
		if err := c.ReadConfigProperties(); err != nil {
			slog.Warn("Cannot read config properties", slog.Any("error", err))
		}
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
	mux.Handle("/api/v2/admin/configuration/set", h.WithAdminAuthentication(h.WithConfigContext(handlers.WithRequestBody(h.SetConfiguration, handlers.ConfigurationMessageValidationErrors))))
	mux.Handle("/api/v2/admin/logs", h.WithAdminAuthentication(h.GetLogs))
	mux.Handle("/api/v2/admin/upload/logo", h.WithAdminAuthentication(h.UploadLogoFile))
}

// registerDebugHandlers registers the profiling and metrics endpoints. They are
// deliberately kept off the API mux and served on a port of their own: a heap
// profile contains session secrets and the database password, and
// /debug/pprof/profile pins a CPU profiler for 30 seconds per request. Keeping
// them on a separate port makes reachability a deployment decision instead of
// something a reverse proxy has to remember to filter out.
func registerDebugHandlers(mux *http.ServeMux) {
	// Register pprof handlers manually
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register expvar handler manually
	mux.Handle("/debug/vars", expvar.Handler())
}

// newDebugServer returns the debug listener for the configured debug.port, or
// nil if no port is set. Note that it has no WriteTimeout: a CPU profile runs
// for 30 seconds by default, so any write deadline shorter than that would
// truncate the very profile it was asked for.
func newDebugServer(c *config.Config) (*http.Server, error) {
	// Read the raw value instead of going through GetInt64, which discards the
	// parse error and so cannot tell "not set" apart from "not a number". A
	// typo in BOOKSYS_DEBUG_PORT would otherwise disable the listener silently,
	// and that only becomes apparent when someone needs to profile a live
	// incident.
	value, _ := c.GetString("debug.port")
	if value == "" {
		return nil, nil
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("debug.port %q is not a number: %w", value, err)
	}
	if port == 0 {
		return nil, nil
	}
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf("debug.port %d is not a valid port", port)
	}
	mux := http.NewServeMux()
	registerDebugHandlers(mux)
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}, nil
}

func printVersionAndExit() {
	version.PrintVersion()
	os.Exit(0)
}

func main() {
	setupLogger()

	v := viper.New()
	if err := getFlags(v); err != nil {
		slog.Error("Cannot process the command line flags", slog.Any("error", err))
		os.Exit(1)
	}

	// If requested, print version and exit.
	if *printVersion {
		printVersionAndExit()
	}

	// Get startup configuration.
	conf, err := config.NewConfig(v)
	if err != nil {
		slog.Error("Invalid config", slog.Any("error", err))
		os.Exit(1)
	}

	// If requested, print config and exit. The plaintext dump has to be asked
	// for explicitly, so that the common case of sharing a config dump does
	// not hand out the database password along with it.
	if *printConfigUnsafe {
		fmt.Fprintln(os.Stderr, "WARNING: --print_config_unsafe prints secrets in plaintext, do not share this output.")
		fmt.Println(conf.ToStringFullUnsafe())
		os.Exit(0)
	}
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
					if err := dbm.ConnectAndReplace(nil); err != nil {
						slog.Warn("Cannot reconnect database", slog.Any("error", err))
						return
					}
					slog.Info("Reconnected database.")
					if err := conf.ReadConfigProperties(); err != nil {
						slog.Warn("Cannot read config properties after reconnect", slog.Any("error", err))
					}
				}()
			case <-chConfigFileUpdate:
				func() {
					slog.Info("Reconnect database after config change.")
					dbSettings := getDBSettings(conf)
					if err := dbm.ConnectAndReplace(&dbSettings); err != nil {
						slog.Warn("Cannot connect with new database settings after config file update", slog.Any("error", err))
						return
					}
					if err := conf.ReadConfigProperties(); err != nil {
						slog.Warn("Cannot read config properties after config file update", slog.Any("error", err))
					}
				}()
			}
		}
	}()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", v.GetUint16("http.port")),
		// Cap the request body once, around the whole mux, so that every route
		// is covered including the ones registered without other middleware.
		Handler:      handlers.WithMaxBodySize(mux, handlers.MaxRequestBodyBytes),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	registerUploadsServer(mux, conf)
	registerHandlers(mux, h)

	debugServer, err := newDebugServer(conf)
	if err != nil {
		slog.Error("Cannot set up the debug listener", slog.Any("error", err))
		os.Exit(1)
	}
	if debugServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Warn("Debug listener enabled, this port must not be reachable publicly", slog.String("address", debugServer.Addr))
			if err := debugServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Debug server failure", slog.Any("error", err))
			}
			slog.Info("Debug listener stopped.")
		}()
	}

	// TODO remove
	// jss := http.FileServer(http.Dir("../../../frontend/dist/"))
	// mux.Handle("/", jss)

	wg.Add(1)
	go func() {
		slog.Info("Server listening", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failure", slog.Any("error", err))
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
	if debugServer != nil {
		if err := debugServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Debug server shutdown failure", slog.Any("error", err))
		}
	}
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown failure", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("HTTP server stopped.")

	// Cancel context and wait for all routines to end.
	cancel()
	wg.Wait()
	slog.Info("Shutdown complete.")
}
