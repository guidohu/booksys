package database

import (
	"github.com/shopspring/decimal"
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
		return err
	}

	// prepare database for initialization

	// initialize values
	err = d.initializeContent()
	if err != nil {
		slog.Error("initialize database failed:", err)
		return err
	}

	// post schema update tasks
	err = d.cleanup()
	if err != nil {
		slog.Error("post migration tasks failed:", err)
	}

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

	// increase id in expenditure_type
	count = 0
	err = d.db.QueryRow("SELECT COUNT(*) FROM expenditure_type WHERE ID = 0").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		slog.Info("migration preflight table `expenditure_type` - skip")
	} else {
		t, err := d.db.Begin()
		if err != nil {
			slog.Error("migration preflight table `expenditure_type` - failed to create transaction")
			return err
		}
		for i, _ := range DefaultExpenseTypes {
			newID := len(DefaultExpenseTypes) - i // idx is current ID + 1
			_, err = t.Query("UPDATE expenditure_type SET id = ? WHERE id = ?", newID, newID-1)
			if err != nil {
				slog.Error("migration preflight table `expenditure_type` - failed to move type ID from N to N+1 with", slog.Int("N", newID-1))
				return err
			}
			slog.Info("migration preflight table `expenditure_type` - move type ID from N to N+1 with", slog.Int("N", newID-1))
		}
		err = t.Commit()
		if err != nil {
			slog.Error("migration preflight table `expenditure_type` - failed to commit transaction")
			return err
		}
		slog.Info("migration preflight table `expenditure_type` - done")
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

	// set is_discounted where a discount was provided
	// for boat_fuel
	err := d.orm.Exec(`UPDATE boat_fuel 
		SET is_discounted = 1
		WHERE 
			cost_chf_brutto IS NOT NULL 
			AND cost_chf_brutto <> cost_chf;`).Error
	if err != nil {
		slog.Error("migration cleanup table `boat_fuel` - failed to set is_discounted")
		return err
	}
	slog.Info("migration cleanup table `boat_fuel` - set is_discounted done")

	return nil
}

func SetupInitialValues() error {
	return nil
}

func SetupUpserts() error {
	return nil
}

const (
	UserRoleGuest = iota + 1
	UserRoleMember
	UserRoleAdmin
)

// DefaultUserRoles contains all the access levels supported by the
// app. The current access levels are:
// - Guest: Limited access to app functionality
// - Member: Regular access but no Boat/Administration functionality
// - Admin: Full access
var DefaultUserRoles = []UserRole{
	{
		ID:          UserRoleGuest,
		Name:        "guest",
		Description: "Guest User Permissions",
	},
	{
		ID:          UserRoleMember,
		Name:        "member",
		Description: "Member User Permissions",
	},
	{
		ID:          UserRoleAdmin,
		Name:        "admin",
		Description: "Administrator User Permissions",
	},
}

// TODO minimum set of user status required for the app to work

const (
	UserStatusGuest = iota + 1
	UserStatusMember
	UserStatusAdmin
	UserStatusCourse
	UserStatusPartner
)

