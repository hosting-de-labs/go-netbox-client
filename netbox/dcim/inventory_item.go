package dcim

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func InventoryItemGetByHostname(netboxClient *client.NetBox, hostname string) ([]*models.InventoryItem, error) {
	params := dcim.NewDcimInventoryItemsListParams()
	params.Device = &hostname

	res, err := netboxClient.Dcim.DcimInventoryItemsList(params, nil)
	if err != nil {
		return []*models.InventoryItem{}, err
	}

	return res.Payload.Results, nil
}

func InventoryItemCreate(netboxClient *client.NetBox) {
	params := dcim.NewDcimInventoryItemsCreateParams()

	manufacturer := item.GetManufacturer()
	model := item.GetModel()

	if manufacturer == "Undefined" || manufacturer == "0000" {
		manufacturer = "Unknown"
	}

	if len(model) == 0 {
		model = "Unknown"
	}

	params.Data = &models.WritableInventoryItem{
		Device:       &device.ID,
		Discovered:   true,
		Manufacturer: swag.Int64(getManufacturer(manufacturer).ID),
		Name:         &model,
		Serial:       item.GetSerialNumber(),
		AssetTag:     GenerateItemHash(item),
		Description:  item.GetDescription(),
		Tags:         []string{},
	}

	_, err := netbox.Dcim.DcimInventoryItemsCreate(params, nil)
	if err != nil {
		panic(err)
	}
}
