package y

import (
	"reflect"

	sq "github.com/masterminds/squirrel"
)

type wrapper interface {
	field(string) reflect.Value
	addr() reflect.Value
	typo() reflect.Type
}

type value interface {
	wrapper
	put(sq.BaseRunner, schema) (int64, error)
	index(int) value
	size() int
	addTo(*Collection)
}

type ptr struct {
	reflect.Value
}

func (v ptr) field(name string) reflect.Value {
	return v.Elem().FieldByName(name)
}

func (v ptr) addr() reflect.Value {
	return v.Value
}

func (v ptr) typo() reflect.Type {
	return v.Elem().Type()
}

type elem struct {
	reflect.Value
}

func (v elem) field(name string) reflect.Value {
	return v.FieldByName(name)
}

func (v elem) addr() reflect.Value {
	return v.Addr()
}

func (v elem) typo() reflect.Type {
	return v.Type()
}

type singular struct {
	wrapper
}

func (v singular) index(i int) value {
	panic("y/value: singular value doesn't support index method.")
}

func (v singular) size() int {
	panic("y/value: singular value doesn't support size method.")
}

func (v singular) addTo(c *Collection) {
	c.add(v)
}

type plural struct {
	reflect.Value
}

func (v plural) field(name string) reflect.Value {
	panic("y/value: plural value doesn't support field method.")
}

func (v plural) addr() reflect.Value {
	panic("y/value: plural value doesn't support addr method.")
}

func (v plural) index(i int) value {
	return valueOf(v.Index(i))
}

func (v plural) size() int {
	return v.Len()
}

func (v plural) addTo(c *Collection) {
	for i, l := 0, v.size(); i < l; i++ {
		v.index(i).addTo(c)
	}
}

func (v plural) typo() reflect.Type {
	t := v.Type().Elem()
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func valueOf(v reflect.Value) value {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return plural{v}
	case reflect.Ptr:
		return singular{ptr{v}}
	default:
		return singular{elem{v}}
	}
}
