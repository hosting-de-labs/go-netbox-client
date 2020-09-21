package ipam

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/ipam"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//IPAddressGet returns an existing ip-address based on the given ip/cidr string.
func (c Client) IPAddressFind(ipAddress types.IPAddress) (*models.IPAddress, error) {
	params := ipam.NewIpamIPAddressesListParams()
	params.WithAddress(swag.String(ipAddress.String()))

	res, err := c.client.Ipam.IpamIPAddressesList(params, nil)

	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("IP Address %s is not unique", ipAddress)
	}

	return res.Payload.Results[0], nil
}

//IPAddressGetByInterfaceID returns all ip addresses assigned to an interface identified by it's ID
func (c Client) IPAddressFindByInterfaceID(interfaceID int64) ([]*models.IPAddress, error) {
	params := ipam.NewIpamIPAddressesListParams()
	params.WithInterfaceID(swag.Int64(interfaceID))

	res, err := c.client.Ipam.IpamIPAddressesList(params, nil)

	if err != nil {
		panic(err)
		return nil, err
	}

	return res.Payload.Results, nil
}

//IPAddressCreate creates an ip-address based on the given ip/cidr string.
func (c Client) IPAddressCreate(ipAddress types.IPAddress) (*models.IPAddress, error) {
	data := new(models.WritableIPAddress)
	data.Address = swag.String(ipAddress.String())

	params := ipam.NewIpamIPAddressesCreateParams()
	params.WithData(data)

	_, err := c.client.Ipam.IpamIPAddressesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.IPAddressFind(ipAddress)
}

//IPAddressGetCreate is a convenience function that looks up an existing ip-address from netbox or creates it
func (c Client) IPAddressFindCreate(ipAddress types.IPAddress) (*models.IPAddress, error) {
	res, err := c.IPAddressFind(ipAddress)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	return c.IPAddressCreate(ipAddress)
}

//IPAddressAssignInterface assigns a ip-address/cidr string to an existing interface.
func (c Client) IPAddressAssignInterface(ipAddress types.IPAddress, interfaceID int64) (*models.IPAddress, error) {
	ipAddress2, err := c.IPAddressFindCreate(ipAddress)
	if err != nil {
		return nil, err
	}

	//Do not update ipAddress if interface is already correct
	if *ipAddress2.AssignedObjectID == interfaceID {
		return ipAddress2, nil
	}

	data := new(models.WritableIPAddress)
	data.Address = swag.String(ipAddress.String())
	data.AssignedObjectID = &interfaceID

	params := ipam.NewIpamIPAddressesPartialUpdateParams()
	params.WithID(ipAddress2.ID)
	params.WithData(data)

	_, err = c.client.Ipam.IpamIPAddressesPartialUpdate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.IPAddressFind(ipAddress)
}

func (c Client) IPAddressDelete(ipAddressID int64) error {
	p := ipam.NewIpamIPAddressesDeleteParams()
	p.WithID(ipAddressID)

	_, err := c.client.Ipam.IpamIPAddressesDelete(p, nil)

	return err
}
