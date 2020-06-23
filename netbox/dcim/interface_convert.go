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
func (c Client) InterfaceConvertFromNetbox(nbIf models.DeviceInterface) (*types.NetworkInterface, error) {
	netIf := types.NewNetworkInterface()
	netIf.SetNetboxEntity(nbIf.ID, nbIf)

	if nbIf.Type != nil {
		netIf.Type = types.InterfaceType(*nbIf.Type.Value)
	}

	netIf.Enabled = nbIf.Enabled

	if nbIf.Name != nil {
		netIf.Name = *nbIf.Name
	}

	if nbIf.MacAddress != nil {
		mac, err := net.ParseMAC(*nbIf.MacAddress)
		if err != nil {
			return nil, err
		}

		netIf.MACAddress = mac
	}

	if nbIf.UntaggedVlan != nil {
		vlan, err := netboxIpam.VlanConvertFromNetbox(*nbIf.UntaggedVlan)
		if err != nil {
			return nil, err
		}

		netIf.UntaggedVlan = vlan
	}

	if len(nbIf.TaggedVlans) > 0 {
		for _, taggedVlan := range nbIf.TaggedVlans {
			vlan, err := netboxIpam.VlanConvertFromNetbox(*taggedVlan)
			if err != nil {
				return nil, err
			}

			netIf.TaggedVlans = append(netIf.TaggedVlans, *vlan)
		}
	}

	netIf.Tags = nbIf.Tags

	ipamClient := netboxIpam.NewClient(c.client)
	netboxAddresses, err := ipamClient.IPAddressFindByInterfaceID(nbIf.ID)
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

	netIf.SetOriginalEntity(*netIf)

	return netIf, nil
}

//InterfaceConvertToNetbox allows to convert a NetworkInterface to a netbox compatible device interface
func (c Client) InterfaceConvertToNetbox(deviceID int64, intf types.NetworkInterface) (out *models.WritableDeviceInterface, err error) {
	device, err := c.DeviceGet(deviceID)
	if err != nil {
		return nil, err
	}

	//d := device.Meta.GetNetboxEntity().DcimDeviceWithConfigContext()

	var siteID int64
	switch device.Meta.NetboxEntity.(type) {
	case models.Device:
		d := device.Meta.NetboxEntity.(models.Device)
		if d.Site == nil {
			return nil, fmt.Errorf("device with ID %d is not assigned to any site", deviceID)
		}

		siteID = d.Site.ID
	case models.DeviceWithConfigContext:
		d := device.Meta.NetboxEntity.(models.DeviceWithConfigContext)
		if d.Site == nil {
			return nil, fmt.Errorf("device with ID %d is not assigned to any site", deviceID)
		}

		siteID = d.Site.ID

	default:
		return nil, fmt.Errorf("Unsupported type for device: %s", reflect.TypeOf(device.Meta.NetboxEntity))
	}

	out = &models.WritableDeviceInterface{}

	out.Device = deviceID
	out.Name = intf.Name
	out.Type = string(intf.Type)
	out.Enabled = intf.Enabled
	out.MgmtOnly = intf.IsManagement

	if intf.MACAddress != nil && intf.MACAddress.String() != "" {
		out.MacAddress = swag.String(intf.MACAddress.String())
	}

	out.Tags = intf.Tags

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
