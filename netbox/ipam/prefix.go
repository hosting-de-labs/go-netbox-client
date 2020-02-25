package ipam

import (
	"github.com/hosting-de-labs/go-netbox/netbox/client/ipam"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c Client) PrefixGet(prefix string) (*models.Prefix, error) {
	params := ipam.NewIpamPrefixesListParams()
	params.WithQ(&prefix)

	res, err := c.client.Ipam.IpamPrefixesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	return res.Payload.Results[0], nil
}
