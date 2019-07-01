package dcim

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c *Client) InventoryItemFindAll(deviceID int64) (out []types.InventoryItem, err error) {
	params := dcim.NewDcimInventoryItemsListParams()
	params.WithDeviceID(&deviceID).WithLimit(swag.Int64(100))

	res, err := c.client.Dcim.DcimInventoryItemsList(params, nil)
	if err != nil {
		return nil, err
	}

	for _, nbItem := range res.Payload.Results {
		out = append(out, c.InventoryItemConvertFromNetbox(*nbItem))
	}

	return out, nil
}

func (c *Client) InventoryItemCreate(deviceID int64, inventoryItem types.InventoryItem) (*models.InventoryItem, error) {
	params := dcim.NewDcimInventoryItemsCreateParams()

	data, err := c.InventoryItemConvertToNetbox(inventoryItem, deviceID)
	if err != nil {
		return nil, err
	}

	params.Data = data

	res, err := c.client.Dcim.DcimInventoryItemsCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

func (c *Client) InventoryItemDelete(id int64) error {
	params := dcim.NewDcimInventoryItemsDeleteParams().WithID(id)

	_, err := c.client.Dcim.DcimInventoryItemsDelete(params, nil)
	if err != nil {
		return err
	}

	return nil
}
