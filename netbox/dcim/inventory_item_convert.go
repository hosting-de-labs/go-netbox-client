package dcim

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c *Client) InventoryItemConvertToNetbox(i types.InventoryItem, deviceID int64) (*models.WritableInventoryItem, error) {
	manufacturerName := i.Manufacturer

	if manufacturerName == "Undefined" || manufacturerName == "0000" {
		manufacturerName = "Unknown"
	}

	manufacturer, err := c.ManufacturerGet(i.Manufacturer)
	if err != nil {
		return nil, err
	}

	if manufacturer == nil {
		return nil, fmt.Errorf("manufacturer %s not found", manufacturerName)
	}

	model := i.Model
	if len(model) == 0 {
		model = "Unknown"
	}

	out := new(models.WritableInventoryItem)
	out.Manufacturer = &manufacturer.ID
	out.Name = swag.String(model)

	out.Device = swag.Int64(deviceID)

	out.PartID = i.PartNumber
	out.Serial = i.SerialNumber
	out.AssetTag = &i.AssetTag
	out.Discovered = true

	out.Tags = []string{}
	if i.Tags != nil {
		out.Tags = i.Tags
	}

	descriptionByte, err := json.Marshal(i.Details)
	if err != nil {
		return nil, err
	}

	description := string(descriptionByte)
	if description == "null" {
		description = ""
	}

	out.Description = description

	return out, nil
}

func (c *Client) InventoryItemConvertFromNetbox(i models.InventoryItem) (out types.InventoryItem) {
	out.OriginalEntity = &i
	out.Manufacturer = *i.Manufacturer.Name
	out.PartNumber = i.PartID
	out.SerialNumber = i.Serial
	out.AssetTag = *i.AssetTag

	result := strings.SplitN(*i.Name, ":", 1)
	switch len(result) {
	case 0, 1:
		out.Type = types.InventoryItemTypeOther
		out.Model = *i.Name
	default:
		out.Type, _ = types.InventoryItemTypeParse(result[0])
		out.Model = result[1]
	}

	err := json.Unmarshal([]byte(i.Description), &out.Details)
	if err != nil {
		out.AddDetail("JSON Unmarshal", "failed")
	}

	return out
}
