package model

// SQLMap is a map of name/query. Each entry has a name and an SQL query used
// in some method.
type SQLMap map[string]string

var queryTemplates = SQLMap{
	"INSERT":  `INSERT INTO ${table} (${columns}) VALUES(${placeholders})`,
	"GetById": `SELECT * FROM ${table} WHERE id=?`,
	"SELECT":  `SELECT ${columns} FROM ${table} LIMIT ${limit} OFFSET ${offset}`,
}
