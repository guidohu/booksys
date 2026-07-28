// Package dbtest provides test doubles for the database layer.
//
// It is meant to be imported from the tests of packages that depend on
// server/database, for example the HTTP handlers or the configuration.
//
// Note: The tests inside package database itself cannot use this package,
// importing it from there would be an import cycle. Those tests have to
// define their own doubles.
package dbtest

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"server/database"
)

// ErrNotFound is the generic "record does not exist" error returned by the
// fake for lookups that have no behaviour configured.
var ErrNotFound = errors.New("record not found")

// FakeDB is a full in-memory implementation of the database.Database
// interface. Every method is backed by an optional function field so that a
// test only has to provide the behaviour it actually cares about. Methods
// without a function return zero values and no error, lookups that address a
// single record return ErrNotFound.
type FakeDB struct {
	// Connection handling
	ConnectFn       func() error
	DisconnectFn    func() error
	PingFn          func() error
	IsConfiguredFn  func() bool
	IsInitializedFn func() (bool, error)

	// Logs
	GetLogsFn func(currency string) ([]database.Log, error)

	// Accounting
	GetYearsFn                 func() ([]uint64, error)
	GetPaymentTotalFn          func(year uint64) (decimal.Decimal, error)
	GetExpenseTotalFn          func(year uint64) (decimal.Decimal, error)
	GetExpenseNoRefundsTotalFn func(year uint64) (decimal.Decimal, error)
	GetHeatCostTotalFn         func(year uint64) (decimal.Decimal, error)
	GetSessionPaymentTotalFn   func(year uint64) (decimal.Decimal, error)
	GetSessionRefundsTotalFn   func(year uint64) (decimal.Decimal, error)
	GetTransactionsFn          func(year uint64) ([]database.TransactionRow, error)
	DeleteTransactionFn        func(tableID, rowID uint64) error
	AddExpenseFn               func(database.Expense) error
	AddIncomeFn                func(database.Income) error

	// Boat engine hours
	GetEngineHourLatestFn func() (database.BoatEngineHour, error)
	GetEngineHoursFn      func() ([]database.BoatEngineHour, error)
	GetEngineHoursEntryFn func(id uint) (database.BoatEngineHour, error)
	AddEngineHoursFn      func(database.BoatEngineHour) error
	UpdateEngineHoursFn   func(database.BoatEngineHour) error

	// Boat fuel
	GetFuelEntriesFn  func() ([]database.BoatFuel, error)
	GetFuelEntryFn    func(id uint) (database.BoatFuel, error)
	AddFuelEntryFn    func(database.BoatFuel) error
	ChangeFuelEntryFn func(database.BoatFuel) error
	RemoveFuelEntryFn func(id uint) error

	// Boat maintenance
	GetMaintenanceFn      func() ([]database.BoatMaintenance, error)
	AddMaintenanceEntryFn func(database.BoatMaintenance) error

	// Browser sessions
	AddBrowserSessionFn    func(database.BrowserSession) (string, error)
	GetBrowserSessionFn    func(id string) (*database.BrowserSession, error)
	UpdateBrowserSessionFn func(database.BrowserSession) error
	DeleteBrowserSessionFn func(database.BrowserSession) error

	// Configuration
	DeletePropertyFn               func(key string) error
	GetPropertyValueFn             func(key string) (database.Configuration, error)
	GetAllPropertyValuesFn         func() ([]database.Configuration, error)
	GetAllPropertyValuesMapFn      func() (map[string]database.Configuration, error)
	UpdateOrInsertPropertyValuesFn func([]database.Configuration) error
	GetTimezoneLocationFn          func() (*time.Location, error)
	GetEmailConfigurationFn        func() (database.EmailConfiguration, error)
	GetMyNautiqueConfigurationFn   func() (database.MyNautiqueConfiguration, error)
	GetConfigurationVersionFn      func() (database.ConfigurationVersion, error)

	// Expenditure / payment
	GetUserSessionPaybacksFn func(userID uint) (decimal.Decimal, error)
	GetUserSessionPaymentsFn func(userID uint) (decimal.Decimal, error)
	GetExpenseTypesFn        func() ([]database.ExpenseType, error)

	// Heats
	GetUserHeatsFn          func(userID uint, size int) ([]database.Heat, error)
	GetUserHeatsBySessionFn func(sessionID, userID uint, size int) ([]database.Heat, error)
	GetUserHeatStatsFn      func(userID uint, start, end time.Time) (int64, decimal.Decimal, error)
	GetHeatInSessionCountFn func(sessionID uint) (int, error)
	GetHeatsInSessionFn     func(sessionID uint) ([]database.Heat, error)
	AddHeatFn               func(*database.Heat) error
	DeleteHeatFn            func(heatID uint) error
	GetHeatFn               func(heatID uint) (database.Heat, error)
	ChangeHeatFn            func(*database.Heat) error

	// Password reset
	AddPasswordResetTokenFn          func(database.PasswordReset) error
	GetPasswordResetEntryFn          func(userID uint, token string) (database.PasswordReset, error)
	InvalidatePasswordResetEntriesFn func(userID uint) error

	// Pricing
	GetPricingsFn                func() ([]database.Pricing, error)
	GetUserStatusToPricingsMapFn func() (map[uint]database.Pricing, error)

	// Sessions
	GetSessionFn             func(sessionID uint) (database.Session, error)
	GetSessionsBetweenFn     func(start, end time.Time) ([]database.Session, error)
	GetSessionsByUserFn      func(userID uint) ([]database.Session, error)
	CreateSessionFn          func(database.Session) (uint, error)
	DeleteUsersFromSessionFn func(sessionID uint) error
	DeleteSessionFn          func(sessionID uint) error
	UpdateSessionFn          func(database.Session) error

	// User groups / roles
	CreateUserGroupFn func(database.UserStatus, database.Pricing) error
	ChangeUserGroupFn func(database.UserStatus, database.Pricing) error
	DeleteUserGroupFn func(id uint) error
	SetUserGroupFn    func(userID, groupID uint) error
	GetUserRolesFn    func() ([]database.UserRole, error)

	// User to session
	GetUsersForSessionFn       func(id uint) ([]database.UserToSession, error)
	AddSessionToUserEntryFn    func(database.UserToSession) error
	DeleteSessionToUserEntryFn func(userID, sessionID uint) error

	// Users
	AddUserFn          func(database.User) (uint, error)
	DeleteUserByIdFn   func(id uint) error
	GetAdminUsersFn    func() ([]database.User, error)
	IsAdminUserFn      func(userID uint) (bool, error)
	GetUserByNameFn    func(name string) (database.User, error)
	GetUserByIdFn      func(id uint) (database.User, error)
	ChangeUserStatusFn func(id, userStatusID uint) error
	ChangeLockFn       func(id uint, locked bool) error
	CountAdminUsersFn  func() int64
	UpdatePasswordFn   func(userID uint, user database.User) error
	UpdateUserFn       func(userID uint, user database.User) error
	UsersExistFn       func() (bool, error)
	UserExistsFn       func(id uint) bool
	GetUsersFn         func(includeDeleted bool) ([]database.User, error)
}

