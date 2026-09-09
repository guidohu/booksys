package database

import (
	"reflect"
	"time"

	"github.com/shopspring/decimal"
)

// BoatEngineHour is a reading of the engine hour meter, taken before and
// after a session.
type BoatEngineHour struct {
	ID          uint            `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Timestamp   time.Time       `gorm:"type:datetime DEFAULT NULL"`
	BeforeHours decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	AfterHours  decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	DeltaHours  decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	TypeID      uint8           `gorm:"column:type;type:mediumint(9) DEFAULT NULL"`
	Type        SessionType     `gorm:"foreignKey:TypeID;references:ID"`
	UserID      uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User        User            `gorm:"foreignKey:UserID;references:ID"`
	Comment     string          `gorm:"type:text CHARACTER SET utf8"`
	CheckedIn   bool            `gorm:"column:checked_in;type:int(8) DEFAULT 0"`
}

// BoatFuel is a refuelling of the boat. ContributeToBalance says whether the
// person who paid gets credited for it.
type BoatFuel struct {
	ID                  uint             `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Timestamp           time.Time        `gorm:"type:datetime DEFAULT NULL"`
	UserID              uint             `gorm:"type:mediumint(9) DEFAULT NULL"`
	User                User             `gorm:"foreignKey:UserID;references:ID"`
	EngineHours         *decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	Liters              *decimal.Decimal `gorm:"type:DECIMAL(10,3) DEFAULT NULL"`
	Cost                *decimal.Decimal `gorm:"column:cost_chf;type:DECIMAL(10,3) DEFAULT NULL"`
	CostBrutto          *decimal.Decimal `gorm:"column:cost_chf_brutto;type:DECIMAL(10,3) DEFAULT NULL"`
	ContributeToBalance bool             `gorm:"column:contributes_to_balance;type:int(8) DEFAULT 1"`
	IsDiscounted        bool             `gorm:"column:is_discounted;type:int(8) DEFAULT 0"`
}

// TableName returns the name of the table BoatFuel maps to.
func (BoatFuel) TableName() string {
	return "boat_fuel"
}

