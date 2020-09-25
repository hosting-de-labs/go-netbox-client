package virtualization

import (
	"fmt"
	"net"
	"reflect"

	"github.com/go-openapi/swag"

	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

func (c Client) InterfaceConvertFromNetbox(netboxInterface models.VMInterface) (*types.NetworkInterface, error) {
	netIf := types.NewNetworkInterface()
	netIf.SetNetboxEntity(netboxInterface.ID, netboxInterface)

	netIf.Type = types.InterfaceTypeVirtualInterfacesVirtual

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
	}

	ipamClient := netboxIpam.NewClient(c.client)
	netboxAddresses, err := ipamClient.IPAddressFindByVMInterfaceID(netboxInterface.ID)
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
		addr.Family = types.IPAddressFamily(*netboxAddress.Family.Label)

		netIf.IPAddresses = append(netIf.IPAddresses, addr)
	}

	netIf.SetOriginalEntity(*netIf)

	return netIf, nil
}

//InterfaceConvertToNetbox allows to convert a NetworkInterface to a netbox compatible device interface
func (c Client) InterfaceConvertToNetbox(vmID int64, intf types.NetworkInterface) (out *models.WritableVMInterface, err error) {
	vm, err := c.VirtualMachineGet(vmID)
	if err != nil {
		return nil, err
	}

	var siteID int64
	switch vm.Meta.NetboxEntity.(type) {
	case models.VirtualMachineWithConfigContext:
		vm := vm.Meta.NetboxEntity.(models.VirtualMachineWithConfigContext)
		if vm.Site == nil {
			return nil, fmt.Errorf("vm with ID %d is not assigned to any site", vmID)
		}

		siteID = vm.Site.ID

	default:
		return nil, fmt.Errorf("Unsupported type for vm: %s", reflect.TypeOf(vm.Meta.NetboxEntity))
	}

	out = &models.WritableVMInterface{}

	out.VirtualMachine = &vmID
	out.Name = &intf.Name

	if intf.MACAddress != nil && intf.MACAddress.String() != "" {
		out.MacAddress = swag.String(intf.MACAddress.String())
	}

	for _, tag := range intf.Tags {
		out.Tags = append(out.Tags, &models.NestedTag{
			Name: &tag,
		})
	}

	if intf.UntaggedVlan != nil && len(intf.TaggedVlans) > 0 {
		//Tagged mode
		out.Mode = "tagged"
	} else if intf.UntaggedVlan != nil {
		//Access mode
		out.Mode = "access"
	} else if len(intf.TaggedVlans) > 0 {
		//All Tagged mode
		out.Mode = "tagged-all"
	}

	ipamClient := netboxIpam.NewClient(c.client)
	if intf.UntaggedVlan != nil {
		vlan, err := ipamClient.VLANGet(intf.UntaggedVlan.ID, &siteID)
		if err != nil {
			return nil, fmt.Errorf("an error occured when fetching vlan with tag %d: %s", intf.UntaggedVlan.ID, err)
		}

		if vlan == nil {
			return nil, fmt.Errorf("vlan %s with ID %d not found", intf.UntaggedVlan.Name, intf.UntaggedVlan.ID)
		}

		out.UntaggedVlan = &vlan.ID
	}

	out.TaggedVlans = []int64{}
	if len(intf.TaggedVlans) > 0 {
		for _, vlanTag := range intf.TaggedVlans {
			vlan, err := ipamClient.VLANGet(vlanTag.ID, &siteID)
			if err != nil {
				return nil, fmt.Errorf("an error occured when fetching vlan with tag %d: %s", vlanTag.ID, err)
			}

			out.TaggedVlans = append(out.TaggedVlans, vlan.ID)
		}
	}

	return out, nil
}
