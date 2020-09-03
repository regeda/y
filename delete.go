package y

import (
	sq "github.com/Masterminds/squirrel"
)

type truncator struct {
	table string
}

func (t truncator) ToSql() (string, []interface{}, error) {
	return "TRUNCATE " + t.table, []interface{}{}, nil
}

// Truncate erases all data from a table
func Truncate(db sq.BaseRunner, p *Proxy) (err error) {
	_, err = sq.ExecWith(db, truncator{p.schema.table})
	return
}

// DeleteBy removes a proxy by values
func DeleteBy(db sq.BaseRunner, p *Proxy, by Values) (int64, error) {
	result, err := builder{p.schema}.forDelete(sq.Eq(by)).RunWith(db).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
