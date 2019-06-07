package virtualization

import (
	"strings"

	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

func (c Client) InterfaceConvertFromNetbox(netboxInterface models.VirtualMachineInterface) (*types.HostNetworkInterface, error) {
	netIf := types.HostNetworkInterface{}

	if netboxInterface.Name != nil {
		netIf.Name = *netboxInterface.Name
	}

	if netboxInterface.MacAddress != nil {
		netIf.MACAddress = *netboxInterface.MacAddress
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
	}

	ipamClient := netboxIpam.NewClient(c.client)
	netboxAddresses, err := ipamClient.IPAddressGetByInterfaceID(netboxInterface.ID)
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

		addr.Type = types.IPAddressTypeIPv4
		if strings.Contains(ip, ":") {
			addr.Type = types.IPAddressTypeIPv6
		}

		netIf.IPAddresses = append(netIf.IPAddresses, addr)
	}

	return &netIf, nil
}
