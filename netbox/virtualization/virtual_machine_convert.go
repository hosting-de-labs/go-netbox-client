package virtualization

import (
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//VirtualMachineConvertFromNetbox converts a netbox virtual machine entity to a VirtualServer entity
func (c Client) VirtualMachineConvertFromNetbox(netboxVM models.VirtualMachine, interfaces []*models.VirtualMachineInterface) (out *types.VirtualServer, err error) {
	out = &types.VirtualServer{}
	out.Metadata.ID = netboxVM.ID
	out.Metadata.NetboxEntity = netboxVM

	out.Hostname = *netboxVM.Name

	if netboxVM.PrimaryIp4 != nil {
		address, cidr, err := utils.SplitCidrFromIP(*netboxVM.PrimaryIp4.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv4.Address = address
		out.PrimaryIPv4.CIDR = cidr
		out.PrimaryIPv4.Family = types.IPAddressFamilyIPv4
	}

	if netboxVM.PrimaryIp6 != nil {
		address, cidr, err := utils.SplitCidrFromIP(*netboxVM.PrimaryIp6.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv6.Address = address
		out.PrimaryIPv6.CIDR = cidr
		out.PrimaryIPv6.Family = types.IPAddressFamilyIPv6
	}

	if netboxVM.Vcpus != nil {
		out.Resources.Cores = int(*netboxVM.Vcpus)
	}

	if netboxVM.Memory != nil {
		out.Resources.Memory = *netboxVM.Memory
	}

	if netboxVM.Disk != nil {
		out.Resources.Disks = append(out.Resources.Disks, types.VirtualServerDisk{
			Size: *netboxVM.Disk * 1024,
		})
	}

	for _, tag := range netboxVM.Tags {
		out.AddTag(tag)

		if tag == "managed" {
			out.IsManaged = true
		}
	}

	if netboxVM.CustomFields != nil {
		customFields := utils.ConvertCustomFields(netboxVM.CustomFields)
		for key, val := range customFields {
			switch key {
			case "hypervisor_label":
				out.Hypervisor = val
			}
		}
	}

	//read comments
	utils.ParseVMComment(netboxVM.Comments, out)

	//interfaces / ips
	for _, netboxInterface := range interfaces {
		netIf, err := c.InterfaceConvertFromNetbox(*netboxInterface)
		if err != nil {
			return nil, err
		}

		out.NetworkInterfaces = append(out.NetworkInterfaces, *netIf)
	}

	out.OriginalEntity = out.Copy()
	return out, nil
}
