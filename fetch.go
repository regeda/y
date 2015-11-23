package y

import sq "github.com/lann/squirrel"

// Finder loads a collection from a database
type Finder struct {
	qualifier Qualifier
	proxy     *Proxy
}

func (f *Finder) prepare() sq.Sqlizer {
	table := f.proxy.schema.table
	cols := f.proxy.schema.fseq.alias(table)
	return f.qualifier(sq.Select(cols...).From(table))
}

// Load fetches an object from db and loads in self proxy
func (f *Finder) Load(db DB) error {
	row := queryRow(f.prepare(), db)
	ptrs := f.proxy.schema.ptrs()
	f.proxy.schema.set(ptrs, f.proxy.v)
	return row.Scan(ptrs...)
}

// Fetch make a query and creates a collection
func (f *Finder) Fetch(db DB) (*Collection, error) {
	rows, err := query(f.prepare(), db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	c := f.proxy.Collection()
	ptrs := f.proxy.schema.ptrs()
	for rows.Next() {
		v := f.proxy.schema.create().Elem()
		f.proxy.schema.set(ptrs, singular{v})
		rows.Scan(ptrs...)
		c.add(v)
	}
	return c, nil
}

func makeFinder(p *Proxy, q Qualifier) *Finder {
	return &Finder{q, p}
}

// Fetch loads a collection of objects
func Fetch(db DB, v interface{}) (*Collection, error) {
	return New(v).Fetch(db)
}
