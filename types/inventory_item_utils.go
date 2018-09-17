package types

import (
	"fmt"
	"sort"
)

//InventoryItemsSort sorts a slice of inventory items by manufacturer, model, asset tag
func InventoryItemsSort(items []*InventoryItem) {
	sort.Slice(items, func(i, j int) bool {
		string1 := items[i].Manufacturer + items[i].Model + items[i].AssetTag
		string2 := items[j].Manufacturer + items[j].Model + items[j].AssetTag

		return string1 < string2
	})
}

//InventoryItemTypeParse returns an InventoryItemType based on the short string
func InventoryItemTypeParse(s string) (out InventoryItemType, err error) {
	switch s {
	case "Processor":
		fallthrough
	case "CPU":
		return InventoryItemTypeProcessor, nil

	case "Harddisk Drive":
		fallthrough
	case "HDD":
		return InventoryItemTypeHarddrive, nil

	case "Solid State Drive":
		fallthrough
	case "SSD":
		return InventoryItemTypeSolidStateDrive, nil

	case "Drive":
		fallthrough
	case "DRIVE":
		return InventoryItemTypeGenericDrive, nil

	case "RAID Controller":
		fallthrough
	case "RAID":
		return InventoryItemTypeRAIDController, nil

	case "Mainboard":
		fallthrough
	case "MB":
		return InventoryItemTypeMainboard, nil

	case "Memory Module":
		fallthrough
	case "RAM":
		return InventoryItemTypeMemoryModule, nil

	case "Baseband Management Controller":
		fallthrough
	case "BMC":
		return InventoryItemTypeBasebandManagementController, nil

	case "Powersupply":
		fallthrough
	case "PSU":
		return InventoryItemTypePowersupply, nil

	case "Other":
		fallthrough
	case "OTHER":
		return InventoryItemTypeOther, nil

	default:
		return InventoryItemTypeOther, fmt.Errorf("cannot parse %s to a matching InventoryItemType", s)
	}
}
