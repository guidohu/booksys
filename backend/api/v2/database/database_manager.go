package database

import (
	"fmt"
	"sync"

	"golang.org/x/exp/slog"
)

type Manager struct {
	generations    int64
	settings       Settings
	latestDatabase *Handle
	// Connections we want to close but there might still be clients
	// using those connections.
	previousDatabases []*Handle
	mu                sync.RWMutex
}

type Handle struct {
	id int64
	db Database
	wg sync.WaitGroup
}

// NewDB function that gets called when a new Database
// is created by the DBManager. Can be overridden for testing.
var InstanceCreator func(settings Settings) Database = func(settings Settings) Database {
	return NewDBMysql(settings)
}

func NewManager(settings Settings) *Manager {
	return &Manager{
		settings: settings,
	}
}

func (dbm *Manager) UpdateSettings(settings Settings) {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()
	dbm.settings = settings
}

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

func (dbm *Manager) ConnectAndReplace(settings *Settings) error {
	err := dbm.Connect(settings)
	if err != nil {
		return err
	}
	go func() {
		slog.Info("Closing database handles in background")
		for _, dbh := range dbm.previousDatabases {
			dbh.wg.Wait()
			slog.Info("Closing database handle", slog.Int64("id", dbh.id))
			dbm.mu.Lock()
			dbh.db.Disconnect()
			dbh.db = nil
			dbm.mu.Unlock()
		}
		dbm.mu.Lock()
		remainingDatabases := []*Handle{}

		for i, dbh := range dbm.previousDatabases {
			fmt.Printf("DEBUG - i: %d, dbh: %+v", i, dbh)
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

// Disconnect disconnects all database connections handled by the DBManager
func (dbm *Manager) Disconnect() {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()
	var disconnectWg sync.WaitGroup
	dbm.previousDatabases = append(dbm.previousDatabases, dbm.latestDatabase)
	dbm.latestDatabase = nil
	for _, dbh := range dbm.previousDatabases {
		disconnectWg.Add(1)
		go func(h *Handle) {
			slog.Info("Wait for clients to close database handle", slog.Int64("id", h.id))
			h.wg.Wait()
			slog.Info("Closing database handle", slog.Int64("id", h.id))
			h.db.Disconnect()
			disconnectWg.Done()
		}(dbh)
	}
	disconnectWg.Wait()
}

func (dbm *Manager) IsInitialized() {

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
