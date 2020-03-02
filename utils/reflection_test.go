package utils_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox-client/utils"
	"github.com/stretchr/testify/assert"
)

//TODO: write more tests

func TestCompareStruct(t *testing.T) {
	item1 := types.CommonEntity{
		Meta: &types.Metadata{
			ID:           1,
			NetboxEntity: interface{}(10),
		},
	}

	item2 := types.CommonEntity{
		Meta: &types.Metadata{
			ID:           1,
			NetboxEntity: interface{}(10),
		},
	}

	assert.True(t, utils.CompareStruct(item1, item2, nil, nil))
}
