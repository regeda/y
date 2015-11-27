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

// Modifier changes a value for update statement
type Modifier func(v interface{}) interface{}

// IncrInt returns a modifier for int/int8/int16/int32/int64 increment
func IncrInt(to interface{}) Modifier {
	return func(v interface{}) interface{} {
		if x, ok := v.(int); ok {
			return x + to.(int)
		}
		if x, ok := v.(int8); ok {
			return x + to.(int8)
		}
		if x, ok := v.(int16); ok {
			return x + to.(int16)
		}
		if x, ok := v.(int32); ok {
			return x + to.(int32)
		}
		if x, ok := v.(int64); ok {
			return x + to.(int64)
		}
		panic("y/update: unknown int type for increment.")
	}
}

// IncrFloat returns a modifier for float32/float64 increment
func IncrFloat(to interface{}) Modifier {
	return func(v interface{}) interface{} {
		if x, ok := v.(float32); ok {
			return x + to.(float32)
		}
		if x, ok := v.(float64); ok {
			return x + to.(float64)
		}
		panic("y/update: unknown float type for increment.")
	}
}

func (c *Changer) modify() sq.Eq {
	modified := sq.Eq{}
	for name, val := range c.values {
		f := c.proxy.Field(name)
		if modifier, ok := val.(Modifier); ok {
			val = modifier(f.Interface())
		}
		if f.Interface() != val {
			modified[name] = val
			f.Set(reflect.ValueOf(val))
		}
	}
	return modified
}

// Update saves object changes in db after version validation
func (c *Changer) Update(db DB) (err error) {
	pk := c.proxy.primary()
	// load origin
	err = c.proxy.loadBy(db, pk)
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
	// add version to search condition
	pk[_version] = oldv
	// save
	result, err := exec(
		builder{c.proxy.schema}.forUpdate(clauses, sq.Eq(pk)), db)
	if err == nil {
		count, _ := result.RowsAffected()
		if count != 1 {
			err = ErrNoAffectedRows
		}
	}
	return
}

func makeChanger(p *Proxy, v Values) *Changer {
	return &Changer{p, v}
}

// Update saves object changes
func Update(db DB, p *Proxy, v Values) error {
	return makeChanger(p, v).Update(db)
}
