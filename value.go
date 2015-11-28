package y

import "reflect"

type value interface {
	put(DB, *schema) (int64, error)
	ptr() reflect.Value
	field(string) reflect.Value
	index(int) value
	size() int
	deploy(*Collection)
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

func (v singular) deploy(c *Collection) {
	c.add(v.Value)
}

func (v singular) ptr() reflect.Value {
	return v.Addr()
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

func (v plural) deploy(c *Collection) {
	for i, l := 0, v.size(); i < l; i++ {
		v.index(i).deploy(c)
	}
}

func (v plural) ptr() reflect.Value {
	return v.Addr()
}
