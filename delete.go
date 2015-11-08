package y

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
