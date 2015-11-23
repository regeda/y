package y

import "reflect"

type value interface {
	put(DB, *schema) (int64, error)
	field(string) reflect.Value
	index(int) value
	size() int
}

type singular struct {
	reflect.Value
}

func (v singular) field(name string) reflect.Value {
	return v.FieldByName(name)
}

func (v singular) index(i int) value {
	panic("y/value: singular value doesn't support index method.")
}

func (v singular) size() int {
	panic("y/value: singular value doesn't support size method.")
}

type plural struct {
	reflect.Value
}

func (v plural) field(name string) reflect.Value {
	panic("y/value: plural value doesn't support field method.")
}

func (v plural) index(i int) value {
	return singular{v.Index(i)}
}

func (v plural) size() int {
	return v.Len()
}
