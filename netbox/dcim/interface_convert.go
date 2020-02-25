package dcim

import (
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

//InterfaceConvertFromNetbox allows to convert a DeviceInterface to a NetworkInterface
func (c Client) InterfaceConvertFromNetbox(netboxInterface models.DeviceInterface) (*types.NetworkInterface, error) {
	netIf := types.NetworkInterface{}

	//pass reference as original entity
	netIf.OriginalEntity = interface{}(&netboxInterface)

	if netboxInterface.Type != nil {
		ff := types.InterfaceType(*netboxInterface.Type.Value)
		netIf.Type = &ff
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

//InterfaceConvertToNetbox allows to convert a NetworkInterface to a netbox compatible device interface
func (c Client) InterfaceConvertToNetbox(deviceID int64, intf types.NetworkInterface) (out *models.WritableDeviceInterface, err error) {
	device, err := c.DeviceGet(deviceID)
	if err != nil {
		return nil, err
	}

	//d := device.Metadata.GetEntity().DcimDeviceWithConfigContext()

	var siteID int64
	switch device.Metadata.NetboxEntity.(type) {
	case models.Device:
		d := device.Metadata.NetboxEntity.(models.Device)
		if d.Site == nil {
			return nil, fmt.Errorf("device with ID %d is not assigned to any site", deviceID)
		}

		siteID = d.Site.ID
	case models.DeviceWithConfigContext:
		d := device.Metadata.NetboxEntity.(models.DeviceWithConfigContext)
		if d.Site == nil {
			return nil, fmt.Errorf("device with ID %d is not assigned to any site", deviceID)
		}

		siteID = d.Site.ID

	default:
		return nil, fmt.Errorf("Unsupported type for device: %s", reflect.TypeOf(device.Metadata.NetboxEntity))
	}

	out = &models.WritableDeviceInterface{}

	out.Device = &deviceID
	out.Name = &intf.Name

	if intf.Type != nil {
		out.Type = swag.String(string(*intf.Type))
	}

	out.MgmtOnly = intf.IsManagement
	out.MacAddress = swag.String(intf.MACAddress.String())

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
