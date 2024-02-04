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

type PageType string

const (
	Regular PageType = "regular"
	Root    PageType = "root"
)

// StraniciArgs is a struct which we will populate by using the Echo binding mechanizm
// to pass these arguments to SQL queries and whereever needed. For binding
// these struct fields to values from different places, we use our custom
// [slovo.Binder].
type StraniciArgs struct {
	// Alias for a page: "вѣра"
	Alias string `param:"stranica"`
	// Alias for an article/paragraph/book/product/content: "чуждият-hôtel"
	Celina string `param:"celina"`
	// Language for the content (for now only "bg")
	Lang string `param:"lang"`
	// Format of the content (for now only "html")
	Format string `param:"format"`
	// Required Published status (0:no|1:preview|2:yes) Default: 2
	Pub uint8
	// UserID - current user_id. Default 2(guest)
	// TODO: Implement authentication via some custom header or JWT - not Cookies
	UserID int32
	// Slovo is a multidomain CMS. Get it from c.Request().Host
	Domain string
	Now    int64
}

// Stranici represents a Record in table stranici.
type Stranici struct {
	ID          int32
	Pid         int32
	DomID       int32
	Alias       string
	PageType    PageType
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

/*
FindForDisplay returns a page from the database to be displayed. The page
must have the given alias, readable by the given user, be in the given
domain  and published(=2). `args` is a struct containing the arguments for
stmt.Get. It is put together in slovo.Binder.Bind().
*/
func (s *Stranici) FindForDisplay(args *StraniciArgs) error {
	table := Record2Table(s)
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)
	// Logger.Debugf("FindForDisplay(GET_PAGE_FOR_DISPLAY) SQL:\n%s", SQL)

	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		return err
	} else {
		return stmt.Get(s, args)
	}
}

// IsDir returns true if the permissions field starts with `d`.
func (s *Stranici) IsDir() bool { return strings.HasPrefix(s.Permissions, "d") }

/*
TemplatePath returns the path to the template file to be used for this page
record or the given `defaultTemplate`.
*/
func (s *Stranici) TemplatePath(defaultTemplate string) string {
	if s.Template.Valid {
		return s.Template.String
	}
	return defaultTemplate
}

type Box string

const (
	Main   Box = "main"
	Left   Box = "left"
	Right  Box = "right"
	Header Box = "header"
	Footer Box = "footer"
)

/*
StrMenuItem holds subset of the fields of a page and some fields from Celini
for pdisplaying a link in the main menu. The query to populate this struct is
SELECT_PAGES_FOR_MAIN_MENU.
*/
type StrMenuItem struct {
	ID          int32
	Pid         int32
	Alias       string
	Title       string
	Language    string
	Permissions string
}

/*
SelectMenuItems populates a []StrMenuItem slice and returns it or an
error from DB().
*/
func SelectMenuItems(args *StraniciArgs) (items []StrMenuItem, err error) {
	SQL := SQLFor("SELECT_PAGES_FOR_MAIN_MENU", "stranici")
	//Logger.Debugf("SelectMenuItems(%#v) SQL:\n%s", args, SQL)
	stmt, err := DB().PrepareNamed(SQL)
	if err != nil {
		return nil, err
	}
	err = stmt.Select(&items, args)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// IsDir returns true if the permissions field starts with `d`.
func (s *StrMenuItem) IsDir() bool { return strings.HasPrefix(s.Permissions, "d") }

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
	Box         Box
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
	args := map[string]any{
		"alias":   alias,
		"box":     box,
		"lang":    language + `%`,
		"now":     time.Now().Unix(),
		"page_id": page.ID,
		"user_id": user.ID,
		"pub":     preview,
	}
	// Logger.Debugf("GET_CELINA_FOR_DISPLAY SQL:\n%s", SQL)
	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		return err
	} else {
		return stmt.Get(ce, args)
	}
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
