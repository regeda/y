package y

import (
	"database/sql"
	"log"

	sq "github.com/lann/squirrel"
)

// Execer decribes exec operation
type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

// Queryer decribes query operation
type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// DB describes db operations
type DB interface {
	Execer
	Queryer
}

// Qualifier updates a select builder if you need
type Qualifier func(sq.SelectBuilder) sq.SelectBuilder

// NoopQualifier returns the select builder "as is"
var NoopQualifier = func(b sq.SelectBuilder) sq.SelectBuilder {
	return b
}

// ByEq returns the filter by squirrel.Eq
var ByEq = func(eq sq.Eq) Qualifier {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(eq)
	}
}

// ByID returns the filter by ID
var ByID = func(id int64) Qualifier {
	return ByEq(sq.Eq{"id": id})
}

// TxPipe run some db statement
type TxPipe func(db DB, v interface{}) error

// Tx executes statements in a transaction
func Tx(db *sql.DB, v interface{}, pipes ...TxPipe) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	for _, pipe := range pipes {
		err = pipe(tx, v)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func sqlize(q sq.Sqlizer) (sql string, args []interface{}) {
	sql, args, _ = q.ToSql()
	if Debug {
		log.Printf("y/db: SQL: %s, args: %#v", sql, args)
	}
	return
}

func exec(q sq.Sqlizer, db DB) (sql.Result, error) {
	sql, args := sqlize(q)
	return db.Exec(sql, args...)
}

func query(q sq.Sqlizer, db DB) (*sql.Rows, error) {
	sql, args := sqlize(q)
	return db.Query(sql, args...)
}

func queryRow(q sq.Sqlizer, db DB) *sql.Row {
	sql, args := sqlize(q)
	return db.QueryRow(sql, args...)
}
