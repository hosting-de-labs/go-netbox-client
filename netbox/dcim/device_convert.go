package dcim

import (
	"strings"

	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c Client) DeviceConvertFromNetbox(device *models.Device) (*types.DedicatedServer, error) {
	out := new(types.DedicatedServer)

	out.ID = device.ID
	out.Hostname = *device.Name
	out.Tags = device.Tags
	out.Comments = strings.Split(device.Comments, "\n")

	//iterate over tags to find managed tag
	for _, tag := range out.Tags {
		if tag == "managed" {
			out.IsManaged = true
			break
		}
	}

	if device.PrimaryIp4 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*device.PrimaryIp4.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv4 = types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv4,
		}
	}

	if device.PrimaryIp6 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*device.PrimaryIp6.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv6 = types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv6,
		}
	}

	out.OriginalEntity = out.Copy()

	return nil, nil
}
