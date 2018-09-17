package ipam

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/ipam"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//IPAddressGet returns an existing ip-address based on the given ip/cidr string.
func IPAddressGet(netboxClient *client.NetBox, ipAddress types.IPAddress) (*models.IPAddress, error) {
	params := ipam.NewIPAMIPAddressesListParams()
	params.WithAddress(swag.String(ipAddress.String()))

	res, err := netboxClient.IPAM.IPAMIPAddressesList(params, nil)

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

//IPAddressGetByInterfaceID returns all interfaces assigned to an interface identified by it's ID
func IPAddressGetByInterfaceID(netboxClient *client.NetBox, interfaceID int64) ([]*models.IPAddress, error) {
	params := ipam.NewIPAMIPAddressesListParams()
	params.WithInterfaceID(swag.Int64(interfaceID))

	res, err := netboxClient.IPAM.IPAMIPAddressesList(params, nil)

	if err != nil {
		return nil, err
	}

	return res.Payload.Results, nil
}

//IPAddressCreate creates an ip-address based on the given ip/cidr string.
//Supports both IPv4 and IPv6.
func IPAddressCreate(netboxClient *client.NetBox, ipAddress *types.IPAddress) (*models.IPAddress, error) {
	data := new(models.WritableIPAddress)
	data.Address = swag.String(ipAddress.String())
	data.Tags = []string{}

	params := ipam.NewIPAMIPAddressesCreateParams()
	params.WithData(data)

	_, err := netboxClient.IPAM.IPAMIPAddressesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return IPAddressGet(netboxClient, ipAddress)
}

//IPAddressGetCreate is a convenience function that looks up an existing ip-address from netbox
//or creates it
func IPAddressGetCreate(netboxClient *client.NetBox, ipAddress *types.IPAddress) (*models.IPAddress, error) {
	res, err := IPAddressGet(netboxClient, ipAddress)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	return IPAddressCreate(netboxClient, ipAddress)
}

//IPAddressAssignInterface assigns a ip-address/cidr string to an existing interface.
func IPAddressAssignInterface(netboxClient *client.NetBox, ipAddress *types.IPAddress, deviceInterface models.Interface) (*models.IPAddress, error) {
	ipAddress2, err := IPAddressGetCreate(netboxClient, ipAddress)
	if err != nil {
		return nil, err
	}

	//Do not update ipAddress if interface is already correct
	if ipAddress2.Interface != nil && ipAddress2.Interface.ID == deviceInterface.ID {
		return ipAddress2, nil
	}

	data := new(models.WritableIPAddress)
	data.Address = swag.String(ipAddress.String())
	data.Tags = []string{}
	data.Interface = deviceInterface.ID

	params := ipam.NewIPAMIPAddressesPartialUpdateParams()
	params.WithID(ipAddress2.ID)
	params.WithData(data)

	_, err = netboxClient.IPAM.IPAMIPAddressesPartialUpdate(params, nil)
	if err != nil {
		return nil, err
	}

	return IPAddressGet(netboxClient, ipAddress)
}

//VlanGet returns a vlan-object based on the given vlanTag
func VlanGet(netboxClient *client.NetBox, vlanTag int64, siteID int64) (*models.VLAN, error) {
	params := ipam.NewIPAMVlansListParams()
	params.SetVid(&vlanTag)
	params.SetSiteID(&siteID)

	res, err := netboxClient.IPAM.IPAMVlansList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	return res.Payload.Results[0], nil
}
