package y

import sq "github.com/lann/squirrel"

// Finder loads a collection from a database
type Finder struct {
	proxy   *Proxy
	builder sq.SelectBuilder
}

// Qualify updates select builder
func (f *Finder) Qualify(q Qualifier) *Finder {
	f.builder = q(f.builder)
	return f
}

// Load fetches an object from db and loads in self proxy
func (f *Finder) Load(db DB) error {
	row := queryRow(f.builder, db)
	ptrs := f.proxy.schema.ptrs()
	f.proxy.schema.set(ptrs, f.proxy.v)
	return row.Scan(ptrs...)
}

// Fetch make a query and creates a collection
func (f *Finder) Fetch(db DB) (*Collection, error) {
	rows, err := query(f.builder, db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	c := f.proxy.blankCollection()
	ptrs := f.proxy.schema.ptrs()
	for rows.Next() {
		v := f.proxy.schema.create()
		f.proxy.schema.set(ptrs, v)
		rows.Scan(ptrs...)
		v.addTo(c)
	}
	return c, nil
}

func makeFinder(p *Proxy, b sq.SelectBuilder) *Finder {
	return &Finder{p, b}
}

// Fetch loads a collection of objects
func Fetch(db DB, v interface{}) (*Collection, error) {
	return New(v).Fetch(db)
}
