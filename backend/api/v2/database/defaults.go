package database

import (
	"golang.org/x/exp/slog"

	"gorm.io/gorm/clause"
)

// Migrate creates the database scheme. It does so by:
// - Prepare the database such that it is save to run gorm migrate.
// - Run a schema update by auto migrate.
// - Set the default values.
// - Remove all configuration settings that are no longer supported.
func (d *DBMysql) Migrate() error {
	// prepare database for migration
	err := d.migrationPreflight()
	if err != nil {
		slog.Error("migrationPreflight failed:", err)
		return err
	}

	// run schema update
	err = d.autoMigrate()
	if err != nil {
		slog.Error("autoMigrate failed:", err)
	}

	// prepare database for initialization

	// initialize values
	err = d.initializeContent()
	if err != nil {
		slog.Error("initialize database failed:", err)
	}

	// cleanup

	return nil
}

func (d *DBMysql) Initialize() error {
	// create schema
	err := d.autoMigrate()
	if err != nil {
		slog.Error("autoMigrate failed:", err)
		return err
	}

	// initialize values
	err = d.initializeContent()
	if err != nil {
		slog.Error("initialize database failed:", err)
		return err
	}

	return nil
}

// migrationPreflight is some ugly code that make existing DBs compatible
// with the gorm based schema management.
func (d *DBMysql) migrationPreflight() error {
	// change all session_type occurrences to have ID 1 and 2 instead of 0 and 1
	count := 0
	err := d.db.QueryRow("SELECT COUNT(*) FROM session_type WHERE ID = 0").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		slog.Info("migration preflight table `session_type` - skip")
	} else {
		result := d.orm.Clauses(clause.OnConflict{DoNothing: true}).Create(DefaultSessionTypes[1])
		if result.Error != nil {
			slog.Error("migration preflight table `session_type` - failed to create id 2")
			return result.Error
		}
		slog.Info("migration preflight table `session_type` - id 2 created done")

		t, err := d.db.Begin()
		if err != nil {
			slog.Error("migration preflight table `session_type` - failed to create transaction")
			return err
		}
		_, err = t.Query("UPDATE session SET type = 2 WHERE type = 1")
		if err != nil {
			slog.Error("migration preflight table `session_type` - failed to move session id 1 to id 2")
			return err
		}
		slog.Info("migration preflight table `session_type` - move session id 1 to id 2 done")

		_, err = t.Query("UPDATE session SET type = 1 WHERE type = 0")
		if err != nil {
			slog.Error("migration preflight table `session_type` - failed to move session id 0 to id 1")
			return err
		}
		slog.Info("migration preflight table `session_type` - move session id 0 to id 1 done")

		_, err = t.Query("DELETE FROM session_type WHERE id = 0")
		if err != nil {
			slog.Error("migration preflight table `session_type` - failed to remove id 0")
			return err
		}
		slog.Info("migration preflight table `session_type` - remove id 0 done")

		err = t.Commit()
		if err != nil {
			slog.Error("migration preflight table `session_type` - failed to commit transaction")
			return err
		}

		slog.Info("migration preflight table `session_type` - done")
	}

	// TODO further fixes come here

	return nil
}

func (d *DBMysql) autoMigrate() error {
	tables := []interface{}{
		&User{},
		&BoatEngineHour{},
		&BoatFuel{},
		&BoatMaintenance{},
		&BrowserSession{},
		&Expense{},
		&ExpenseType{},
		&Heat{},
		&Invitation{},
		&InvitationStatus{},
		&PasswordReset{},
		&Income{},
		&Pricing{},
		&Session{},
		&SessionType{},
		&UserRole{},
		&UserStatus{},
		&UserToSession{},
		&Configuration{},
	}
	return d.orm.AutoMigrate(tables...)
}

func (d *DBMysql) initializeContent() error {
	defaultValues := [][]interface{}{
		{DefaultUserRoles},
		{DefaultUserStatus},
		{DefaultSessionTypes},
		{DefaultPricing},
		{DefaultInvitationStatus},
		{DefaultExpenseTypes},
		{DefaultConfiguration},
	}

	// Setup default values
	for _, i := range defaultValues {
		for _, r := range i {
			result := d.orm.Clauses(clause.OnConflict{DoNothing: true}).Create(r)
			if result.Error != nil {
				return result.Error
			}
		}
	}
	return nil
}

func (d *DBMysql) cleanup() error {
	return nil
}

func SetupInitialValues() error {
	return nil
}

func SetupUpserts() error {
	return nil
}

// DefaultUserRoles contains all the access levels supported by the
// app.
var DefaultUserRoles = []UserRole{
	{
		ID:          1,
		Name:        "guest",
		Description: "Guest User Permissions",
	},
	{
		ID:          2,
		Name:        "member",
		Description: "Member User Permissions",
	},
	{
		ID:          3,
		Name:        "admin",
		Description: "Administrator User Permissions",
	},
}

// TODO minimum set of user status required for the app to work
var DefaultUserStatus = []UserStatus{
	{
		ID:          1,
		Name:        "Guest",
		Description: "Guest Access",
		UserRoleID:  1,
	},
	{
		ID:          2,
		Name:        "Member",
		Description: "Member Access",
		UserRoleID:  2,
	},
	{
		ID:          3,
		Name:        "Admin",
		Description: "Administrator and Boat Community Access",
		UserRoleID:  3,
	},
	{
		ID:          4,
		Name:        "Course",
		Description: "Course Status and only Guest Access Permissions",
		UserRoleID:  1,
	},
	{
		ID:          5,
		Name:        "Partner",
		Description: "Partners with Member Permissions",
		UserRoleID:  2,
	},
}

