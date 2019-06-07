package ipam

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/client/ipam"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//VlanGet returns a vlan-object based on the given vlanID and an optional siteID
func (c Client) VLANGet(vlanID uint16, siteID *int64) (*models.VLAN, error) {
	params := ipam.NewIPAMVlansListParams()
	params.SetVid(swag.Int64(int64(vlanID)))

	if siteID != nil {
		params.SetSiteID(siteID)
	}

	res, err := c.client.IPAM.IPAMVlansList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	return res.Payload.Results[0], nil
}

//VlanGet returns a vlan-object based on the given vlanTag and siteID
//
// Deprecated: use VLANGet instead
func (c Client) VlanGet(vlanTag int64, siteID int64) (*models.VLAN, error) {
	return c.VLANGet(uint16(vlanTag), &siteID)
}
