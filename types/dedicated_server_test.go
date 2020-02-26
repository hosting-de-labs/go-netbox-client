package types_test

import (
	"strconv"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"
)

func TestDedicatedServer_Copy(t *testing.T) {
	host1 := client_types.MockDedicatedServer()
	host2 := host1.Copy()
	assert.Equal(t, host1, host2)

	host1.Inventory = append(host1.Inventory, &types.InventoryItem{
		Type:         types.InventoryItemTypeMainboard,
		Manufacturer: "Supermicro",
		Model:        "X9SCL-F",
	})
	assert.NotEqual(t, host1, host2)
}

//TODO: move to inventory_item.go
func TestDedicatedServer_IsEqual(t *testing.T) {
	cases := []struct {
		host1   types.DedicatedServer
		host2   types.DedicatedServer
		isEqual bool
	}{
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{
						Manufacturer: "unknown",
						Model:        "unknown",
						AssetTag:     "asset tag",
						PartNumber:   "part number",
						SerialNumber: "serial number",
					},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{
						Manufacturer: "unknown",
						Model:        "unknown",
						AssetTag:     "asset tag",
						PartNumber:   "part number",
						SerialNumber: "serial number",
					},
				},
			},
			isEqual: true,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Manufacturer: "unknown"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Manufacturer: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Model: "unknown"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Model: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{PartNumber: "unknown"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{PartNumber: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{AssetTag: "unknown"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{AssetTag: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{SerialNumber: "unknown"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{SerialNumber: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Manufacturer: "Intel"},
					{Manufacturer: "AMD"},
				},
			},
			host2: types.DedicatedServer{
				Inventory: []*types.InventoryItem{
					{Manufacturer: "AMD"},
					{Manufacturer: "Intel"},
				},
			},
			isEqual: true,
		},
		//TODO: Details
	}

	for key, testcase := range cases {
		if testcase.isEqual {
			assert.True(t, testcase.host1.IsEqual(testcase.host2, true), "Case ID: "+strconv.Itoa(key))
		} else {
			assert.False(t, testcase.host1.IsEqual(testcase.host2, true), "Case ID: "+strconv.Itoa(key))
		}
	}
}
