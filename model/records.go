package model

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// In this file we store records by table name. Each type represents a row in
// the respective table after which it is named.

// Domove is a record from table domove.
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

func (d *Domove) GetByName(domain string) error {
	table := Record2Table(d)
	SQL := SQLFor("GET_DOMAIN", table)
	// Logger.Debugf("domain: %#v GetByName(GET_DOMAIN) SQL:\n%s", domain, SQL)
	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		return err
	} else {
		args := struct{ Domain string }{Domain: domain}
		return stmt.Get(d, args)
	}
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
	// Alias for an article/paragraph/book/product/content. Example: "чуждият-hôtel"
	Celina string `param:"celina"`
	// Box on the page where to be put  the body of a celina
	Box Box
	// Language for the content (for now only "bg")
	Lang string `param:"lang"`
	// Format of the content (for now only "html")
	Format string `param:"format"`
	// Required Published status (0:no|1:preview|2:yes) Default: 2
	Pub int
	// UserID - current user_id. Default 2(guest)
	// TODO: Implement authentication via some custom header or JWT - not Cookies
	UserID int32
	// Slovo is a multidomain CMS. Get it from c.Request().Host
	Domain string
	Now    int64
	Limit  int `query:"limit" form:"limit"`
	Offset int `query:"offset" form:"offset"`
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
	Published   uint8
	Hidden      bool
	Deleted     bool
	ChangedBy   string
	/*
		Here are fields from Celini. We may get them with some Get, when we
		select from both tables. When selecting only from stranici table, these
		will be empty. When selecting default content for a page these are
		populated from the celina with the same `alias` field as the page. This
		celina is created when the page is created and holds all the other
		celini in this page.
	*/
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
func (s *Stranici) FindForDisplay(args StraniciArgs) (err error) {
	table := Record2Table(s)
	args.Lang = args.Lang + `%`
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)
	// Logger.Debugf("args: %#v FindForDisplay(GET_PAGE_FOR_DISPLAY) SQL:\n%s", args, SQL)
	var stmt *sqlx.NamedStmt
	if stmt, err = DB().PrepareNamed(SQL); err != nil {
		return err
	}
	return stmt.Get(s, args)
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

// ListStranici returns a slice of pages which are children of the page with
// StraniciArgs.Alias.
func ListStranici(args StraniciArgs) (items []Stranici) {

	SQL := SQLFor("SELECT_CHILD_PAGES", "stranici")
	Logger.Debugf("ListStranici(%#v) SQL:\n%s", args, SQL)

	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		Logger.Panicf(`error from ListStranici/PrepareNamed(SQL):%v; args: %#v`, err, args)
	} else if err = stmt.Select(&items, args); err != nil {
		Logger.Panicf(`error from ListStranici/Select(&items, args):%v; args: %#v`, err, args)
	}
	return
}

type Box string

const (
	MainBox   Box = "main"
	LeftBox   Box = "left"
	RightBox  Box = "right"
	HeaderBox Box = "header"
	FooterBox Box = "footer"
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
func SelectMenuItems(args StraniciArgs) (items []StrMenuItem) {
	SQL := SQLFor("SELECT_PAGES_FOR_MAIN_MENU", "stranici")
	// Logger.Debugf("SelectMenuItems(%#v) SQL:\n%s", args, SQL)
	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		Logger.Panicf(`error from SelectMenuItems/PrepareNamed(SQL):%v; args: %#v`, err, args)
	} else if err = stmt.Select(&items, args); err != nil {
		Logger.Panicf(`error from SelectMenuItems/Select(&items, args):%v; args: %#v`, err, args)
	}
	return items
}

// IsDir returns true if the permissions field starts with `d`.
func (s *StrMenuItem) IsDir() bool { return strings.HasPrefix(s.Permissions, "d") }

// CelDataType represents the logical data type for a celina. Based on it we
// choose the appropriate template for displaying it.
// "default": "note",
// "enum":["title", "note", "writing", "book", "chapter", "question", "answer", "paragraph"]
type CelDataType string

const (
	Title     CelDataType = "title"
	Note      CelDataType = "note"
	Writing   CelDataType = "writing"
	Book      CelDataType = "book"
	Chapter   CelDataType = "chapter"
	Question  CelDataType = "question"
	Answer    CelDataType = "Answer"
	Paragraph             = "paragraph"
)

func (d CelDataType) String() string {
	return string(d)
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
	DataType    CelDataType
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
	Bad         bool
	Deleted     bool
	Start       int32
	Stop        int32
	ChangedBy   sql.NullString
	Published   uint8
}

func (cel *Celini) FindForDisplay(args StraniciArgs) error {
	SQL := SQLFor("GET_CELINA_FOR_DISPLAY", Record2Table(cel))
	args.Lang = args.Lang + `%`
	//Logger.Debugf("GET_CELINA_FOR_DISPLAY SQL:\n%s", SQL)
	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		return err
	} else {
		return stmt.Get(cel, args)
	}
}

var celiniTemplatesDir = "celini/"

func (cel *Celini) TemplatePath(defaultTemplate string) string {

	switch cel.DataType {
	case Writing, Book, Note:
		return celiniTemplatesDir + cel.DataType.String()
	default:
		// TODO:define templates and behaviour for all data types
		return celiniTemplatesDir + `writing`
	}
}

// ListCelini returns a slice of celini which are children of the page with
// StraniciArgs.Alias. The celina used for title of the page in the respective
// language is `pid` for these celini.
func ListCelini(args StraniciArgs) (items []Celini) {
	SQL := SQLFor("CELINI_FOR_LIST_IN_PAGE", "celini")
	Logger.Debugf("ListCelini(%#v) SQL:\n%s", args, SQL)

	if stmt, err := DB().PrepareNamed(SQL); err != nil {
		Logger.Panicf("error from model.ListStranici/PrepareNamed(SQL): %v", err)
	} else if err = stmt.Select(&items, args); err != nil {
		Logger.Panicf("error from model.ListStranici/Select(&items, args): %v", err)
	}
	return items
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
