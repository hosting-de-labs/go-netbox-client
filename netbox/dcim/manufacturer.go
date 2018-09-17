package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/netbox"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//ManufacturerGet retrieves a manufacturer from netbox
func (c *Client) ManufacturerGet(manufacturerName string) (*models.Manufacturer, error) {
	params := dcim.NewDcimManufacturersListParams().WithName(&manufacturerName)

	res, err := c.getClient().Dcim.DcimManufacturersList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("Manufacturer %s is not unique", manufacturerName)
	}

	//TODO: Use cache

	return res.Payload.Results[0], nil
}

//ManufacturerCreate creates a manufacturer in netbox
func (c *Client) ManufacturerCreate(manufacturerName string) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	var manufacturerSlug = netbox.GenerateSlug(manufacturerName)

	manufacturer.Name = &manufacturerName
	manufacturer.Slug = &manufacturerSlug

	params := dcim.NewDcimManufacturersCreateParams().WithData(&manufacturer)

	res, err := c.getClient().Dcim.DcimManufacturersCreate(params, nil)
	if err != nil {
		panic(err)
	}

	return res.Payload, nil
}
