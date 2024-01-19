// Package model is where we keep our database tables representations as
// structures.
package model

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
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

// Logger must be instantiated before using any function from this package
var Logger *log.Logger

// DSN must be set before using DB() funstion
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

	table = strings.Replace(table, "ID", "_id", 1)
	var snakeCase strings.Builder
	for _, r := range table {
		if unicode.IsUpper(r) {
			snakeCase.Write([]byte{'_', byte(unicode.ToLower(r))})
			continue
		}
		snakeCase.WriteRune(r)
	}
	// remove the prefixed underscore for the first uppercase letter
	record2Table[typestr] = snakeCase.String()[1:]
	return record2Table[typestr]
}

var field2Column = map[string]string{}

// Field2Column converts ColumName to column_name and returns it. Works only
// with latin letters. Caches the converted field for subsequent calls.
func Field2Column(field string) string {
	if field == "ID" {
		return "id"
	}
	if f, ok := field2Column[field]; ok {
		return f
	}
	field = strings.Replace(field, "ID", "_id", -1)
	var snakeCase strings.Builder
	for _, r := range field {
		if unicode.IsUpper(r) {
			snakeCase.Write([]byte{'_', byte(unicode.ToLower(r))})
			continue
		}
		snakeCase.WriteRune(r)
	}
	// remove the prefixed underscore for the first uppercase letter
	field2Column[field] = snakeCase.String()[1:]
	return field2Column[field]
}

/*
Table is the base implementation for all tables in the database
*/
type Table struct {
	queries SQLMap
	table   string
}

func GetByID[T Record](r *T, id int32) error {
	table := Record2Table(r)
	println(table)
	return nil
}

var globalConnection *sqlx.DB

func DB() *sqlx.DB {
	if globalConnection != nil {
		return globalConnection
	}
	Logger.Debug("database:", DSN)

	globalConnection = sqlx.MustConnect("sqlite3", DSN)
	globalConnection.MapperFunc(Field2Column)
	return globalConnection
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
