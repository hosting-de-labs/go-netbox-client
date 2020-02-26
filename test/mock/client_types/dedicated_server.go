package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockDedicatedServer() types.DedicatedServer {
	return types.DedicatedServer{
		Host: types.Host{
			Hostname:  "host1",
			IsManaged: false,
		},
		Inventory: []*types.InventoryItem{
			{
				Type:         types.InventoryItemTypeProcessor,
				Manufacturer: "Intel",
				Model:        "Xeon X5660",
			},
		},
	}
}
