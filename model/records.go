package model

import (
	"database/sql"
	"strings"
	"time"
)

// Domove is a records from table domove.
// In this file we strore records by table name. Each type represents a row in
// the respective table after which it is named.
type Domove struct {
	ID          int32
	Domain      string
	SiteName    string
	Description string
	OwnerID     int32
	GroupID     int32
	Permissions string
	Published   int32
	Ips         string
	Aliases     string
	Templates   string
}

type Stranici struct {
	ID          int32
	Pid         int32
	DomID       int32
	Alias       string
	PageType    string
	Permissions string
	Sorting     int32
	Template    sql.NullString
	UserID      int32
	GroupID     int32
	Tstamp      int32
	Start       int32
	Stop        int32
	Published   int32
	Hidden      int32
	Deleted     int32
	ChangedBy   string
	// Here are fields from Celini. We may get them with some Get, when we
	// select from both tables. When selecting only from stranici table, these
	// will be empty. These are populated from the celina which is created when
	// the page is created and holds all the other celini in this page.
	Title       string
	Description string
	Keywords    string
	Tags        string
	Language    string
	Body        string
	DataType    string
	DataFormat  string
}

// FindForDisplay returns a page from the database to be displayed. The page
// must have the given alias, readable by the given user, be in the given
// domain  and published(=2).
func (s *Stranici) FindForDisplay(alias string, user *Users, preview uint8, domain string, lang string) error {
	table := Record2Table(s)
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)
	now := time.Now().Unix()
	//Logger.Debugf("FindForDisplay(alias:%s,user.ID:%d, domain:%s) SQL:\n%s", alias, user.ID, domain, SQL)
	return DB().Get(s, SQL, lang[:2], user.ID, user.ID, preview, alias, alias, alias, alias, domain, domain, domain, now, now)
}

// IsDir returns true if the permissions field starts with `d`.
func (s *Stranici) IsDir() bool { return strings.HasPrefix(s.Permissions, "d") }

/*
TemplatePath returns the path to the template file to be used for this page
record.
*/
func (s *Stranici) TemplatePath() string {
	return s.Template.String
}

/*
HasTemplate tells if this page has a custom template to be used for this
page record.
*/
func (s *Stranici) HasTemplate() bool {
	return len(s.Template.String) > 0
}

type Celini struct {
	ID          int32
	Alias       string
	Pid         int32
	FromID      int32
	PageID      int32
	UserID      int32
	GroupID     int32
	Sorting     int32
	DataType    string
	DataFormat  string
	CreatedAt   int32
	Tstamp      int32
	Title       string
	Description string
	Keywords    string
	Tags        string
	Body        string
	Box         string
	Language    string
	Permissions string
	Featured    int32
	Accepted    int32
	Bad         int32
	Deleted     int32
	Start       int32
	Stop        int32
	ChangedBy   string
	Published   int32
}

func (ce *Celini) FindForDisplay(page *Stranici, alias string, user *Users, preview uint8, language string, box string) error {
	SQL := SQLFor("GET_CELINA_FOR_DISPLAY", Record2Table(ce))
	now := time.Now().Unix()
	// Logger.Debugf("GET_CELINA_FOR_DISPLAY(page.ID:%#v, user.ID:%d, preview:%v, alias:%s,  language: %s, now:%d) SQL:\n%s",
	//	page.ID, user.ID, preview, alias, language, now, SQL)
	return DB().Get(ce, SQL, page.ID, language+`%`, box, user.ID, user.ID, preview, alias, alias, alias, alias, now, now)
}

type Aliases struct {
	ID         int32
	OldAlias   string
	NewAlias   string
	AliasID    string
	AliasTable string
}

type Users struct {
	ID            int32
	LoginName     string
	LoginPassword string
	FirstName     string
	LastName      string
	Email         string
	Description   string
	Tstamp        int32
	RegTime       int32
	Disabled      int32
	StartDate     int32
	StopDate      int32
	GroupID       int32
	CreatedBy     int32
	ChangedBy     int32
}

type Groups struct {
	ID          int32
	Name        string
	Description string
	Disabled    int32
	CreatedBy   int32
	ChangedBy   int32
}

type UserGroup struct {
	UserID  int32
	GroupID int32
}

type FirstLogin struct {
	ID        int32
	Token     string
	FromUID   int32
	ToUID     int32
	StartDate int32
	StopDate  int32
}

type PasswLogin struct {
	ID        int32
	Token     string
	ToUID     int32
	StartDate int32
	StopDate  int32
}

type Products struct {
	ID         int32
	Sku        string
	Alias      string
	Title      string
	PType      string
	Properties string
}

type Orders struct {
	ID          int32
	Name        string
	Email       string
	Phone       string
	Deliverer   string
	DelivererID string
	CityName    string
	Poruchka    string
	WayBill     string
	CreatedAt   int32
	Tstamp      int32
	WayBillID   string
	Executed    int32
}
type Invoices struct {
	ID            int32
	UserInvoiceID int32
	UserID        int32
	OrderID       int32
	ChangedBy     int32
}

type UsersInvoicesLastID struct {
	UserID    int32
	InvoiceID int32
}