// Compile time check that the fake covers the whole interface.
var _ database.Database = (*FakeDB)(nil)

func (f *FakeDB) Connect() error {
	if f.ConnectFn != nil {
		return f.ConnectFn()
	}
	return nil
}

func (f *FakeDB) Disconnect() error {
	if f.DisconnectFn != nil {
		return f.DisconnectFn()
	}
	return nil
}

func (f *FakeDB) Ping() error {
	if f.PingFn != nil {
		return f.PingFn()
	}
	return nil
}

func (f *FakeDB) IsConfigured() bool {
	if f.IsConfiguredFn != nil {
		return f.IsConfiguredFn()
	}
	return true
}

func (f *FakeDB) IsInitialized() (bool, error) {
	if f.IsInitializedFn != nil {
		return f.IsInitializedFn()
	}
	return true, nil
}

func (f *FakeDB) GetLogs(currency string) ([]database.Log, error) {
	if f.GetLogsFn != nil {
		return f.GetLogsFn(currency)
	}
	return nil, nil
}

func (f *FakeDB) GetYears() ([]uint64, error) {
	if f.GetYearsFn != nil {
		return f.GetYearsFn()
	}
	return nil, nil
}

func (f *FakeDB) GetPaymentTotal(year uint64) (decimal.Decimal, error) {
	if f.GetPaymentTotalFn != nil {
		return f.GetPaymentTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetExpenseTotal(year uint64) (decimal.Decimal, error) {
	if f.GetExpenseTotalFn != nil {
		return f.GetExpenseTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetExpenseNoRefundsTotal(year uint64) (decimal.Decimal, error) {
	if f.GetExpenseNoRefundsTotalFn != nil {
		return f.GetExpenseNoRefundsTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetHeatCostTotal(year uint64) (decimal.Decimal, error) {
	if f.GetHeatCostTotalFn != nil {
		return f.GetHeatCostTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetSessionPaymentTotal(year uint64) (decimal.Decimal, error) {
	if f.GetSessionPaymentTotalFn != nil {
		return f.GetSessionPaymentTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetSessionRefundsTotal(year uint64) (decimal.Decimal, error) {
	if f.GetSessionRefundsTotalFn != nil {
		return f.GetSessionRefundsTotalFn(year)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetTransactions(year uint64) ([]database.TransactionRow, error) {
	if f.GetTransactionsFn != nil {
		return f.GetTransactionsFn(year)
	}
	return nil, nil
}

func (f *FakeDB) DeleteTransaction(tableID, rowID uint64) error {
	if f.DeleteTransactionFn != nil {
		return f.DeleteTransactionFn(tableID, rowID)
	}
	return nil
}

func (f *FakeDB) AddExpense(e database.Expense) error {
	if f.AddExpenseFn != nil {
		return f.AddExpenseFn(e)
	}
	return nil
}

func (f *FakeDB) AddIncome(i database.Income) error {
	if f.AddIncomeFn != nil {
		return f.AddIncomeFn(i)
	}
	return nil
}

func (f *FakeDB) GetEngineHourLatest() (database.BoatEngineHour, error) {
	if f.GetEngineHourLatestFn != nil {
		return f.GetEngineHourLatestFn()
	}
	return database.BoatEngineHour{}, nil
}

func (f *FakeDB) GetEngineHours() ([]database.BoatEngineHour, error) {
	if f.GetEngineHoursFn != nil {
		return f.GetEngineHoursFn()
	}
	return nil, nil
}

func (f *FakeDB) GetEngineHoursEntry(id uint) (database.BoatEngineHour, error) {
	if f.GetEngineHoursEntryFn != nil {
		return f.GetEngineHoursEntryFn(id)
	}
	return database.BoatEngineHour{}, nil
}

func (f *FakeDB) AddEngineHours(b database.BoatEngineHour) error {
	if f.AddEngineHoursFn != nil {
		return f.AddEngineHoursFn(b)
	}
	return nil
}

func (f *FakeDB) UpdateEngineHours(b database.BoatEngineHour) error {
	if f.UpdateEngineHoursFn != nil {
		return f.UpdateEngineHoursFn(b)
	}
	return nil
}

func (f *FakeDB) GetFuelEntries() ([]database.BoatFuel, error) {
	if f.GetFuelEntriesFn != nil {
		return f.GetFuelEntriesFn()
	}
	return nil, nil
}

func (f *FakeDB) GetFuelEntry(id uint) (database.BoatFuel, error) {
	if f.GetFuelEntryFn != nil {
		return f.GetFuelEntryFn(id)
	}
	return database.BoatFuel{}, nil
}

func (f *FakeDB) AddFuelEntry(e database.BoatFuel) error {
	if f.AddFuelEntryFn != nil {
		return f.AddFuelEntryFn(e)
	}
	return nil
}

func (f *FakeDB) ChangeFuelEntry(e database.BoatFuel) error {
	if f.ChangeFuelEntryFn != nil {
		return f.ChangeFuelEntryFn(e)
	}
	return nil
}

func (f *FakeDB) RemoveFuelEntry(id uint) error {
	if f.RemoveFuelEntryFn != nil {
		return f.RemoveFuelEntryFn(id)
	}
	return nil
}

func (f *FakeDB) GetMaintenance() ([]database.BoatMaintenance, error) {
	if f.GetMaintenanceFn != nil {
		return f.GetMaintenanceFn()
	}
	return nil, nil
}

func (f *FakeDB) AddMaintenanceEntry(m database.BoatMaintenance) error {
	if f.AddMaintenanceEntryFn != nil {
		return f.AddMaintenanceEntryFn(m)
	}
	return nil
}

func (f *FakeDB) AddBrowserSession(b database.BrowserSession) (string, error) {
	if f.AddBrowserSessionFn != nil {
		return f.AddBrowserSessionFn(b)
	}
	return b.SessionSecret, nil
}

func (f *FakeDB) GetBrowserSession(id string) (*database.BrowserSession, error) {
	if f.GetBrowserSessionFn != nil {
		return f.GetBrowserSessionFn(id)
	}
	return nil, ErrNotFound
}

func (f *FakeDB) UpdateBrowserSession(b database.BrowserSession) error {
	if f.UpdateBrowserSessionFn != nil {
		return f.UpdateBrowserSessionFn(b)
	}
	return nil
}

func (f *FakeDB) DeleteBrowserSession(b database.BrowserSession) error {
	if f.DeleteBrowserSessionFn != nil {
		return f.DeleteBrowserSessionFn(b)
	}
	return nil
}

func (f *FakeDB) DeleteProperty(key string) error {
	if f.DeletePropertyFn != nil {
		return f.DeletePropertyFn(key)
	}
	return nil
}

func (f *FakeDB) GetPropertyValue(key string) (database.Configuration, error) {
	if f.GetPropertyValueFn != nil {
		return f.GetPropertyValueFn(key)
	}
	return database.Configuration{}, nil
}

func (f *FakeDB) GetAllPropertyValues() ([]database.Configuration, error) {
	if f.GetAllPropertyValuesFn != nil {
		return f.GetAllPropertyValuesFn()
	}
	return nil, nil
}

func (f *FakeDB) GetAllPropertyValuesMap() (map[string]database.Configuration, error) {
	if f.GetAllPropertyValuesMapFn != nil {
		return f.GetAllPropertyValuesMapFn()
	}
	return map[string]database.Configuration{}, nil
}

func (f *FakeDB) UpdateOrInsertPropertyValues(conf []database.Configuration) error {
	if f.UpdateOrInsertPropertyValuesFn != nil {
		return f.UpdateOrInsertPropertyValuesFn(conf)
	}
	return nil
}

func (f *FakeDB) GetTimezoneLocation() (*time.Location, error) {
	if f.GetTimezoneLocationFn != nil {
		return f.GetTimezoneLocationFn()
	}
	return time.UTC, nil
}

func (f *FakeDB) GetEmailConfiguration() (database.EmailConfiguration, error) {
	if f.GetEmailConfigurationFn != nil {
		return f.GetEmailConfigurationFn()
	}
	return database.EmailConfiguration{}, nil
}

func (f *FakeDB) GetMyNautiqueConfiguration() (database.MyNautiqueConfiguration, error) {
	if f.GetMyNautiqueConfigurationFn != nil {
		return f.GetMyNautiqueConfigurationFn()
	}
	return database.MyNautiqueConfiguration{}, nil
}

func (f *FakeDB) GetConfigurationVersion() (database.ConfigurationVersion, error) {
	if f.GetConfigurationVersionFn != nil {
		return f.GetConfigurationVersionFn()
	}
	return database.ConfigurationVersion{}, nil
}

func (f *FakeDB) GetUserSessionPaybacks(userID uint) (decimal.Decimal, error) {
	if f.GetUserSessionPaybacksFn != nil {
		return f.GetUserSessionPaybacksFn(userID)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetUserSessionPayments(userID uint) (decimal.Decimal, error) {
	if f.GetUserSessionPaymentsFn != nil {
		return f.GetUserSessionPaymentsFn(userID)
	}
	return decimal.Zero, nil
}

func (f *FakeDB) GetExpenseTypes() ([]database.ExpenseType, error) {
	if f.GetExpenseTypesFn != nil {
		return f.GetExpenseTypesFn()
	}
	return nil, nil
}

func (f *FakeDB) GetUserHeats(userID uint, size int) ([]database.Heat, error) {
	if f.GetUserHeatsFn != nil {
		return f.GetUserHeatsFn(userID, size)
	}
	return nil, nil
}

func (f *FakeDB) GetUserHeatsBySession(sessionID, userID uint, size int) ([]database.Heat, error) {
	if f.GetUserHeatsBySessionFn != nil {
		return f.GetUserHeatsBySessionFn(sessionID, userID, size)
	}
	return nil, nil
}

func (f *FakeDB) GetUserHeatStats(userID uint, start, end time.Time) (int64, decimal.Decimal, error) {
	if f.GetUserHeatStatsFn != nil {
		return f.GetUserHeatStatsFn(userID, start, end)
	}
	return 0, decimal.Zero, nil
}

func (f *FakeDB) GetHeatInSessionCount(sessionID uint) (int, error) {
	if f.GetHeatInSessionCountFn != nil {
		return f.GetHeatInSessionCountFn(sessionID)
	}
	return 0, nil
}

func (f *FakeDB) GetHeatsInSession(sessionID uint) ([]database.Heat, error) {
	if f.GetHeatsInSessionFn != nil {
		return f.GetHeatsInSessionFn(sessionID)
	}
	return nil, nil
}

func (f *FakeDB) AddHeat(h *database.Heat) error {
	if f.AddHeatFn != nil {
		return f.AddHeatFn(h)
	}
	return nil
}

func (f *FakeDB) DeleteHeat(heatID uint) error {
	if f.DeleteHeatFn != nil {
		return f.DeleteHeatFn(heatID)
	}
	return nil
}

func (f *FakeDB) GetHeat(heatID uint) (database.Heat, error) {
	if f.GetHeatFn != nil {
		return f.GetHeatFn(heatID)
	}
	return database.Heat{}, nil
}

func (f *FakeDB) ChangeHeat(h *database.Heat) error {
	if f.ChangeHeatFn != nil {
		return f.ChangeHeatFn(h)
	}
	return nil
}

func (f *FakeDB) AddPasswordResetToken(entry database.PasswordReset) error {
	if f.AddPasswordResetTokenFn != nil {
		return f.AddPasswordResetTokenFn(entry)
	}
	return nil
}

func (f *FakeDB) GetPasswordResetEntry(userID uint, token string) (database.PasswordReset, error) {
	if f.GetPasswordResetEntryFn != nil {
		return f.GetPasswordResetEntryFn(userID, token)
	}
	return database.PasswordReset{}, nil
}

func (f *FakeDB) InvalidatePasswordResetEntries(userID uint) error {
	if f.InvalidatePasswordResetEntriesFn != nil {
		return f.InvalidatePasswordResetEntriesFn(userID)
	}
	return nil
}

func (f *FakeDB) GetPricings() ([]database.Pricing, error) {
	if f.GetPricingsFn != nil {
		return f.GetPricingsFn()
	}
	return nil, nil
}

func (f *FakeDB) GetUserStatusToPricingsMap() (map[uint]database.Pricing, error) {
	if f.GetUserStatusToPricingsMapFn != nil {
		return f.GetUserStatusToPricingsMapFn()
	}
	return map[uint]database.Pricing{}, nil
}

func (f *FakeDB) GetSession(sessionID uint) (database.Session, error) {
	if f.GetSessionFn != nil {
		return f.GetSessionFn(sessionID)
	}
	return database.Session{}, nil
}

func (f *FakeDB) GetSessionsBetween(start, end time.Time) ([]database.Session, error) {
	if f.GetSessionsBetweenFn != nil {
		return f.GetSessionsBetweenFn(start, end)
	}
	return nil, nil
}

func (f *FakeDB) GetSessionsByUser(userID uint) ([]database.Session, error) {
	if f.GetSessionsByUserFn != nil {
		return f.GetSessionsByUserFn(userID)
	}
	return nil, nil
}

func (f *FakeDB) CreateSession(s database.Session) (uint, error) {
	if f.CreateSessionFn != nil {
		return f.CreateSessionFn(s)
	}
	return 0, nil
}

func (f *FakeDB) DeleteUsersFromSession(sessionID uint) error {
	if f.DeleteUsersFromSessionFn != nil {
		return f.DeleteUsersFromSessionFn(sessionID)
	}
	return nil
}

func (f *FakeDB) DeleteSession(sessionID uint) error {
	if f.DeleteSessionFn != nil {
		return f.DeleteSessionFn(sessionID)
	}
	return nil
}

func (f *FakeDB) UpdateSession(s database.Session) error {
	if f.UpdateSessionFn != nil {
		return f.UpdateSessionFn(s)
	}
	return nil
}

func (f *FakeDB) CreateUserGroup(us database.UserStatus, p database.Pricing) error {
	if f.CreateUserGroupFn != nil {
		return f.CreateUserGroupFn(us, p)
	}
	return nil
}

func (f *FakeDB) ChangeUserGroup(us database.UserStatus, p database.Pricing) error {
	if f.ChangeUserGroupFn != nil {
		return f.ChangeUserGroupFn(us, p)
	}
	return nil
}

func (f *FakeDB) DeleteUserGroup(id uint) error {
	if f.DeleteUserGroupFn != nil {
		return f.DeleteUserGroupFn(id)
	}
	return nil
}

func (f *FakeDB) SetUserGroup(userID, groupID uint) error {
	if f.SetUserGroupFn != nil {
		return f.SetUserGroupFn(userID, groupID)
	}
	return nil
}

func (f *FakeDB) GetUserRoles() ([]database.UserRole, error) {
	if f.GetUserRolesFn != nil {
		return f.GetUserRolesFn()
	}
	return nil, nil
}

func (f *FakeDB) GetUsersForSession(id uint) ([]database.UserToSession, error) {
	if f.GetUsersForSessionFn != nil {
		return f.GetUsersForSessionFn(id)
	}
	return nil, nil
}

func (f *FakeDB) AddSessionToUserEntry(u database.UserToSession) error {
	if f.AddSessionToUserEntryFn != nil {
		return f.AddSessionToUserEntryFn(u)
	}
	return nil
}

func (f *FakeDB) DeleteSessionToUserEntry(userID, sessionID uint) error {
	if f.DeleteSessionToUserEntryFn != nil {
		return f.DeleteSessionToUserEntryFn(userID, sessionID)
	}
	return nil
}

func (f *FakeDB) AddUser(u database.User) (uint, error) {
	if f.AddUserFn != nil {
		return f.AddUserFn(u)
	}
	return 0, nil
}

func (f *FakeDB) DeleteUserById(id uint) error {
	if f.DeleteUserByIdFn != nil {
		return f.DeleteUserByIdFn(id)
	}
	return nil
}

func (f *FakeDB) GetAdminUsers() ([]database.User, error) {
	if f.GetAdminUsersFn != nil {
		return f.GetAdminUsersFn()
	}
	return nil, nil
}

func (f *FakeDB) IsAdminUser(userID uint) (bool, error) {
	if f.IsAdminUserFn != nil {
		return f.IsAdminUserFn(userID)
	}
	return false, nil
}

func (f *FakeDB) GetUserByName(name string) (database.User, error) {
	if f.GetUserByNameFn != nil {
		return f.GetUserByNameFn(name)
	}
	return database.User{}, ErrNotFound
}

func (f *FakeDB) GetUserById(id uint) (database.User, error) {
	if f.GetUserByIdFn != nil {
		return f.GetUserByIdFn(id)
	}
	return database.User{}, ErrNotFound
}

func (f *FakeDB) ChangeUserStatus(id, userStatusID uint) error {
	if f.ChangeUserStatusFn != nil {
		return f.ChangeUserStatusFn(id, userStatusID)
	}
	return nil
}

func (f *FakeDB) ChangeLock(id uint, locked bool) error {
	if f.ChangeLockFn != nil {
		return f.ChangeLockFn(id, locked)
	}
	return nil
}

func (f *FakeDB) CountAdminUsers() int64 {
	if f.CountAdminUsersFn != nil {
		return f.CountAdminUsersFn()
	}
	return 0
}

func (f *FakeDB) UpdatePassword(userID uint, user database.User) error {
	if f.UpdatePasswordFn != nil {
		return f.UpdatePasswordFn(userID, user)
	}
	return nil
}

func (f *FakeDB) UpdateUser(userID uint, user database.User) error {
	if f.UpdateUserFn != nil {
		return f.UpdateUserFn(userID, user)
	}
	return nil
}

func (f *FakeDB) UsersExist() (bool, error) {
	if f.UsersExistFn != nil {
		return f.UsersExistFn()
	}
	return false, nil
}

func (f *FakeDB) UserExists(id uint) bool {
	if f.UserExistsFn != nil {
		return f.UserExistsFn(id)
	}
	return false
}

func (f *FakeDB) GetUsers(includeDeleted bool) ([]database.User, error) {
	if f.GetUsersFn != nil {
		return f.GetUsersFn(includeDeleted)
	}
	return nil, nil
}
