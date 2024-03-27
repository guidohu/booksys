package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	Connect() error
	Disconnect()
	Ping() error
	IsConfigured() bool
	// views
	LogsView
	// collections
	AccountingCollection
	// tables
	BoatEngineHoursTable
	BoatFuelTable
	BoatMaintenanceTable
	BrowserSessionTable
	ConfigurationTable
	ExpenditureTable
	ExpenditureTypeTable
	HeatTable
	SessionTable
	PaymentTable
	PricingTable
	UserGroupTable
	UserRoleTable
	UserTable
	UserToSessionTable
}

type LogsView interface {
	GetLogs() ([]Log, error)
}

type AccountingCollection interface {
	GetYears() ([]uint64, error)
	GetPaymentTotal(uint64) (decimal.Decimal, error)
	GetExpenseTotal(uint64) (decimal.Decimal, error)
	GetExpenseNoRefundsTotal(uint64) (decimal.Decimal, error)
	GetHeatCostTotal(uint64) (decimal.Decimal, error)
	GetSessionPaymentTotal(uint64) (decimal.Decimal, error)
	GetSessionRefundsTotal(uint64) (decimal.Decimal, error)
	GetTransactions(uint64) ([]TransactionRow, error)
	DeleteTransaction(uint64, uint64) error
	AddExpense(Expense) error
	AddIncome(Income) error
}

type BoatEngineHoursTable interface {
	GetEngineHourLatest() (BoatEngineHour, error)
	GetEngineHours() ([]BoatEngineHour, error)
	GetEngineHoursEntry(id uint) (BoatEngineHour, error)
	AddEngineHours(b BoatEngineHour) error
	UpdateEngineHours(b BoatEngineHour) error
}

type BoatFuelTable interface {
	GetFuelEntries() ([]BoatFuel, error)
	GetFuelEntry(id uint) (BoatFuel, error)
	AddFuelEntry(e BoatFuel) error
	ChangeFuelEntry(e BoatFuel) error
}

type BoatMaintenanceTable interface {
	GetMaintenance() ([]BoatMaintenance, error)
	AddMaintenanceEntry(m BoatMaintenance) error
}
type BrowserSessionTable interface {
	AddBrowserSession(b BrowserSession) (string, error)
	GetBrowserSession(id string) (*BrowserSession, error)
	UpdateBrowserSession(b BrowserSession) error
	DeleteBrowserSession(b BrowserSession) error
}

type ConfigurationTable interface {
	GetPropertyValue(key string) (Configuration, error)
	GetAllPropertyValues() ([]Configuration, error)
	UpdateOrInsertPropertyValues(conf []Configuration) error
	GetTimezoneLocation() (*time.Location, error)
	GetMyNautiqueConfiguration() (MyNautiqueConfiguration, error)
}

type ExpenditureTable interface {
	GetUserSessionPaybacks(userID uint) (decimal.Decimal, error)
}

type ExpenditureTypeTable interface {
	GetExpenseTypes() ([]ExpenseType, error)
}

type HeatTable interface {
	GetUserHeats(userID uint, size int) ([]Heat, error)
	GetUserHeatStats(userID uint, start time.Time, end time.Time) (int64, decimal.Decimal, error)
	GetHeatInSessionCount(sessionID uint) (int, error)
	GetHeatsInSession(sessionID uint) ([]Heat, error)
	AddHeat(h *Heat) error
	DeleteHeat(heatID uint) error
	GetHeat(heatID uint) (Heat, error)
	ChangeHeat(h *Heat) error
}

type PaymentTable interface {
	GetUserSessionPayments(userID uint) (decimal.Decimal, error)
}

type PricingTable interface {
	GetPricings() ([]Pricing, error)
}

type SessionTable interface {
	GetSession(sessionID uint) (Session, error)
	// GetSessionsBetween returns all sessions between start and end time
	GetSessionsBetween(start, end time.Time) ([]Session, error)
	GetSessionsByUser(userID uint) ([]Session, error)
	CreateSession(s Session) (uint, error)
	DeleteUsersFromSession(sessionID uint) error
	DeleteSession(sessionID uint) error
	UpdateSession(s Session) error
}

type UserGroupTable interface {
	CreateUserGroup(us UserStatus, p Pricing) error
	ChangeUserGroup(us UserStatus, p Pricing) error
	DeleteUserGroup(id uint) error
	SetUserGroup(userID uint, groupID uint) error
}

