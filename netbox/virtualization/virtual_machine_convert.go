package virtualization

import (
	"fmt"
	"reflect"

	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//VirtualMachineConvertFromNetbox converts a netbox virtual machine entity to a VirtualServer entity
func (c Client) VirtualMachineConvertFromNetbox(netboxVM interface{}, interfaces []*models.VirtualMachineInterface) (out *types.VirtualServer, err error) {
	out = types.NewVirtualServer()

	var cf interface{}
	primaryIPv4 := &models.NestedIPAddress{}
	primaryIPv6 := &models.NestedIPAddress{}

	switch netboxVM.(type) {
	case models.VirtualMachineWithConfigContext:
		vm := netboxVM.(models.VirtualMachineWithConfigContext)

		out.SetNetboxEntity(vm.ID, netboxVM)
		err = out.SetCustomFields(vm.CustomFields)
		if err != nil {
			return nil, err
		}

		out.Hostname = *vm.Name
		out.Tags = vm.Tags

		if vm.Vcpus != nil {
			out.Resources.Cores = int(*vm.Vcpus)
		}

		if vm.Memory != nil {
			out.Resources.Memory = *vm.Memory
		}

		if vm.Disk != nil {
			out.Resources.Disks = append(out.Resources.Disks, types.VirtualServerDisk{
				Size: *vm.Disk * 1024,
			})
		}

		primaryIPv4 = vm.PrimaryIp4
		primaryIPv6 = vm.PrimaryIp6
		cf = vm.CustomFields

		//read comments
		utils.ParseVMComment(vm.Comments, out)
	default:
		return nil, fmt.Errorf("unsupported type for device: %s", reflect.TypeOf(netboxVM))
	}

	for _, tag := range out.Tags {
		if tag == "managed" {
			out.IsManaged = true
			break
		}
	}

	if primaryIPv4 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*primaryIPv4.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv4 = &types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv4,
		}
	}

	if primaryIPv6 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*primaryIPv6.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv6 = &types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv6,
		}
	}

	if cf != nil {
		customFields := utils.ConvertCustomFields(cf)
		for key, val := range customFields {
			switch key {
			case "hypervisor_label":
				out.Hypervisor = val
			}
		}
	}

	if interfaces != nil {
		//interfaces / ips
		for _, netboxInterface := range interfaces {
			netIf, err := c.InterfaceConvertFromNetbox(*netboxInterface)
			if err != nil {
				return nil, err
			}

			out.NetworkInterfaces = append(out.NetworkInterfaces, *netIf)
		}
	}

	out.SetOriginalEntity(*out)

	return out, nil
}