var DefaultUserStatus = []UserStatus{
	{
		ID:          UserStatusGuest,
		Name:        "Guest",
		Description: "Guest Access",
		UserRoleID:  UserRoleGuest,
	},
	{
		ID:          UserStatusMember,
		Name:        "Member",
		Description: "Member Access",
		UserRoleID:  UserRoleMember,
	},
	{
		ID:          UserStatusAdmin,
		Name:        "Admin",
		Description: "Administrator and Boat Community Access",
		UserRoleID:  UserRoleAdmin,
	},
	{
		ID:          UserStatusCourse,
		Name:        "Course",
		Description: "Course Status and only Guest Access Permissions",
		UserRoleID:  UserRoleGuest,
	},
	{
		ID:          UserStatusPartner,
		Name:        "Partner",
		Description: "Partners with Member Permissions",
		UserRoleID:  UserRoleMember,
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

var DefaultSessionTypesMap = map[int]SessionType{
	DefaultSessionType: DefaultSessionTypes[0],
	CourseSessionType:  DefaultSessionTypes[1],
}

var DefaultPriceGuest, _ = decimal.NewFromString("2.80")
var DefaultPriceMember, _ = decimal.NewFromString("2.80")
var DefaultPriceCommunity, _ = decimal.NewFromString("1.30")
var DefaultPriceCourse, _ = decimal.NewFromString("0.00")

var DefaultPricing = []Pricing{
	{
		ID:             1,
		UserStatusID:   1,
		PricePerMinute: DefaultPriceGuest,
		Comment:        "Guest Price",
	},
	{
		ID:             2,
		UserStatusID:   2,
		PricePerMinute: DefaultPriceMember,
		Comment:        "Member Price",
	},
	{
		ID:             3,
		UserStatusID:   3,
		PricePerMinute: DefaultPriceCommunity,
		Comment:        "Boat Community Price",
	},
	{
		ID:             4,
		UserStatusID:   4,
		PricePerMinute: DefaultPriceCourse,
		Comment:        "Course Price",
	},
	{
		ID:             5,
		UserStatusID:   5,
		PricePerMinute: DefaultPriceCommunity,
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

const (
	ExpenseTypeFuelDirect = iota + 1
	ExpenseTypeMaintenance
	ExpenseTypeMaterial
	ExpenseTypeInvestment
	ExpenseTypeSession
	ExpenseTypeOther
	ExpenseTypeMembershipFee
	ExpenseTypeSalary
	ExpenseTypeOwnersRefund
	ExpenseTypeFuelBill
)

var DefaultExpenseTypes = []ExpenseType{
	{
		ID:      ExpenseTypeFuelDirect,
		Name:    "fuel",
		Comment: "Costs for fuel",
	},
	{
		ID:      ExpenseTypeMaintenance,
		Name:    "maintenance",
		Comment: "Costs for maintenance.",
	},
	{
		ID:      ExpenseTypeMaterial,
		Name:    "material",
		Comment: "Material which is needed for the boat.",
	},
	{
		ID:      ExpenseTypeInvestment,
		Name:    "invest",
		Comment: "Investment in new Features, ...",
	},
	{
		ID:      ExpenseTypeSession,
		Name:    "session",
		Comment: "Payments for sessions.",
	},
	{
		ID:      ExpenseTypeOther,
		Name:    "other",
		Comment: "Other costs.",
	},
	{
		ID:      ExpenseTypeMembershipFee,
		Name:    "membership fee",
		Comment: "Membership Fee.",
	},
	{
		ID:      ExpenseTypeSalary,
		Name:    "salary",
		Comment: "Compensation for towing people or do whatever in the name of the community.",
	},
	{
		ID:      ExpenseTypeOwnersRefund,
		Name:    "owners refund",
		Comment: "Payback of investments from owners.",
	},
	{
		ID:      ExpenseTypeFuelBill,
		Name:    "fuel bill",
		Comment: "Costs for fuel paid through invoices.",
	},
}

var DefaultConfiguration = []Configuration{
	{
		ID:       1,
		Property: "schema.version",
		Value:    "2.0",
	},
	{
		ID:       2,
		Property: "browser.session.timeout.default",
		Value:    "10800",
	},
	{
		ID:       3,
		Property: "browser.session.timeout.max",
		Value:    "604800",
	},
	{
		ID:       4,
		Property: "location.longitude",
		Value:    "8.542939",
	},
	{
		ID:       5,
		Property: "location.latitude",
		Value:    "47.367658",
	},
	{
		ID:       6,
		Property: "location.gmt_offset",
		Value:    "1",
	},
	{
		ID:       7,
		Property: "location.timezone",
		Value:    "Europe/Berlin",
	},
	{
		ID:       8,
		Property: "business.day.start",
		Value:    "08:00:00",
	},
	{
		ID:       9,
		Property: "business.day.end",
		Value:    "21:00:00",
	},
	{
		ID:       10,
		Property: "business.day.startatsunrise",
		Value:    "false",
	},
	{
		ID:       11,
		Property: "business.day.endatsunset",
		Value:    "false",
	},
	{
		ID:       12,
		Property: "session.cancel.graceperiod",
		Value:    "86400",
	},
	{
		ID:       13,
		Property: "recaptcha.privatekey",
		Value:    "",
	},
	{
		ID:       14,
		Property: "recaptcha.publickey",
		Value:    "",
	},
	{
		ID:       15,
		Property: "currency",
		Value:    "CHF",
	},
	{
		ID:       16,
		Property: "location.address",
		Value:    "",
	},
	{
		ID:       17,
		Property: "location.map",
		Value:    "",
	},
	{
		ID:       18,
		Property: "payment.account.owner",
		Value:    "",
	},
	{
		ID:       19,
		Property: "payment.account.iban",
		Value:    "",
	},
	{
		ID:       20,
		Property: "payment.account.bic",
		Value:    "",
	},
	{
		ID:       21,
		Property: "payment.account.comment",
		Value:    "",
	},
	{
		ID:       22,
		Property: "smtp.sender",
		Value:    "",
	},
	{
		ID:       23,
		Property: "smtp.server",
		Value:    "",
	},
	{
		ID:       24,
		Property: "smtp.username",
		Value:    "",
	},
	{
		ID:       25,
		Property: "smtp.password",
		Value:    "",
	},
	{
		ID:       26,
		Property: "logo.file",
		Value:    "",
	},
	{
		ID:       27,
		Property: "engine.hour.format",
		Value:    "hh.h",
	},
	{
		ID:       28,
		Property: "fuel.payment.type",
		Value:    "instant",
	},
}
