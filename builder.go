package y

import (
	"fmt"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

type provider interface {
	grabInsertID(schema, value, sq.InsertBuilder) (int64, error)
	placeholder() sq.PlaceholderFormat
}

type statementBuilder struct {
	provider
	sq.StatementBuilderType
}

var (
	// MySQL describes behavior for LastInsertId and question placeholder
	MySQL = mysqlProvider{}
	// Postgres describes behavior for SERIAL data type and dollar placeholder
	Postgres = postgresProvider{}

	stmtBuilder statementBuilder
)

type mysqlProvider struct{}

func (m mysqlProvider) grabInsertID(s schema, v value, b sq.InsertBuilder) (int64, error) {
	result, err := b.Exec()
	if err != nil {
		return 0, err
	}
	for _, pk := range s.xinfo.pk {
		f := s.fields[pk]
		if f.autoincr {
			id, _ := result.LastInsertId()
			v.field(f.Name).SetInt(id)
			break
		}
	}
	return result.RowsAffected()
}

func (m mysqlProvider) placeholder() sq.PlaceholderFormat {
	return sq.Question
}

type postgresProvider struct{}

func (p postgresProvider) grabInsertID(s schema, v value, b sq.InsertBuilder) (int64, error) {
	for _, pk := range s.xinfo.pk {
		f := s.fields[pk]
		if f.autoincr {
			var id int64
			row := b.Suffix(fmt.Sprintf("RETURNING \"%s\"", pk)).QueryRow()
			if err := row.Scan(&id); err != nil {
				return 0, err
			}
			v.field(f.Name).SetInt(id)
			return 1, nil
		}
	}

	result, err := b.Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (p postgresProvider) placeholder() sq.PlaceholderFormat {
	return sq.Dollar
}

// SetBuilderProvider defines statement builder type
func SetBuilderProvider(rdb provider) {
	stmtBuilder = statementBuilder{
		rdb,
		sq.StatementBuilder.PlaceholderFormat(rdb.placeholder()),
	}
}

type builder struct {
	schema
}

func (b builder) grabInsertID(v value, stmt sq.InsertBuilder) (int64, error) {
	return stmtBuilder.grabInsertID(b.schema, v, stmt)
}

func (b builder) forInsert() sq.InsertBuilder {
	return stmtBuilder.Insert(b.table)
}

func (b builder) forFinder() sq.SelectBuilder {
	table := b.table
	cols := b.fseq.alias(table)
	return stmtBuilder.Select(cols...).From(table)
}

func (b builder) forUpdate(set sq.Eq, where sq.Eq) sq.UpdateBuilder {
	return stmtBuilder.Update(b.table).SetMap(set).Where(where)
}

func (b builder) forDelete(where sq.Eq) sq.DeleteBuilder {
	return stmtBuilder.Delete(b.table).Where(where)
}

func init() {
	SetBuilderProvider(MySQL)
}
