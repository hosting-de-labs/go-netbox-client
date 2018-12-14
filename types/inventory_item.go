package types

import (
	"encoding/json"
	"fmt"
	"strings"

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
	str := fmt.Sprintf("%s:%s:%s:%s:%s", i.Manufacturer, i.Model, i.PartNumber, i.AssetTag, i.SerialNumber)

	//TODO: order of map entries is not always stable
	if len(i.Details) > 0 {
		str = fmt.Sprintf("%s:details{", str)

		for key, val := range i.Details {
			str = fmt.Sprintf("%s%s:%s;", str, key, val)
		}

		str = fmt.Sprintf("%s}", str)
	}

	str = strings.Replace(str, " ", "", -1)

	return str
}

func (i InventoryItem) Copy() (out InventoryItem) {
	out = InventoryItem{
		Type:         i.Type,
		Manufacturer: i.Manufacturer,
		Model:        i.Model,
		PartNumber:   i.PartNumber,
		SerialNumber: i.SerialNumber,
		AssetTag:     i.AssetTag,
	}

	if len(i.Details) > 0 {
		out.Details = make(map[string]string, len(i.Details))
		for key, val := range i.Details {
			out.Details[key] = val
		}
	}

	return out
}

//IsEqual compares an InventoryItem with another one
func (i InventoryItem) IsEqual(i2 InventoryItem) bool {
	if !utils.CompareStruct(i, i2, []string{}, []string{"Details"}) {
		return false
	}

	for key, val1 := range i.Details {
		if val2, ok := i2.Details[key]; !ok || val1 != val2 {
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
