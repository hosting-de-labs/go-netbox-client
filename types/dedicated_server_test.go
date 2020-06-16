package types_test

import (
	"strconv"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/stretchr/testify/assert"
)

func mockDedicatedServer() types.DedicatedServer {
	d := types.NewDedicatedServer()
	d.Hostname = "host1"
	d.IsManaged = false
	d.Inventory = []*types.InventoryItem{
		{
			Type:         types.InventoryItemTypeProcessor,
			Manufacturer: "Intel",
			Model:        "Xeon X5660",
		},
	}

	return *d
}

func TestDedicatedServer_Copy(t *testing.T) {
	host1 := mockDedicatedServer()
	host2 := host1.Copy()
	assert.True(t, host1.IsEqual(host2, true))

	host1.Inventory = append(host1.Inventory, &types.InventoryItem{
		Type:         types.InventoryItemTypeMainboard,
		Manufacturer: "Supermicro",
		Model:        "X9SCL-F",
	})
	assert.False(t, host1.IsEqual(host2, true))
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
				},
			},
			isEqual: false,
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