// BoatMaintenance is a maintenance job carried out on the boat.
type BoatMaintenance struct {
	ID          uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	Timestamp   time.Time       `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP"`
	UserID      uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User        User            `gorm:"foreignKey:UserID;references:ID"`
	EngineHours decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	Description string          `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table BoatMaintenance maps to.
func (BoatMaintenance) TableName() string {
	return "boat_maintenance"
}

// BrowserSession is a login session. It denormalizes the parts of the user
// that the middleware needs so that authentication costs a single lookup.
type BrowserSession struct {
	SessionSecret string       `gorm:"primaryKey;type:varchar(512) NOT NULL"`
	ValidUntil    time.Time    `gorm:"column:valid_thru;type:datetime DEFAULT NULL"`
	MaxValidUntil time.Time    `gorm:"column:max_valid_thru;type:datetime DEFAULT NULL"`
	LastActivity  time.Time    `gorm:"type:datetime DEFAULT NULL"`
	UserID        uint         `gorm:"type:mediumint(9) DEFAULT NULL"`
	User          User         `gorm:"foreignKey:UserID;references:ID"`
	Username      string       `gorm:"type:varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL"`
	FirstName     string       `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	LastName      string       `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	UserStatus    uint         `gorm:"type:int(11) DEFAULT '0'"`
	UserRoleID    UserRoleType `gorm:"type:int(11) NOT NULL"`
	SessionData   string       `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table BrowserSession maps to.
func (BrowserSession) TableName() string {
	return "browser_session"
}

// IsEmpty returns whether this object is empty.
func (b *BrowserSession) IsEmpty() bool {
	emptySession := BrowserSession{}
	return reflect.DeepEqual(*b, emptySession)
}

// Expired returns whether the session is still valid or expired.
func (b *BrowserSession) Expired() bool {
	return b.ValidUntil.Before(time.Now())
}

// Valid returns true if a session is neither expired nor empty.
func (b *BrowserSession) Valid() bool {
	return !b.Expired() && !b.IsEmpty()
}

// Expense is money that left the account.
type Expense struct {
	ID            uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserID        uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User          User            `gorm:"foreignKey:UserID;references:ID"`
	ExpenseTypeID uint            `gorm:"column:type_id;type:int(11) NOT NULL DEFAULT '0'"`
	ExpenseType   ExpenseType     `gorm:"foreignKey:ExpenseTypeID;references:ID;constraint:OnUpdate:CASCADE"`
	Timestamp     time.Time       `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Amount        decimal.Decimal `gorm:"column:amount_chf;type:DECIMAL(10,3) DEFAULT NULL"`
	Comment       string          `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table Expense maps to.
func (Expense) TableName() string {
	return "expenditure"
}

// ExpenseType categorizes an Expense or an Income.
type ExpenseType struct {
	ID      uint   `gorm:"type:int(11) NOT NULL AUTO_INCREMENT" json:"id"`
	Name    string `gorm:"type:text CHARACTER SET utf8" json:"name"`
	Comment string `gorm:"type:text CHARACTER SET utf8" json:"comment"`
}

// TableName returns the name of the table ExpenseType maps to.
func (ExpenseType) TableName() string {
	return "expenditure_type"
}

// Heat is a single ride, the unit users are billed for.
type Heat struct {
	ID              uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserID          uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User            User            `gorm:"foreignKey:UserID;references:ID"`
	SessionID       uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	Session         Session         `gorm:"foreignKey:SessionID;references:ID"`
	Timestamp       time.Time       `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	DurationSeconds int32           `gorm:"column:duration_s;type:int(11) DEFAULT NULL"`
	Cost            decimal.Decimal `gorm:"column:cost_chf;type:DECIMAL(10,3) DEFAULT NULL"`
	Comment         string          `gorm:"type:text CHARACTER SET utf8 DEFAULT NULL"`
}

// TableName returns the name of the table Heat maps to.
func (Heat) TableName() string {
	return "heat"
}

// Invitation is an invite of one user to another user's session.
type Invitation struct {
	ID        uint             `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	SessionID uint             `gorm:"type:mediumint(9) DEFAULT NULL"`
	Session   Session          `gorm:"foreignKey:SessionID;references:ID"`
	InviteeID uint             `gorm:"column:user_id;type:mediumint(9) DEFAULT NULL"`
	Invitee   User             `gorm:"foreignKey:InviteeID;references:ID"`
	InviterID uint             `gorm:"type:mediumint(9) DEFAULT NULL"`
	Inviter   User             `gorm:"foreignKey:InviterID;references:ID"`
	Timestamp time.Time        `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	StatusID  uint             `gorm:"column:status;type:int(11) DEFAULT NULL"`
	Status    InvitationStatus `gorm:"foreignKey:StatusID;references:ID"`
}

// TableName returns the name of the table Invitation maps to.
func (Invitation) TableName() string {
	return "invitation"
}

// InvitationStatus is the state an Invitation is in.
type InvitationStatus struct {
	ID      uint   `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name    string `gorm:"type:text CHARACTER SET utf8"`
	Comment string `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table InvitationStatus maps to.
func (InvitationStatus) TableName() string {
	return "invitation_status"
}

// PasswordReset is a password reset token issued to a user.
// A user has at most one live token: issuing a new one invalidates the
// previous ones. IssuedAt drives the cooldown between requests, ValidUntil the
// expiry, and Attempts counts wrong guesses so that a token can be burned long
// before its keyspace could be searched.
//
// Note that the superseded `timestamp` column is deliberately not mapped here.
// It carries ON UPDATE CURRENT_TIMESTAMP, so every write to a row rewrote it,
// which makes it unusable both as an issue time and as an expiry.
type PasswordReset struct {
	ID         uint
	UserID     uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	Token      string    `gorm:"type:varchar(255) NOT NULL"`
	IssuedAt   time.Time `gorm:"type:datetime DEFAULT NULL"`
	ValidUntil time.Time `gorm:"type:datetime DEFAULT NULL"`
	Attempts   uint      `gorm:"type:int(11) NOT NULL DEFAULT 0"`
	Valid      bool      `gorm:"type:tinyint(1) DEFAULT '0'"`
}

// TableName returns the name of the table PasswordReset maps to.
func (PasswordReset) TableName() string {
	return "password_reset"
}

// Income is money that entered the account. It shares the type table with
// Expense.
type Income struct {
	ID            uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserID        uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User          User            `gorm:"foreignKey:UserID;references:ID"`
	Timestamp     time.Time       `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Amount        decimal.Decimal `gorm:"column:amount_chf;type:DECIMAL(10,3) DEFAULT NULL"`
	ExpenseTypeID uint            `gorm:"column:type_id;type:int(11) NOT NULL DEFAULT '0'"`
	ExpenseType   ExpenseType     `gorm:"foreignKey:ExpenseTypeID;references:ID;constraint:OnUpdate:CASCADE"`
	Comment       string          `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table Income maps to.
func (Income) TableName() string {
	return "payment"
}

// Pricing is the price per minute charged to a user group.
type Pricing struct {
	ID             uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserStatusID   uint            `gorm:"type:int(11) DEFAULT NULL"`
	UserStatus     UserStatus      `gorm:"foreignKey:UserStatusID;references:ID"`
	PricePerMinute decimal.Decimal `gorm:"column:price_chf_min;type:DECIMAL(10,3) DEFAULT NULL"`
	Comment        string          `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table Pricing maps to.
func (Pricing) TableName() string {
	return "pricing"
}

// Session is a bookable slot on the water.
type Session struct {
	ID            uint        `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	StartTime     time.Time   `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP"`
	EndTime       time.Time   `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP"`
	Title         string      `gorm:"type:text CHARACTER SET utf8"`
	Comment       string      `gorm:"type:text CHARACTER SET utf8"`
	SessionTypeID uint        `gorm:"column:type;type:mediumint(9) DEFAULT '0'"`
	SessionType   SessionType `gorm:"foreignKey:SessionTypeID;references:ID"`
	FreeSpaces    uint        `gorm:"column:free;type:int(11) DEFAULT NULL"`
	CreatorID     uint        `gorm:"type:mediumint(9) DEFAULT NULL"`
	Creator       User        `gorm:"foreignKey:CreatorID;references:ID"`
}

// TableName returns the name of the table Session maps to.
func (Session) TableName() string {
	return "session"
}

// Equal reports whether two sessions are identical in every field.
func (s Session) Equal(t Session) bool {
	return reflect.DeepEqual(s, t)
}

// SessionType categorizes a Session.
type SessionType struct {
	ID      uint   `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	Name    string `gorm:"type:text CHARACTER SET utf8"`
	Comment string `gorm:"type:text CHARACTER SET utf8"`
}

// TableName returns the name of the table SessionType maps to.
func (SessionType) TableName() string {
	return "session_type"
}

// User is a person using the application. Deleting a user clears the personal
// fields and sets IsDeleted rather than removing the row, so that the
// transactions that reference it stay intact.
type User struct {
	ID            uint       `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	Username      string     `gorm:"type:varchar(255) DEFAULT NULL"`
	PasswordSalt  int        `gorm:"type:int(11) DEFAULT NULL"`
	PasswordHash  string     `gorm:"type:text CHARACTER SET utf8"`
	FirstName     string     `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	LastName      string     `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	Address       string     `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	City          string     `gorm:"type:varchar(255) CHARACTER SET utf8 DEFAULT NULL"`
	ZipCode       int        `gorm:"column:plz;type:int(11) DEFAULT NULL"`
	MobilePhoneNr string     `gorm:"column:mobile;type:varchar(255) DEFAULT NULL"`
	Email         string     `gorm:"type:varchar(255) DEFAULT NULL"`
	BoatLicense   bool       `gorm:"column:license"`
	UserStatusID  uint       `gorm:"column:status;type:int(11) DEFAULT '0'"`
	UserStatus    UserStatus `gorm:"foreignKey:UserStatusID;references:ID"`
	Locked        bool       `gorm:"type:tinyint(1) DEFAULT '1'"`
	Comment       string     `gorm:"type:text CHARACTER SET utf8"`
	IsDeleted     bool       `gorm:"column:deleted;type:tinyint(1) DEFAULT '0'"`
}

// TableName returns the name of the table User maps to.
func (User) TableName() string {
	return "user"
}

// UserRole is an access level.
type UserRole struct {
	ID          UserRoleType `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name        string       `gorm:"type:text COLLATE utf8_bin NOT NULL"`
	Description string       `gorm:"type:text COLLATE utf8_bin NOT NULL"`
}

// TableName returns the name of the table UserRole maps to.
func (UserRole) TableName() string {
	return "user_role"
}

// UserStatus is a user group. It ties a set of users to a UserRole and,
// through Pricing, to a price.
type UserStatus struct {
	ID          uint         `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name        string       `gorm:"type:text CHARACTER SET utf8"`
	Description string       `gorm:"type:text CHARACTER SET utf8"`
	UserRoleID  UserRoleType `gorm:"type:int(11) NOT NULL"`
	UserRole    UserRole     `gorm:"foreignKey:UserRoleID;references:ID"`
}

// TableName returns the name of the table UserStatus maps to.
func (UserStatus) TableName() string {
	return "user_status"
}

// UserToSession records that a user takes part in a session.
type UserToSession struct {
	ID        uint      `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserID    uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	SessionID uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	Session   Session   `gorm:"foreignKey:SessionID;references:ID"`
	TimeAdded time.Time `gorm:"column:time;type:datetime DEFAULT NULL"`
}

// TableName returns the name of the table UserToSession maps to.
func (UserToSession) TableName() string {
	return "user_to_session"
}

// Configuration is a single configuration property stored in the database.
type Configuration struct {
	// ID       uint   `gorm:"type:int(11) DEFAULT 0"`
	Property string `gorm:"primaryKey;type:VARCHAR(255) NOT NULL"`
	Value    string `gorm:"type:TEXT"`
}

// TableName returns the name of the table Configuration maps to.
func (Configuration) TableName() string {
	return "configuration"
}

// ConfigurationVersion is bumped whenever a property changes, so that readers
// can detect updates without re-reading every property.
type ConfigurationVersion struct {
	ID        uint       `gorm:"type:int(11);autoIncrement:false"`
	Version   uint       `gorm:"type:int(11) NOT NULL DEFAULT 1"`
	Timestamp *time.Time `gorm:"column:time;type:datetime DEFAULT NULL"`
}

// TableName returns the name of the table ConfigurationVersion maps to.
func (ConfigurationVersion) TableName() string {
	return "configuration_version"
}
