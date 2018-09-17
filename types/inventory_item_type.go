package types

type InventoryItemType int

const (
	_ InventoryItemType = iota
	InventoryItemTypeProcessor
	InventoryItemTypeHarddrive
	InventoryItemTypeSolidStateDrive
	InventoryItemTypeGenericDrive
	InventoryItemTypeRAIDController
	InventoryItemTypeMainboard
	InventoryItemTypeMemoryModule
	InventoryItemTypeBasebandManagementController
	InventoryItemTypePowersupply
	InventoryItemTypeOther
)

//String implements the Stringer interface
func (i InventoryItemType) String() string {
	switch i {
	case InventoryItemTypeProcessor:
		return "Processor"
	case InventoryItemTypeHarddrive:
		return "Harddisk Drive"
	case InventoryItemTypeSolidStateDrive:
		return "Solid State Drive"
	case InventoryItemTypeGenericDrive:
		return "Drive"
	case InventoryItemTypeRAIDController:
		return "RAID Controller"
	case InventoryItemTypeMainboard:
		return "Mainboard"
	case InventoryItemTypeMemoryModule:
		return "Memory Module"
	case InventoryItemTypeBasebandManagementController:
		return "Baseband Management Controller"
	case InventoryItemTypePowersupply:
		return "Powersupply"
	default:
		return "Other"
	}
}

//MarshalJSON implements the Marshaller interface
func (i InventoryItemType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + i.String() + "\""), nil
}

//Short returns a short identifier for a given InventoryItemType
func (i InventoryItemType) Short() string {
	switch i {
	case InventoryItemTypeProcessor:
		return "CPU"
	case InventoryItemTypeHarddrive:
		return "HDD"
	case InventoryItemTypeSolidStateDrive:
		return "SSD"
	case InventoryItemTypeGenericDrive:
		return "DRIVE"
	case InventoryItemTypeRAIDController:
		return "RAID"
	case InventoryItemTypeMainboard:
		return "MB"
	case InventoryItemTypeMemoryModule:
		return "RAM"
	case InventoryItemTypeBasebandManagementController:
		return "BMC"
	case InventoryItemTypePowersupply:
		return "PSU"
	default:
		return "OTHER"
	}
}