const (
	DefaultSessionType = iota + 1
	CourseSessionType
)

// TODO change all session type occurrences from 0 to 1 and from 1 to two in the code and in the database
var DefaultSessionTypes = []SessionType{
	{
		ID:      DefaultSessionType,
		Name:    "default",
		Comment: "A normal session.",
	},
	{
		ID:      CourseSessionType,
		Name:    "course",
		Comment: "A course session.",
	},
}

var DefaultPricing = []Pricing{
	{
		ID:             1,
		UserStatusID:   1,
		PricePerMinute: 2.80,
		Comment:        "Guest Price",
	},
	{
		ID:             2,
		UserStatusID:   2,
		PricePerMinute: 2.80,
		Comment:        "Member Price",
	},
	{
		ID:             3,
		UserStatusID:   3,
		PricePerMinute: 1.30,
		Comment:        "Boat Community Price",
	},
	{
		ID:             4,
		UserStatusID:   4,
		PricePerMinute: 0.00,
		Comment:        "Course Price",
	},
	{
		ID:             5,
		UserStatusID:   5,
		PricePerMinute: 1.30,
		Comment:        "Partner of Boat Community",
	},
}

var DefaultInvitationStatus = []InvitationStatus{
	{
		ID:      1,
		Name:    "full",
		Comment: "Session is full by now",
	},
	{
		ID:      2,
		Name:    "query",
		Comment: "The invitation is still valid",
	},
	{
		ID:      3,
		Name:    "accepted",
		Comment: "User accepted the invitation",
	},
	{
		ID:      4,
		Name:    "declined",
		Comment: "User declined the invitation",
	},
	{
		ID:      5,
		Name:    "unknown",
		Comment: "No invitation sent yet.",
	},
}

var DefaultExpenseTypes = []ExpenseType{
	{
		ID:      0,
		Name:    "fuel",
		Comment: "Costs for fuel",
	},
	{
		ID:      1,
		Name:    "maintenance",
		Comment: "Costs for maintenance.",
	},
	{
		ID:      2,
		Name:    "material",
		Comment: "Material which is needed for the boat.",
	},
	{
		ID:      3,
		Name:    "invest",
		Comment: "Investment in new Features, ...",
	},
	{
		ID:      4,
		Name:    "session",
		Comment: "Payments for sessions.",
	},
	{
		ID:      5,
		Name:    "other",
		Comment: "Other costs.",
	},
	{
		ID:      6,
		Name:    "membership fee",
		Comment: "Membership Fee.",
	},
	{
		ID:      7,
		Name:    "salary",
		Comment: "Compensation for towing people or do whatever in the name of the community.",
	},
	{
		ID:      8,
		Name:    "owners refund",
		Comment: "Payback of investments from owners.",
	},
	{
		ID:      9,
		Name:    "fuel bill",
		Comment: "Costs for fuel paid through invoices.",
	},
}

var DefaultConfiguration = []Configuration{
	{
		ID:       0,
		Property: "schema.version",
		Value:    "1.16",
	},
	{
		ID:       1,
		Property: "browser.session.timeout.default",
		Value:    "10800",
	},
	{
		ID:       2,
		Property: "browser.session.timeout.max",
		Value:    "604800",
	},
	{
		ID:       3,
		Property: "location.longitude",
		Value:    "8.542939",
	},
	{
		ID:       4,
		Property: "location.latitude",
		Value:    "47.367658",
	},
	{
		ID:       5,
		Property: "location.gmt_offset",
		Value:    "1",
	},
	{
		ID:       6,
		Property: "location.timezone",
		Value:    "Europe/Berlin",
	},
	{
		ID:       7,
		Property: "business.day.start",
		Value:    "08:00:00",
	},
	{
		ID:       8,
		Property: "business.day.end",
		Value:    "21:00:00",
	},
	{
		ID:       9,
		Property: "business.day.startatsunrise",
		Value:    "false",
	},
	{
		ID:       10,
		Property: "business.day.endatsunset",
		Value:    "false",
	},
	{
		ID:       11,
		Property: "session.cancel.graceperiod",
		Value:    "86400",
	},
	{
		ID:       12,
		Property: "recaptcha.privatekey",
		Value:    "",
	},
	{
		ID:       13,
		Property: "recaptcha.publickey",
		Value:    "",
	},
	{
		ID:       14,
		Property: "currency",
		Value:    "CHF",
	},
	{
		ID:       15,
		Property: "location.address",
		Value:    "",
	},
	{
		ID:       16,
		Property: "location.map",
		Value:    "",
	},
	{
		ID:       17,
		Property: "payment.account.owner",
		Value:    "",
	},
	{
		ID:       18,
		Property: "payment.account.iban",
		Value:    "",
	},
	{
		ID:       19,
		Property: "payment.account.bic",
		Value:    "",
	},
	{
		ID:       20,
		Property: "payment.account.comment",
		Value:    "",
	},
	{
		ID:       21,
		Property: "smtp.sender",
		Value:    "",
	},
	{
		ID:       22,
		Property: "smtp.server",
		Value:    "",
	},
	{
		ID:       23,
		Property: "smtp.username",
		Value:    "",
	},
	{
		ID:       24,
		Property: "smtp.password",
		Value:    "",
	},
	{
		ID:       25,
		Property: "logo.file",
		Value:    "",
	},
	{
		ID:       26,
		Property: "engine.hour.format",
		Value:    "hh.h",
	},
	{
		ID:       27,
		Property: "fuel.payment.type",
		Value:    "instant",
	},
}
