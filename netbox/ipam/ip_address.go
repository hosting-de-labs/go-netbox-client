package ipam

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/ipam"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//IPAddressGet returns an existing ip-address based on the given ip/cidr string.
func IPAddressGet(netboxClient *client.NetBox, ipWithCIDR string) (*models.IPAddress, error) {
	params := ipam.NewIPAMIPAddressesListParams()
	params.WithQ(&ipWithCIDR)

	res, err := netboxClient.IPAM.IPAMIPAddressesList(params, nil)

	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("IP Address %s is not unique", ipWithCIDR)
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
func IPAddressCreate(netboxClient *client.NetBox, ipWithCIDR string) (*models.IPAddress, error) {
	data := new(models.WritableIPAddress)
	data.Address = &ipWithCIDR
	data.Tags = []string{}

	params := ipam.NewIPAMIPAddressesCreateParams()
	params.WithData(data)

	_, err := netboxClient.IPAM.IPAMIPAddressesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return IPAddressGet(netboxClient, ipWithCIDR)
}

//IPAddressGetCreate is a convenience function that looks up an existing ip-address from netbox
//or creates it
func IPAddressGetCreate(netboxClient *client.NetBox, ipWithCIDR string) (*models.IPAddress, error) {
	res, err := IPAddressGet(netboxClient, ipWithCIDR)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	return IPAddressCreate(netboxClient, ipWithCIDR)
}

//IPAddressAssignInterface assigns a ip-address/cidr string to an existing interface.
func IPAddressAssignInterface(netboxClient *client.NetBox, ipWithCIDR string, deviceInterface models.Interface) (*models.IPAddress, error) {
	ipAddress, err := IPAddressGetCreate(netboxClient, ipWithCIDR)
	if err != nil {
		return nil, err
	}

	//Do not update ipAddress if interface is already correct
	if ipAddress.Interface != nil && ipAddress.Interface.ID == deviceInterface.ID {
		return ipAddress, nil
	}

	data := new(models.WritableIPAddress)
	data.Address = &ipWithCIDR
	data.Tags = []string{}
	data.Interface = deviceInterface.ID

	params := ipam.NewIPAMIPAddressesPartialUpdateParams()
	params.WithID(ipAddress.ID)
	params.WithData(data)

	_, err = netboxClient.IPAM.IPAMIPAddressesPartialUpdate(params, nil)
	if err != nil {
		return nil, err
	}

	return IPAddressGet(netboxClient, ipWithCIDR)
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
