// Package model is where we keep our database tables representations as
// structures.
package model

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kberov/slovo2/slovo"
	"github.com/labstack/gommon/random"
	_ "github.com/mattn/go-sqlite3"
)

// import "github.com/jmoiron/sqlx"

/*
ModelI is an interface - a table. It must be implemented by all table representations. We will
generate them from the existing database.
type ModelI interface {
	Migrate() error
	Create(data Record) (*Record, error)
	All(where string, limit int, offset int) (Record, error)
	GetBy(where string) (*Record, error)
	GetById(id int64) (*Record, error)
	Update(id int64, updated Record) (*Record, error)
	Delete(id int64) error
	Data() *Record
}
*/

/*
Record is a generic constraint on the allowed types. Each type here is the
record type in a table with the same name in lowercase. Note! these types
*must* be named exactly after the tables, because we use them to guess the
respective table name.
Example: UsersInvoicesLastID === users_invoices_last_id
*/
type Record interface {
	Aliases | Celini | Domove | FirstLogin | Groups | Invoices | Orders |
		PasswLogin | Products | Stranici | UserGroup | Users | UsersInvoicesLastID
}

var spf = fmt.Sprintf

var record2Table = map[string]string{}

func Record2Table[T Record](r *T) string {
	typestr := spf("%T", r)
	if table, ok := record2Table[typestr]; ok {
		return table
	}
	// TODO: benchmark if this is faster than a regex
	table, _ := strings.CutPrefix(typestr, "*model.")
	table = strings.Replace(table, "ID", "_id", 1)
	for i, l := range random.Uppercase {
		lstr := string(l)
		if strings.Contains(table, lstr) {
			table = strings.ReplaceAll(table, lstr, "_"+string(random.Lowercase[i]))
		}
	}
	record2Table[typestr] = table[1:]
	return record2Table[typestr]
}

/*
Table is the base implementation for all tables in the database
*/
type Table struct {
	queries SQLMap
	table   string
}

func GetByID[T Record](r *T, id int32) error {
	table := fmt.Sprintf("%T", *r)
	println(table)
	return nil
}

var globalConnection *sqlx.DB

func DB() *sqlx.DB {
	if globalConnection != nil {
		return globalConnection
	}
	globalConnection = sqlx.MustConnect("sqlite3", slovo.Cfg.DB.DSN)
	return globalConnection
}

type TableInfo struct {
	TableColumns []TableColumn
}

type TableColumn struct {
	CID       int    `db:cid`
	Name      string `db:name`
	Type      string `db:type`
	NotNull   uint8  `db:notnull`
	DfltValue string `db:dflt_value`
}
