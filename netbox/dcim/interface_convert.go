package dcim

import (
	"net"
	"strings"

	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

func (c Client) InterfaceConvertFromNetbox(netboxInterface models.DeviceInterface) (*types.NetworkInterface, error) {
	netIf := types.NetworkInterface{}

	//pass reference as original entity
	netIf.OriginalEntity = interface{}(&netboxInterface)

	if netboxInterface.FormFactor != nil {
		netIf.FormFactor = types.InterfaceFormFactor(*netboxInterface.FormFactor.Value)
	}

	if netboxInterface.Name != nil {
		netIf.Name = *netboxInterface.Name
	}

	if netboxInterface.MacAddress != nil {
		mac, err := net.ParseMAC(*netboxInterface.MacAddress)
		if err != nil {
			return nil, err
		}

		netIf.MACAddress = mac
	}

	if netboxInterface.UntaggedVlan != nil {
		vlan, err := netboxIpam.VlanConvertFromNetbox(*netboxInterface.UntaggedVlan)
		if err != nil {
			return nil, err
		}

		netIf.UntaggedVlan = vlan
	}

	if len(netboxInterface.TaggedVlans) > 0 {
		for _, taggedVlan := range netboxInterface.TaggedVlans {
			vlan, err := netboxIpam.VlanConvertFromNetbox(*taggedVlan)
			if err != nil {
				return nil, err
			}

			netIf.TaggedVlans = append(netIf.TaggedVlans, *vlan)
		}
	} else {
		netboxInterface.TaggedVlans = []*models.NestedVLAN{}
	}

	ipamClient := netboxIpam.NewClient(c.client)
	netboxAddresses, err := ipamClient.IPAddressFindByInterfaceID(netboxInterface.ID)
	if err != nil {
		return nil, err
	}

	for _, netboxAddress := range netboxAddresses {
		var addr types.IPAddress
		ip, cidr, err := utils.SplitCidrFromIP(*netboxAddress.Address)
		if err != nil {
			return nil, err
		}

		addr.Address = ip
		addr.CIDR = cidr

		addr.Family = types.IPAddressFamilyIPv4
		if strings.Contains(ip, ":") {
			addr.Family = types.IPAddressFamilyIPv6
		}

		netIf.IPAddresses = append(netIf.IPAddresses, addr)
	}

	return &netIf, nil
}
