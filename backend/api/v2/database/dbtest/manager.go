package dbtest

import (
	"testing"

	"server/database"
)

// NewManager returns a connected database.Manager that hands out the provided
// database to every caller of GetHandler.
//
// It swaps the package level database.InstanceCreator for the duration of the
// test and restores it afterwards, so tests using it must not run in parallel
// with each other.
func NewManager(tb testing.TB, db database.Database) *database.Manager {
	tb.Helper()
	previous := database.InstanceCreator
	tb.Cleanup(func() { database.InstanceCreator = previous })
	database.InstanceCreator = func(database.Settings) database.Database { return db }

	manager := database.NewManager(database.Settings{})
	if err := manager.Connect(nil); err != nil {
		tb.Fatalf("cannot connect test database manager: %v", err)
	}
	return manager
}
