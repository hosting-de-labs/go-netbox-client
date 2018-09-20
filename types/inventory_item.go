package types

import (
	"encoding/json"
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

//InventoryItem represents a device / module that is attached to / built into a dedicated server
type InventoryItem struct {
	CommonEntity

	Type         InventoryItemType
	Manufacturer string
	Model        string
	PartNumber   string
	SerialNumber string
	AssetTag     string
	Details      map[string]string
}

func (i *InventoryItem) AddDetail(key string, val string) {
	if i.Details == nil {
		i.Details = make(map[string]string)
	}

	i.Details[key] = val
}

func (i InventoryItem) GetHashableString() string {
	return fmt.Sprintf("%s%s%s%s%s%s", i.Manufacturer, i.Model, i.PartNumber, i.AssetTag, i.SerialNumber, i.Details)
}

func (i InventoryItem) Copy() *InventoryItem {
	return &InventoryItem{
		Type:         i.Type,
		Manufacturer: i.Manufacturer,
		Model:        i.Model,
		PartNumber:   i.PartNumber,
		SerialNumber: i.SerialNumber,
		AssetTag:     i.AssetTag,
		Details:      i.Details,
	}
}

//IsEqual compares an InventoryItem with another one
func (i InventoryItem) IsEqual(i2 InventoryItem) bool {
	if !utils.CompareStruct(i, i2, []string{}, []string{"Details"}) {
		return false
	}

	for key, val1 := range i.Details {
		if val1 != i2.Details[key] {
			return false
		}
	}

	return true
}

func (i InventoryItem) String() string {
	res, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(res)
}
