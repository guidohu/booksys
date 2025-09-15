package database

import (
	"reflect"
	"time"

	"github.com/shopspring/decimal"
)

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

func (BoatFuel) TableName() string {
	return "boat_fuel"
}

type BoatMaintenance struct {
	ID          uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	Timestamp   time.Time       `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP"`
	UserID      uint            `gorm:"type:mediumint(9) DEFAULT NULL"`
	User        User            `gorm:"foreignKey:UserID;references:ID"`
	EngineHours decimal.Decimal `gorm:"type:DECIMAL(10,5) DEFAULT NULL"`
	Description string          `gorm:"type:text CHARACTER SET utf8"`
}

func (BoatMaintenance) TableName() string {
	return "boat_maintenance"
}

type BrowserSession struct {
	SessionSecret string       `gorm:"primaryKey;type:varchar(512) NOT NULL"`
	ValidUntil    time.Time    `gorm:"column:valid_thru;type:datetime DEFAULT NULL"`
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

func (BrowserSession) TableName() string {
	return "browser_session"
}

// IsEmpty returns whether this oobject is empty.
func (b *BrowserSession) IsEmpty() bool {
	emptySession := BrowserSession{}
	return reflect.DeepEqual(*b, emptySession)
}

// Expired returns whether the session is still valid or expired.
func (b *BrowserSession) Expired() bool {
	return b.ValidUntil.Before(time.Now())
}

// Valid returns true if a session is both, not empty and
func (b *BrowserSession) Valid() bool {
	return !b.Expired() && !b.IsEmpty()
}

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

func (Expense) TableName() string {
	return "expenditure"
}

type ExpenseType struct {
	ID      uint   `gorm:"type:int(11) NOT NULL AUTO_INCREMENT" json:"id"`
	Name    string `gorm:"type:text CHARACTER SET utf8" json:"name"`
	Comment string `gorm:"type:text CHARACTER SET utf8" json:"comment"`
}

func (ExpenseType) TableName() string {
	return "expenditure_type"
}

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

func (Heat) TableName() string {
	return "heat"
}

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

func (Invitation) TableName() string {
	return "invitation"
}

type InvitationStatus struct {
	ID      uint   `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name    string `gorm:"type:text CHARACTER SET utf8"`
	Comment string `gorm:"type:text CHARACTER SET utf8"`
}

func (InvitationStatus) TableName() string {
	return "invitation_status"
}

type PasswordReset struct {
	ID        uint
	UserID    uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Token     string    `gorm:"type:varchar(255) NOT NULL"`
	Timestamp time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Valid     bool      `gorm:"type:tinyint(1) DEFAULT '0'"`
}

func (PasswordReset) TableName() string {
	return "password_reset"
}

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

func (Income) TableName() string {
	return "payment"
}

type Pricing struct {
	ID             uint            `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserStatusID   uint            `gorm:"type:int(11) DEFAULT NULL"`
	UserStatus     UserStatus      `gorm:"foreignKey:UserStatusID;references:ID"`
	PricePerMinute decimal.Decimal `gorm:"column:price_chf_min;type:DECIMAL(10,3) DEFAULT NULL"`
	Comment        string          `gorm:"type:text CHARACTER SET utf8"`
}

func (Pricing) TableName() string {
	return "pricing"
}

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

func (Session) TableName() string {
	return "session"
}

func (s Session) Equal(t Session) bool {
	return reflect.DeepEqual(s, t)
}

type SessionType struct {
	ID      uint   `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	Name    string `gorm:"type:text CHARACTER SET utf8"`
	Comment string `gorm:"type:text CHARACTER SET utf8"`
}

func (SessionType) TableName() string {
	return "session_type"
}

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

func (User) TableName() string {
	return "user"
}

type UserRole struct {
	ID          UserRoleType `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name        string       `gorm:"type:text COLLATE utf8_bin NOT NULL"`
	Description string       `gorm:"type:text COLLATE utf8_bin NOT NULL"`
}

func (UserRole) TableName() string {
	return "user_role"
}

type UserStatus struct {
	ID          uint         `gorm:"type:int(11) NOT NULL AUTO_INCREMENT"`
	Name        string       `gorm:"type:text CHARACTER SET utf8"`
	Description string       `gorm:"type:text CHARACTER SET utf8"`
	UserRoleID  UserRoleType `gorm:"type:int(11) NOT NULL"`
	UserRole    UserRole     `gorm:"foreignKey:UserRoleID;references:ID"`
}

func (UserStatus) TableName() string {
	return "user_status"
}

type UserToSession struct {
	ID        uint      `gorm:"type:mediumint(9) NOT NULL AUTO_INCREMENT"`
	UserID    uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	SessionID uint      `gorm:"type:mediumint(9) DEFAULT NULL"`
	Session   Session   `gorm:"foreignKey:SessionID;references:ID"`
	TimeAdded time.Time `gorm:"column:time;type:datetime DEFAULT NULL"`
}

func (UserToSession) TableName() string {
	return "user_to_session"
}

type Configuration struct {
	// ID       uint   `gorm:"type:int(11) DEFAULT 0"`
	Property string `gorm:"primaryKey;type:VARCHAR(255) NOT NULL"`
	Value    string `gorm:"type:TEXT"`
}

func (Configuration) TableName() string {
	return "configuration"
}

type ConfigurationVersion struct {
	ID        uint       `gorm:"type:int(11);autoIncrement:false"`
	Version   uint       `gorm:"type:int(11) NOT NULL DEFAULT 1"`
	Timestamp *time.Time `gorm:"column:time;type:datetime DEFAULT NULL"`
}

func (ConfigurationVersion) TableName() string {
	return "configuration_version"
}
