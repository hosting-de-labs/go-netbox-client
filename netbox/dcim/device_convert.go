package dcim

import (
	"github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//TODO: DeviceConvertToNetbox

func (c *Client) DeviceConvertToNetboxObject(device *types.DedicatedServer, createUpdateObject bool) (*models.WritableDevice, error) {
	out := new(models.WritableDevice)

	if !createUpdateObject || (createUpdateObject && device.Hostname != device.OriginalHost.Hostname) {
		out.Name = device.Hostname
	}

	if !createUpdateObject || (createUpdateObject && !device.PrimaryIPv4.IsEqual(device.OriginalHost.PrimaryIPv4)) {
		ip, err := ipam.IPAddressGet(c.getClient(), device.PrimaryIPv4)
	}

	out.PrimaryIp4
	out.PrimaryIp6
	out.Serial
	out.Tags

	return out, nil
}

func (c *Client) DeviceConvertFromNetbox(device *models.Device) (*types.DedicatedServer, error) {
	out := new(types.DedicatedServer)

	out.ID = device.ID
	out.Hostname = device.Name
	out.Tags = device.Tags

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
			Type:    types.IPAddressTypeIPv4,
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
			Type:    types.IPAddressTypeIPv6,
		}
	}

	//TODO: Comments

	out.OriginalHost = out.Copy()

	return nil, nil
}
