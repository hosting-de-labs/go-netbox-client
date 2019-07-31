package utils_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox-client/utils"
	"github.com/stretchr/testify/assert"
)

func TestCompareStruct(t *testing.T) {
	item1 := types.CommonEntity{OriginalEntity: interface{}(10)}
	item2 := types.CommonEntity{OriginalEntity: interface{}(10)}

	assert.True(t, utils.CompareStruct(item1, item2, nil, nil))
}
