package y

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnderscore(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("i_love_go_and_json_so_much", underscore("ILoveGoAndJSONSoMuch"))
	assert.Equal("camel_case", underscore("CamelCase"))
	assert.Equal("camel", underscore("Camel"))
	assert.Equal("camel", underscore("CAMEL"))
	assert.Equal("big_case", underscore("BIGCase"))
	assert.Equal("small_case", underscore("SmallCASE"))
}
