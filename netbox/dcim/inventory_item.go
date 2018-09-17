package dcim

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c *Client) InventoryItemGetAllByHostname(hostname string) ([]*models.InventoryItem, error) {
	params := dcim.NewDcimInventoryItemsListParams()
	params.Device = &hostname
	params.Limit = swag.Int64(100)

	res, err := c.getClient().Dcim.DcimInventoryItemsList(params, nil)
	if err != nil {
		return []*models.InventoryItem{}, err
	}

	return res.Payload.Results, nil
}

func (c *Client) InventoryItemCreate(deviceID int64, inventoryItem types.InventoryItem) (*models.InventoryItem, error) {
	params := dcim.NewDcimInventoryItemsCreateParams()

	data, err := c.InventoryItemConvertToNetbox(inventoryItem, deviceID)
	if err != nil {
		return nil, err
	}

	params.Data = data

	res, err := c.getClient().Dcim.DcimInventoryItemsCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

func (c *Client) InventoryItemDelete(id int64) error {
	params := dcim.NewDcimInventoryItemsDeleteParams().WithID(id)

	_, err := c.getClient().Dcim.DcimInventoryItemsDelete(params, nil)
	if err != nil {
		return err
	}

	return nil
}
