package y

import sq "gopkg.in/Masterminds/squirrel.v1"

type builder struct {
	schema
}

func (b builder) forFinder() sq.SelectBuilder {
	table := b.table
	cols := b.fseq.alias(table)
	return sq.Select(cols...).From(table)
}

func (b builder) forUpdate(set sq.Eq, where sq.Eq) sq.UpdateBuilder {
	return sq.Update(b.table).SetMap(set).Where(where)
}

func (b builder) forDelete(where sq.Eq) sq.DeleteBuilder {
	return sq.Delete(b.table).Where(where)
}
