package database

import (
	"database/sql"
	"fmt"
	"time"

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
	// tables
	BrowserSessionTable
	ConfigurationTable
	UserTable
	SessionTable
	UserToSessionTable
}

type LogsView interface {
	GetLogs() ([]Log, error)
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
}

type SessionTable interface {
	// GetSessionsBetween returns all sessions between start and end time
	GetSessionsBetween(start, end time.Time) ([]Session, error)
}

type UserToSessionTable interface {
	// GetUsersForSession returns all users from a specific session
	GetUsersForSession(id uint) ([]UserToSession, error)
}

type UserTable interface {
	// AddUser adds a user to the database and returns an error if it failed
	AddUser(u User) (uint, error)
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
	// Returns if users are present
	UsersExist() (bool, error)
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

func (d *DBMysql) Connect() error {
	// d.config = mysql.Config{
	// 	User:                 d.User,
	// 	Passwd:               d.Password,
	// 	Net:                  d.Protocol,
	// 	Addr:                 fmt.Sprintf("%s:%s", d.Host, d.Port),
	// 	DBName:               d.DBName,
	// 	AllowNativePasswords: true,
	// }

	// db, err := sql.Open("mysql", d.config.FormatDSN())
	// if err != nil {
	// 	return err
	// }
	// d.db = db

	// Premigration steps if needed
	// - TODO change all session_type occurrences to have ID 1 and 2 instead of 0 and 1

	datetimePrecision := 2
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", d.User, d.Password, d.Protocol, d.Host, d.Port, d.DBName)
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
		return err
	}
	d.orm = orm
	d.db, err = orm.DB()
	if err != nil {
		return err
	}

	// If db exists we migrate otherwise we setup the tables
	dbIsSetup, err := d.tableExists("user")
	if err != nil {
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
	fmt.Println(d.DBName)
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
