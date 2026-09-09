package database

import (
	"log/slog"
	"sync"
)

// Manager owns the database connections. It hands out the most recent
// connection to callers and keeps superseded ones alive until their last
// client is done with them.
type Manager struct {
	generations    int64
	settings       Settings
	latestDatabase *Handle
	// Connections we want to close but there might still be clients
	// using those connections.
	previousDatabases []*Handle
	mu                sync.RWMutex
}

// Handle is a single database connection together with a wait group that
// tracks how many callers are still using it.
type Handle struct {
	id int64
	db Database
	wg sync.WaitGroup
}

// InstanceCreator creates the Database used by a Manager. It is a variable so
// that tests can substitute a fake implementation.
var InstanceCreator = func(settings Settings) Database {
	return NewDBMysql(settings)
}

// NewManager returns a Manager for the given settings. It does not connect,
// call Connect for that.
func NewManager(settings Settings) *Manager {
	return &Manager{
		settings: settings,
	}
}

// UpdateSettings replaces the settings used for the next connection attempt.
// It does not affect connections that are already established.
func (dbm *Manager) UpdateSettings(settings Settings) {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()
	dbm.settings = settings
}

// Connect establishes a new connection and makes it the one GetHandler returns.
// If settings is nil the settings the Manager was created with are used,
// otherwise the given settings are used and remembered.
func (dbm *Manager) Connect(settings *Settings) error {
	dbSettings := dbm.settings
	if settings != nil {
		dbSettings = *settings
	}
	newDB := InstanceCreator(dbSettings)
	err := newDB.Connect()
	if err != nil {
		return err
	}
	newDBHandle := &Handle{
		id: dbm.generations,
		db: newDB,
	}
	slog.Info("Successfully connected to database", slog.Int64("id", newDBHandle.id))
	dbm.mu.Lock()
	defer dbm.mu.Unlock()
	if dbm.latestDatabase != nil {
		dbm.previousDatabases = append(dbm.previousDatabases, dbm.latestDatabase)
	}
	dbm.latestDatabase = newDBHandle
	if settings != nil {
		dbm.settings = *settings
	}
	dbm.generations++
	return nil
}

// ConnectAndReplace connects like Connect and then closes the superseded
// connections in the background, once their last client is done with them.
func (dbm *Manager) ConnectAndReplace(settings *Settings) error {
	err := dbm.Connect(settings)
	if err != nil {
		return err
	}
	go func() {
		slog.Info("Closing database handles in background")
		// Work on a snapshot, the list can grow while we wait for the
		// clients of the old connections to finish.
		dbm.mu.RLock()
		pending := make([]*Handle, len(dbm.previousDatabases))
		copy(pending, dbm.previousDatabases)
		dbm.mu.RUnlock()

		for _, dbh := range pending {
			dbh.wg.Wait()
			slog.Info("Closing database handle", slog.Int64("id", dbh.id))
			dbm.mu.Lock()
			dbh.db.Disconnect()
			dbh.db = nil
			dbm.mu.Unlock()
		}
		dbm.mu.Lock()
		var remainingDatabases []*Handle
		for i, dbh := range dbm.previousDatabases {
			if dbh.db != nil {
				remainingDatabases = append(remainingDatabases, dbm.previousDatabases[i])
			}
		}
		dbm.previousDatabases = remainingDatabases
		dbm.mu.Unlock()
		slog.Info("Closing database handles in background done")
	}()
	return nil
}

// Disconnect disconnects all database connections handled by the Manager.
func (dbm *Manager) Disconnect() {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()
	var disconnectWg sync.WaitGroup
	dbm.previousDatabases = append(dbm.previousDatabases, dbm.latestDatabase)
	dbm.latestDatabase = nil
	for _, dbh := range dbm.previousDatabases {
		disconnectWg.Add(1)
		go func() {
			defer disconnectWg.Done()
			slog.Info("Wait for clients to close database handle", slog.Int64("id", dbh.id))
			dbh.wg.Wait()
			slog.Info("Closing database handle", slog.Int64("id", dbh.id))
			dbh.db.Disconnect()
		}()
	}
	disconnectWg.Wait()
}

// GetHandler returns a database instance and a callback to call when
// the actions on the database are done.
func (dbm *Manager) GetHandler() (Database, func()) {
	dbm.mu.RLock()
	defer dbm.mu.RUnlock()
	// Always return the latest connection
	dbh := dbm.latestDatabase
	if dbh == nil {
		return nil, func() {}
	}
	dbh.wg.Add(1)
	return dbh.db, func() {
		dbh.wg.Done()
	}
}
