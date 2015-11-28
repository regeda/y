package y

import sq "github.com/lann/squirrel"

// Put inserts new objects
func Put(db DB, p *Proxy) (int64, error) {
	return p.v.put(db, p.schema)
}

func (v plural) put(db DB, s *schema) (int64, error) {
	l := v.size()
	if l == 0 {
		return 0, nil
	}
	q := sq.Insert(s.table)
	for i := 0; i < l; i++ {
		ptrs := s.ptrs()
		s.set(ptrs, v.index(i))
		q = q.Values(ptrs...)
	}
	result, err := exec(q, db)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (v singular) put(db DB, s *schema) (int64, error) {
	set := sq.Eq{}
	for name, f := range s.fields {
		if !f.autoincr {
			set[name] = v.field(f.Name).Interface()
		}
	}
	result, err := exec(sq.Insert(s.table).SetMap(set), db)
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
