// Package model is where we keep our database tables representations as
// structures.
package model

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kberov/slovo2/util"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasttemplate"
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

const defaultLogHeader = `${prefix}:${time_rfc3339}:${level}:${short_file}:${line}`

// Logger must be instantiated before using any function from this package.
var Logger *log.Logger

// DSN must be set before using DB() function.
var DSN string

var spf = fmt.Sprintf

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

var record2Table = map[string]string{}

// Record2Table converts struct type name like *model.UsersInvoicesLastID to
// users_invoices_last_id and returns it. Caches the converted name for
// subsequent calls.
func Record2Table[T Record](record *T) string {
	typestr := spf("%T", record)
	if table, ok := record2Table[typestr]; ok {
		return table
	}
	table, _ := strings.CutPrefix(typestr, "*model.")
	table = util.CamelToSnakeCase(table)
	record2Table[typestr] = table
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
	sql := SQLFor("GetByID", Record2Table(r))
	return DB().Get(r, sql, id)
}

var global *sqlx.DB

func DB() *sqlx.DB {
	if global != nil {
		return global
	}
	Logger.Debug("database:", DSN)

	global = sqlx.MustConnect("sqlite3", DSN)
	global.MapperFunc(util.CamelToSnakeCase)
	return global
}

/*
SQLFor compooses an SQL query for the given key. Returns the composed query.
*/
func SQLFor(query, table string) string {
	q := queryTemplates[query].(string)
	queryTemplates["table"] = table
	for strings.Contains(q, "${") {
		q = fasttemplate.ExecuteStringStd(q, "${", "}", queryTemplates)
	}
	delete(queryTemplates, "table")
	return q
}
