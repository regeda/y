package y

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Versionable mixins a version to a model
type Versionable struct {
	Version sql.NullInt64 `json:"-" y:"_version"`
}

// MakeVersionable inits a new version
func MakeVersionable(n int64) Versionable {
	return Versionable{
		sql.NullInt64{Int64: n, Valid: true},
	}
}

// Qualifier updates a select builder if you need
type Qualifier func(sq.SelectBuilder) sq.SelectBuilder

// ByEq returns the filter by squirrel.Eq
var ByEq = func(eq sq.Eq) Qualifier {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(eq)
	}
}

// ByID returns the filter by ID
var ByID = func(id interface{}) Qualifier {
	return ByEq(sq.Eq{"id": id})
}

// TxPipe run some db statement
type TxPipe func(db sq.BaseRunner, v interface{}) error

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
