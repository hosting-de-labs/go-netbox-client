package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockDedicatedServer() types.DedicatedServer {
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
