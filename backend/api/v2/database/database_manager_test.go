package database

import (
	"sync"
	"testing"
)

type FakeDB struct {
	Database
	ID                int64
	IsConnected       bool
	AllowDisconnectWG sync.WaitGroup
	DisconnectWg      sync.WaitGroup
}

func NewFakeDB(settings Settings, id int64) *FakeDB {
	return &FakeDB{
		ID: id,
	}
}

func (f *FakeDB) Connect() error {
	f.IsConnected = true
	f.AllowDisconnectWG.Add(1)
	return nil
}

func (f *FakeDB) Disconnect() error {
	f.AllowDisconnectWG.Wait()
	f.IsConnected = false
	f.DisconnectWg.Done()
	return nil
}

func TestManager_Connect(t *testing.T) {
	testSettings := Settings{
		User:     "testuser",
		Password: "testpassword",
		Protocol: "tcp",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
	}
	dbm := NewManager(testSettings)
	fakeDB := NewFakeDB(testSettings, 0)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB
	}
	err := dbm.Connect(nil)
	if err != nil {
		t.Errorf("Connect returned error: %v", err)
	}
	if fakeDB.IsConnected != true {
		t.Errorf("Connect was not called on the database")
	}
	dbm.mu.RLock()
	defer dbm.mu.RUnlock()
	if dbm.latestDatabase.id != 0 {
		t.Errorf("latestDatabase.id is not 0")
	}
	if len(dbm.previousDatabases) > 0 {
		t.Errorf("previousDatabases is not empty")
	}
}

func TestManager_ConnectAndReplace(t *testing.T) {
	testSettings := Settings{
		User:     "testuser",
		Password: "testpassword",
		Protocol: "tcp",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
	}
	dbm := NewManager(testSettings)
	fakeDB1 := NewFakeDB(testSettings, 0)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB1
	}
	err := dbm.Connect(nil)
	if err != nil {
		t.Errorf("Connect returned error: %v", err)
	}
	fakeDB2 := NewFakeDB(testSettings, 0)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB2
	}
	fakeDB1.DisconnectWg.Add(1)
	err = dbm.ConnectAndReplace(nil)
	if err != nil {
		t.Errorf("ConnectAndReplace returned error: %v", err)
	}
	if fakeDB2.IsConnected != true {
		t.Errorf("ConnectAndReplace was not called on the database")
	}
	if !fakeDB1.IsConnected {
		t.Errorf("Original DB got disconnected")
	}
	fakeDB1.AllowDisconnectWG.Done()
	fakeDB1.DisconnectWg.Wait()
	if fakeDB1.IsConnected {
		t.Errorf("Original DB did not disconnect")
	}
	dbm.mu.RLock()
	defer dbm.mu.RUnlock()
	if dbm.latestDatabase.id != 1 {
		t.Errorf("latestDatabase.id is not 1")
	}
	if len(dbm.previousDatabases) > 0 {
		t.Errorf("previousDatabases did not cleanup")
	}
}

func TestManager_GetHandler(t *testing.T) {
	testSettings := Settings{
		User:     "testuser",
		Password: "testpassword",
		Protocol: "tcp",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
	}
	// Setup DB 1 connection
	dbm := NewManager(testSettings)
	fakeDB1 := NewFakeDB(testSettings, 0)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB1
	}
	dbm.Connect(nil)
	fakeDB1.AllowDisconnectWG.Done()
	dbh1, done1 := dbm.GetHandler()
	if dbh1.(*FakeDB).ID != 0 {
		t.Errorf("Error returning correct handler with one handler.")
	}

	// Replace with DB 2 connection
	fakeDB2 := NewFakeDB(testSettings, 1)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB2
	}
	fakeDB1.DisconnectWg.Add(1)
	dbm.ConnectAndReplace(nil)
	fakeDB2.AllowDisconnectWG.Done()
	dbh2, _ := dbm.GetHandler()
	if dbh2.(*FakeDB).ID != 1 {
		t.Errorf("Error returning correct handler with one handler.")
	}
	if len(dbm.previousDatabases) != 1 && dbm.previousDatabases[0].db.(*FakeDB).ID != 0 {
		t.Errorf("Error keeping connection 0 around while clients are still connected.")
	}

	// Finish all transactions on DB 1 -> ready for disconnect
	done1()
	fakeDB1.DisconnectWg.Wait()
	if len(dbm.previousDatabases) > 0 {
		t.Errorf("Connection 0 was not cleaned up.")
	}
}

func TestManager_Disconnect(t *testing.T) {
	testSettings := Settings{
		User:     "testuser",
		Password: "testpassword",
		Protocol: "tcp",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
	}
	dbm := NewManager(testSettings)

	// Connect DB1
	fakeDB1 := NewFakeDB(testSettings, 0)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB1
	}
	dbm.Connect(nil)
	fakeDB1.DisconnectWg.Add(1)

	// Connect DB2, moving DB1 to previousDatabases
	fakeDB2 := NewFakeDB(testSettings, 1)
	InstanceCreator = func(settings Settings) Database {
		return fakeDB2
	}
	dbm.Connect(nil)
	fakeDB2.DisconnectWg.Add(1)

	// Get a handler for DB2 to simulate an active connection
	_, done2 := dbm.GetHandler()

	// Allow DBs to be disconnected
	fakeDB1.AllowDisconnectWG.Done()
	fakeDB2.AllowDisconnectWG.Done()

	// Initiate a disconnect while one connection
	// is still active.
	var waitForDisconnect sync.WaitGroup
	waitForDisconnect.Add(1)
	go func() {
		dbm.Disconnect()
		waitForDisconnect.Done()
	}()

	// TODO check that the db is still unclosed
	if !fakeDB2.IsConnected {
		t.Errorf("DB2 is still connected")
	}

	// Release the handler for DB2
	done2()
	// Wait for both DBs to be disconnected
	fakeDB1.DisconnectWg.Wait()
	fakeDB2.DisconnectWg.Wait()
	waitForDisconnect.Wait()
}
