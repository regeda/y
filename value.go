package y

import "reflect"

type value interface {
	put(DB, *schema) (int64, error)
	field(string) reflect.Value
	addr() reflect.Value
	index(int) value
	size() int
	addTo(*Collection)
}

type singular struct {
	reflect.Value
}

func (v singular) field(name string) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem().FieldByName(name)
	}
	return v.FieldByName(name)
}

func (v singular) addr() reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Value
	}
	return v.Addr()
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
	return singular{v.Index(i)}
}

func (v plural) size() int {
	return v.Len()
}

func (v plural) addTo(c *Collection) {
	for i, l := 0, v.size(); i < l; i++ {
		v.index(i).addTo(c)
	}
}
