package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO: move to inventory_item.go
func TestIsEqual(t *testing.T) {
	cases := []struct {
		host1   DedicatedServer
		host2   DedicatedServer
		isEqual bool
	}{
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
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
				InventoryItems: []*InventoryItem{
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
				InventoryItems: []*InventoryItem{
					{
						Manufacturer: "unknown",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						Manufacturer: "u. n. owen",
					},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						Model: "unknown",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						Model: "u. n. owen",
					},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						PartNumber: "unknown",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						PartNumber: "u. n. owen",
					},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						AssetTag: "unknown",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						AssetTag: "u. n. owen",
					},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						SerialNumber: "unknown",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						SerialNumber: "u. n. owen",
					},
				},
			},
			isEqual: false,
		},
		{
			host1: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						Manufacturer: "Intel",
					},
					{
						Manufacturer: "AMD",
					},
				},
			},
			host2: DedicatedServer{
				InventoryItems: []*InventoryItem{
					{
						Manufacturer: "AMD",
					},
					{
						Manufacturer: "Intel",
					},
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
