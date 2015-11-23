package y

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaParseAuto(t *testing.T) {
	type something struct {
		X int64
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Len(p.schema.fields, 1)
	assert.NotNil(p.schema.fields["x"])
}

func TestSchemaParseAutoWithOpts(t *testing.T) {
	type something struct {
		X int64 `db:",pk"`
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Len(p.schema.fields, 1)
	assert.NotNil(p.schema.fields["x"])
	assert.Equal([]string{"x"}, p.schema.xinfo.pk)
}

func TestSchemaParseEmpty(t *testing.T) {
	type something struct {
		X int64 `db:"-"`
	}
	p := New(something{})
	assert.Empty(t, p.schema.fields)
}

func TestCompositeSchema(t *testing.T) {
	type foo struct {
		x int64 `db:"x"`
		y int64 `db:"y"`
	}
	type bar struct {
		foo
		z int64 `db:"z"`
	}
	p := New(bar{})
	assert := assert.New(t)
	assert.Len(p.schema.fields, 3)
	assert.NotNil(p.schema.fields["x"])
	assert.NotNil(p.schema.fields["y"])
	assert.NotNil(p.schema.fields["z"])
}

func TestSchemaParseArray(t *testing.T) {
	type something struct {
		x int64
	}
	p := New([]something{})
	assert := assert.New(t)
	assert.Len(p.schema.fields, 1)
	assert.NotNil(p.schema.fields["x"])
}

func TestSchemaParseSinglePK(t *testing.T) {
	type something struct {
		X int64 `db:"x,pk"`
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Equal([]string{"x"}, p.schema.xinfo.pk)
	assert.Equal(1, p.schema.xinfo.idx["x"])
}

func TestSchemaParseCompositePK(t *testing.T) {
	type something struct {
		X int64 `db:"x,pk"`
		Y int64 `db:"y,pk"`
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Equal([]string{"x", "y"}, p.schema.xinfo.pk)
	assert.Equal(1, p.schema.xinfo.idx["x"])
	assert.Equal(0, p.schema.xinfo.idx["y"])
}

func TestSchemaParseAutoincr(t *testing.T) {
	type something struct {
		X int64 `db:"x,autoincr"`
	}
	p := New(something{})
	assert.True(t, p.schema.fields["x"].opts.autoincr)
}

func TestSchemaParseImplicitFK(t *testing.T) {
	type something struct {
		XID int64 `db:"x_id,fk"`
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Equal("x_id", p.schema.fks["x"].from)
	assert.Equal("id", p.schema.fks["x"].target)
}

func TestSchemaParseExplicitFK(t *testing.T) {
	type something struct {
		XID int64 `db:"xid,fk:xoo.id"`
	}
	p := New(something{})
	assert := assert.New(t)
	assert.Equal("xid", p.schema.fks["xoo"].from)
	assert.Equal("id", p.schema.fks["xoo"].target)
}
