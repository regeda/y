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
	items []reflect.Value
	idx   map[string]*index
	schema
}

func (c *Collection) lookidx(name string) *index {
	idx, ok := c.idx[name]
	if !ok {
		log.Panicf(
			"y/collection: The index \"%s\" not found in collection \"%s\".",
			name, c.table)
	}
	return idx
}

func (c *Collection) add(v value) {
	cell := len(c.items)
	c.items = append(c.items, v.addr())

	for name := range c.schema.xinfo.idx {
		key := c.schema.fval(v, name).Int()
		c.lookidx(name).add(key, cell)
	}
}

func (c *Collection) cells(cells []int) []reflect.Value {
	items := make([]reflect.Value, len(cells))
	for i, cell := range cells {
		items[i] = c.items[cell]
	}
	return items
}

// Get returns an item by primary key
func (c *Collection) Get(pk ...interface{}) interface{} {
	fields := c.schema.xinfo.pk
	flen := len(fields)
	if flen == 0 {
		log.Panicf(
			"y/colleciton: no primary key in the schema definition \"%s\".", c.table)
	}
	if flen != len(pk) {
		log.Panicln("y/collection: missing primary key parameters.")
	}
	idx := c.lookidx(fields[0])
CellLoop:
	for _, cell := range idx.cells[pk[0].(int64)] {
		item := valueOf(c.items[cell])
		// matching by a composite primary key
		for i := 1; i < flen; i++ {
			if c.schema.fval(item, fields[i]).Interface() != pk[i] {
				continue CellLoop
			}
		}
		return item.addr().Interface()
	}
	return nil
}

// Empty returns false if no items exist
func (c *Collection) Empty() bool {
	return c.Size() == 0
}

// First returns the first item
func (c *Collection) First() interface{} {
	return c.items[0].Interface()
}

// Size returns count of items
func (c *Collection) Size() int {
	return len(c.items)
}

// List returns all items
func (c *Collection) List() interface{} {
	size := c.Size()
	items := reflect.MakeSlice(c.schema.sliceOf(), size, size)
	for i, item := range c.items {
		items.Index(i).Set(item)
	}
	return items.Interface()
}

// Join links related collection
func (c *Collection) Join(j *Collection) {
	fk := j.schema.fk(c.schema)

	cidx := c.lookidx(fk.target)
	jidx := j.lookidx(fk.from)

	name := j.schema.t.Name()

	for jkey, jcells := range jidx.cells {
		if ccells, ok := cidx.cells[jkey]; ok {
			for _, ccell := range ccells {
				citem := c.items[ccell].Elem()
				// one-to-many
				target := citem.FieldByName(name + "Array")
				if target.CanSet() {
					items := j.cells(jcells)
					target.Set(reflect.Append(target, items...))
					continue
				}
				// one-to-one
				target = citem.FieldByName(name)
				if target.CanSet() && len(jcells) == 1 {
					target.Set(j.items[jcells[0]])
				}
			}
		}
	}
}

func makeCollection(p *Proxy) *Collection {
	// create the index map
	idx := make(map[string]*index)
	for name := range p.schema.xinfo.idx {
		idx[name] = makeIndex()
	}
	return &Collection{
		idx:    idx,
		schema: p.schema,
	}
}
