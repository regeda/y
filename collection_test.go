package y

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCollection(t *testing.T) {
	type something struct {
	}
	p := New(something{})
	c := p.Collection()
	assert.True(t, c.Empty())
}

func TestAddOneItemToCollection(t *testing.T) {
	type something struct {
		ID int64 `db:"id"`
	}
	p := New(something{})
	c := p.Collection()

	v1 := p.schema.create().Elem()
	ptr1 := v1.Addr().Interface().(*something)
	ptr1.ID = 1
	c.add(v1)

	assert := assert.New(t)
	assert.False(c.Empty())
	assert.Equal(ptr1, c.First())
}

func TestAddTwoItemsToCollection(t *testing.T) {
	type something struct {
		ID int64 `db:"id"`
	}
	p := New(something{})
	c := p.Collection()

	ptrs := []interface{}{}
	for _, id := range []int64{1, 2} {
		v := p.schema.create().Elem()
		c.add(v)
		ptr := v.Addr().Interface().(*something)
		ptrs = append(ptrs, ptr)
		ptr.ID = id
	}

	assert := assert.New(t)
	assert.Equal(2, c.Size())
	assert.Equal(ptrs, c.List())
}
