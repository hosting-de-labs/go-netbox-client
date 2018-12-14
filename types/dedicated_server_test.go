package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockDedicatedServer() DedicatedServer {
	return DedicatedServer{
		Host: Host{
			ID:        10,
			Hostname:  "host1",
			IsManaged: false,
		},
		Inventory: []*InventoryItem{
			{
				Type:         InventoryItemTypeProcessor,
				Manufacturer: "Intel",
				Model:        "Xeon X5660",
			},
		},
	}
}

func TestDedicatedServer_Copy(t *testing.T) {
	host1 := mockDedicatedServer()
	host2 := host1.Copy()
	assert.Equal(t, host1, host2)

	host2.Inventory = append(host2.Inventory, &InventoryItem{
		Type:         InventoryItemTypeMainboard,
		Manufacturer: "Supermicro",
		Model:        "X9SCL-F",
	})
	assert.NotEqual(t, host1, host2)
}

//TODO: move to inventory_item.go
func TestDedicatedServer_IsEqual(t *testing.T) {
	cases := []struct {
		host1   DedicatedServer
		host2   DedicatedServer
		isEqual bool
	}{
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{
						Manufacturer: "unknown",
						Model:        "unknown",
						AssetTag:     "asset tag",
						PartNumber:   "part number",
						SerialNumber: "serial number",
					},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
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
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{Manufacturer: "unknown"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
					{Manufacturer: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{Model: "unknown"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
					{Model: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{PartNumber: "unknown"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
					{PartNumber: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{AssetTag: "unknown"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
					{AssetTag: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{SerialNumber: "unknown"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
					{SerialNumber: "u. n. owen"},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				Inventory: []*InventoryItem{
					{Manufacturer: "Intel"},
					{Manufacturer: "AMD"},
				},
			},
			host2: DedicatedServer{
				Inventory: []*InventoryItem{
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