type UserRoleTable interface {
	GetUserRoles() ([]UserRole, error)
}

type UserToSessionTable interface {
	// GetUsersForSession returns all users from a specific session
	GetUsersForSession(id uint) ([]UserToSession, error)
	AddSessionToUserEntry(u UserToSession) error
	DeleteSessionToUserEntry(userID uint, sessionID uint) error
}

type UserTable interface {
	// AddUser adds a user to the database and returns an error if it failed
	AddUser(u User) (uint, error)
	// Flags a user as deleted and removes all personal data.
	DeleteUserById(id uint) error
	// Returns all the admin users
	GetAdminUsers() ([]User, error)
	// GetUserByUsername find the user that has either the given username
	// or the given email address. Returns an error in case the user was not found.
	GetUserByName(name string) (User, error)
	GetUserById(id uint) (User, error)
	// ChangUserStatus changes the status of a single user.
	ChangeUserStatus(id uint, userStatusId uint) error
	// ChangeLock changes the locked status of a single user.
	ChangeLock(id uint, locked bool) error
	// Count users that have status_id of an admin
	CountAdminUsers() int64
	UpdatePassword(userID uint, user User) error
	UpdateUser(userID uint, user User) error
	// Returns if users are present
	UsersExist() (bool, error)
	// Returns true if user exists
	UserExists(id uint) bool
	// Returns all users
	GetUsers(includeDeleted bool) ([]User, error)
}

type DBMysql struct {
	// config   mysql.Config
	db       *sql.DB
	orm      *gorm.DB
	User     string
	Password string
	Protocol string
	Host     string
	Port     string
	DBName   string
}

func (d *DBMysql) String() string {
	return d.getDSN( /*hidePassword=*/ true)
}

func (d *DBMysql) getDSN(hidePassword bool) string {
	password := d.Password
	if hidePassword {
		password = "***hidden***"
	}
	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", d.User, password, d.Protocol, d.Host, d.Port, d.DBName)
}

func (d *DBMysql) Connect() error {
	// Premigration steps if needed
	// - TODO change all session_type occurrences to have ID 1 and 2 instead of 0 and 1

	datetimePrecision := 2
	dsn := d.getDSN( /*hidePassword=*/ false)
	orm, err := gorm.Open(gormMysql.New(gormMysql.Config{
		DSN:                       dsn,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}), &gorm.Config{})
	if err != nil {
		slog.Error("Cannot connect with gorm", slog.String("error", err.Error()))
		return err
	}
	d.orm = orm
	d.db, err = orm.DB()
	if err != nil {
		slog.Error("Cannot assign db handler", slog.String("error", err.Error()))
		return err
	}

	// If db exists we migrate otherwise we setup the tables
	dbIsSetup, err := d.tableExists("user")
	if err != nil {
		slog.Error("Cannot check if table 'user' exists", slog.String("error", err.Error()))
		return err
	}

	switch {
	case !dbIsSetup:
		slog.Info("New database setup detected.")
		err = d.Initialize()
	case dbIsSetup:
		slog.Info("Existing database setup detected.")
		err = d.Migrate()
	}
	if err != nil {
		slog.Error("Database setup or migration failed", err)
		return err
	}
	slog.Info("Updated DB schema")

	return d.db.Ping()
}

func (d *DBMysql) Disconnect() {
	// ormDB, err := d.orm.DB()
	// if err != nil {
	// 	slog.Error("Cannot close orm database handle", slog.String("error", err.Error()))
	// } else {
	// 	ormDB.Close()
	// }
	if d.db == nil {
		return
	}
	d.db.Close()
}

func (d *DBMysql) Ping() error {
	return d.db.Ping()
}

func (d *DBMysql) IsConfigured() bool {
	switch {
	case d.User == "":
		return false
	case d.Password == "":
		return false
	case d.Protocol == "":
		return false
	case d.Host == "":
		return false
	case d.Port == "":
		return false
	case d.DBName == "":
		return false
	}

	return true
}

func (d *DBMysql) tableExists(tableName string) (bool, error) {
	if d.DBName == "" {
		return false, errors.New("database name not provided")
	}
	rows, err := d.db.Query(fmt.Sprintf("SHOW TABLE STATUS FROM %s WHERE Name = ?", d.DBName), tableName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
