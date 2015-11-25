package y

import (
	"reflect"

	sq "github.com/lann/squirrel"
)

// Proxy contains a schema of a type
type Proxy struct {
	v      value
	schema *schema
}

// Put creates a new object
func (p *Proxy) Put(db DB) (int64, error) {
	return Put(db, p)
}

// Fetch returns a collection of objects
func (p *Proxy) Fetch(db DB) (*Collection, error) {
	return p.Find(NoopQualifier).Fetch(db)
}

// Find returns Finder
func (p *Proxy) Find(q Qualifier) *Finder {
	return makeFinder(p, q)
}

// Collection creates an empty collection of the object
func (p *Proxy) Collection() *Collection {
	return makeCollection(p)
}

// Join adds related collection to self
func (p *Proxy) Join(db DB, in *Collection) (*Collection, error) {
	fk := p.schema.fk(in.schema)

	idx := in.getIdx(fk.target)

	c, err := p.Find(ByEq(sq.Eq{fk.from: idx.keys})).Fetch(db)
	if err != nil {
		return nil, err
	}
	if !c.Empty() {
		in.Join(c)
	}
	return c, nil
}

// Load fetches an object from db by primary key
func (p *Proxy) Load(db DB) error {
	return p.loadBy(db, p.primary())
}

// Version returns a revision of object
func (p *Proxy) Version() int64 {
	return p.Field(_version).Int()
}

// Field returns reflected field by name
func (p *Proxy) Field(name string) reflect.Value {
	return p.schema.field(p.v, name)
}

// Map returns a simple map of struct values
func (p *Proxy) Map() Values {
	return p.schema.mapped(p.v)
}

// Update saves changes of the object
func (p *Proxy) Update(db DB, v Values) error {
	return Update(db, p, v)
}

// Truncate erases all data
func (p *Proxy) Truncate(db DB) error {
	return Truncate(db, p)
}

func (p *Proxy) primary() Values {
	return p.schema.pk(p.v)
}

func (p *Proxy) loadBy(db DB, eq Values) error {
	return p.Find(ByEq(sq.Eq(eq))).Load(db)
}

func proxyOf(v reflect.Value) *Proxy {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return &Proxy{plural{v}, loadSchema(v.Type().Elem())}
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	default:
		return &Proxy{singular{v}, loadSchema(v.Type())}
	}
}
