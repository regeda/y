package y

import (
	"log"
	"reflect"
)

type index struct {
	cells map[int64][]int
	keys  []int64
}

func (idx *index) add(key int64, cell int) {
	cells, ok := idx.cells[key]
	if !ok {
		idx.keys = append(idx.keys, key)
	}
	idx.cells[key] = append(cells, cell)
}

func makeIndex() *index {
	return &index{
		cells: make(map[int64][]int),
	}
}

// Collection contains items and indexes
type Collection struct {
	items  []interface{}
	idx    map[string]*index
	schema *schema
}

func (c *Collection) getIdx(name string) *index {
	if idx, ok := c.idx[name]; ok {
		return idx
	}
	c.idx[name] = makeIndex()
	return c.idx[name]
}

func (c *Collection) add(v reflect.Value) {
	cell := len(c.items)
	c.items = append(c.items, v.Addr().Interface())

	for name := range c.schema.xinfo.idx {
		key := v.Field(c.schema.fields[name].i).Int()
		c.getIdx(name).add(key, cell)
	}
}

// Empty returns false if no items exist
func (c *Collection) Empty() bool {
	return len(c.items) == 0
}

// First returns the first item
func (c *Collection) First() interface{} {
	return c.items[0]
}

// List returns all items
func (c *Collection) List() []interface{} {
	return c.items
}

// Join links related collection
func (c *Collection) Join(j *Collection) {
	fk := j.schema.fk(c.schema)

	cidx, ok := c.idx[fk.target]
	if !ok {
		log.Panicf(
			"y/collection: The index \"%s\" not found in collection \"%s\".",
			fk.target, c.schema.table)
	}
	jidx, ok := j.idx[fk.from]
	if !ok {
		log.Panicf(
			"y/collection: The index \"%s\" not found in collection \"%s\".",
			fk.from, j.schema.table)
	}

	name := j.schema.t.Name()

	for jkey, jcells := range jidx.cells {
		if ccells, ok := cidx.cells[jkey]; ok {
			for _, ccell := range ccells {
				citem := reflect.ValueOf(c.items[ccell]).Elem()
				// one-to-many
				target := citem.FieldByName(name + "s")
				if target.CanSet() {
					reflected := make([]reflect.Value, len(jcells))
					for i, jcell := range jcells {
						reflected[i] = reflect.ValueOf(j.items[jcell])
					}
					target.Set(reflect.Append(target, reflected...))
					continue
				}
				// one-to-one
				target = citem.FieldByName(name)
				if target.CanSet() && len(jcells) == 1 {
					target.Set(reflect.ValueOf(j.items[jcells[0]]))
				}
			}
		}
	}
}

func makeCollection(p *Proxy) *Collection {
	return &Collection{
		idx:    make(map[string]*index),
		schema: p.schema,
	}
}
