package y

import (
	"database/sql"
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

// Proxy contains a schema of a type
type Proxy struct {
	v value
	schema
}

// Put creates a new object
func (p *Proxy) Put(db sq.BaseRunner) (int64, error) {
	return Put(db, p)
}

// Query returns Finder with a custom query
func (p *Proxy) Query(b sq.SelectBuilder) *Finder {
	return makeFinder(p, b)
}

// Fetch returns a collection of objects
func (p *Proxy) Fetch(db sq.BaseRunner) (*Collection, error) {
	return p.Find().Fetch(db)
}

// FindBy returns Finder with qualified query
func (p *Proxy) FindBy(q Qualifier) *Finder {
	return p.Find().Qualify(q)
}

// Find returns Finder with qualified query
func (p *Proxy) Find() *Finder {
	return p.Query(builder{p.schema}.forFinder())
}

// Collection creates a collection of proxy values
func (p *Proxy) Collection() *Collection {
	c := p.blankCollection()
	p.v.addTo(c)
	return c
}

// Join adds related collection to self
func (p *Proxy) Join(db sq.BaseRunner, in *Collection) (*Collection, error) {
	fk := p.schema.fk(in.schema)

	idx := in.lookidx(fk.target)

	c, err := p.findByEq(sq.Eq{fk.from: idx.keys}).Fetch(db)
	if err != nil {
		return nil, err
	}
	if !c.Empty() {
		in.Join(c)
	}
	return c, nil
}

// Load fetches an object from db by primary key
func (p *Proxy) Load(db sq.BaseRunner) error {
	return p.loadBy(db, p.primary())
}

// MustLoad fetches an object and panic if an error occurred
func (p *Proxy) MustLoad(db sq.BaseRunner) *Proxy {
	err := p.Load(db)
	if err != nil {
		panic(err.Error())
	}
	return p
}

// Field returns reflected field by name
func (p *Proxy) Field(name string) reflect.Value {
	return p.schema.fval(p.v, name)
}

// Map returns a simple map of struct values
func (p *Proxy) Map() Values {
	return p.schema.mapped(p.v)
}

// Update saves changes of the object
func (p *Proxy) Update(db sq.BaseRunner, v Values) error {
	return Update(db, p, v)
}

// Truncate erases all data
func (p *Proxy) Truncate(db sq.BaseRunner) error {
	return Truncate(db, p)
}

// Delete removes a proxy by primary
func (p *Proxy) Delete(db sq.BaseRunner) (int64, error) {
	return p.DeleteBy(db, p.primary())
}

// DeleteBy removes a proxy by values
func (p *Proxy) DeleteBy(db sq.BaseRunner, by Values) (int64, error) {
	return DeleteBy(db, p, by)
}

func (p *Proxy) blankCollection() *Collection {
	return makeCollection(p)
}

func (p *Proxy) findByEq(eq sq.Eq) *Finder {
	return p.FindBy(ByEq(eq))
}

func (p *Proxy) primary() Values {
	return p.schema.pk(p.v)
}

func (p *Proxy) version() *sql.NullInt64 {
	return p.Field(_version).Addr().Interface().(*sql.NullInt64)
}

func (p *Proxy) loadBy(db sq.BaseRunner, eq Values) error {
	return p.findByEq(sq.Eq(eq)).Load(db)
}

func proxyOf(v reflect.Value) *Proxy {
	val := valueOf(v)
	return &Proxy{val, schemaOf(val.typo())}
}
