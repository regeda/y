package y

import (
	"reflect"

	sq "github.com/lann/squirrel"
)

// Changer updates object values
type Changer struct {
	proxy  *Proxy
	values Values
}

func (c *Changer) modify() sq.Eq {
	modified := sq.Eq{}
	for name, val := range c.values {
		f := c.proxy.Field(name)
		if f.Interface() != val {
			modified[name] = val
			f.Set(reflect.ValueOf(val))
		}
	}
	return modified
}

func (c *Changer) primary() sq.Eq {
	pks := sq.Eq{}
	for _, pk := range c.proxy.schema.xinfo.pk {
		pks[pk] = c.proxy.Field(pk).Interface()
	}
	return pks
}

func (c *Changer) prepare(clauses sq.Eq, where sq.Eq) sq.UpdateBuilder {
	return sq.Update(c.proxy.schema.table).SetMap(clauses).Where(where)
}

// Update saves object changes in db after version validation
func (c *Changer) Update(db DB) (updated bool, err error) {
	// load origin
	err = c.proxy.Load(db)
	if err != nil {
		return
	}
	oldv := c.proxy.Version()
	newv := oldv + 1
	// find changes
	clauses := c.modify()
	if len(clauses) == 0 {
		return
	}
	// set new version
	clauses[_version] = newv
	c.proxy.Field(_version).SetInt(newv)
	// set where condition
	where := c.primary()
	where[_version] = oldv
	// save
	result, err := exec(c.prepare(clauses, where), db)
	if err != nil {
		return
	}
	count, err := result.RowsAffected()
	updated = count == 1
	return
}

func makeChanger(p *Proxy, v Values) *Changer {
	return &Changer{p, v}
}

// Update saves object changes
func Update(db DB, p *Proxy, v Values) (bool, error) {
	return makeChanger(p, v).Update(db)
}
