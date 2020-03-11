package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockInventoryItem() types.InventoryItem {
	return types.InventoryItem{
		Type:         types.InventoryItemTypeProcessor,
		Manufacturer: "Intel",
		Model:        "Xeon X5670",
		AssetTag:     "Asset Tag",
		PartNumber:   "Part Number",
		SerialNumber: "Serial Number",
	}
}
