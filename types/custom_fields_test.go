package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hosting-de-labs/go-netbox-client/types"
)

func createNetboxCustomFields() interface{} {
	fields := make(map[string]interface{}, 2)

	extfield := make(map[string]interface{}, 2)
	extfield["Value"] = 1
	extfield["Label"] = "Value 1"

	fields["ext_field"] = extfield
	fields["string_field"] = "string"

	var out interface{}
	out = fields

	return out
}

func TestCustomFields_Load(t *testing.T) {
	cf := createNetboxCustomFields()
	c := types.CustomFields{}
	assert.Nil(t, c.Load(cf))
}

func TestCustomFields_Val(t *testing.T) {
	cf := createNetboxCustomFields()
	c := types.CustomFields{}
	assert.Nil(t, c.Load(cf))

	assert.NotNil(t, c.Val("string_field"))
	assert.Equal(t, "string", *c.Val("string_field"))
	assert.Equal(t, "Value 1", *c.Val("ext_field"))
}

func TestCustomFields_ValMap(t *testing.T) {
	cf := createNetboxCustomFields()
	c := types.CustomFields{}
	assert.Nil(t, c.Load(cf))

	m := c.ValMap()

	assert.Equal(t, "string", m["string_field"].(string))
	assert.Equal(t, 1, m["ext_field"].(int))
}
