package y

import sq "github.com/lann/squirrel"

// Put inserts a new object
func Put(db DB, p *Proxy) error {
	set := sq.Eq{}
	for name, f := range p.schema.fields {
		if !f.opts.autoincr {
			set[name] = p.v.FieldByName(f.Name).Interface()
		}
	}
	result, err := exec(sq.Insert(p.schema.table).SetMap(set), db)
	if err != nil {
		return err
	}
	for _, pk := range p.schema.xinfo.pk {
		f := p.schema.fields[pk]
		if f.opts.autoincr {
			id, _ := result.LastInsertId()
			p.v.FieldByName(f.Name).SetInt(id)
			break
		}
	}
	return nil
}
