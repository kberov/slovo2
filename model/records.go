package model

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
	Template    string
	UserID      int32
	GroupID     int32
	Tstamp      int32
	Start       int32
	Stop        int32
	Published   int32
	Hidden      int32
	Deleted     int32
	ChangedBy   string
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
