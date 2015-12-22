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
func Truncate(db DB, p *Proxy) (err error) {
	_, err = exec(truncator{p.schema.table}, db)
	return
}

// DeleteBy removes a proxy by values
func DeleteBy(db DB, p *Proxy, by Values) (int64, error) {
	result, err := exec(builder{p.schema}.forDelete(sq.Eq(by)), db)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
