package virtualization

import (
	"fmt"

	"github.com/go-openapi/swag"
	netboxDcim "github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//VirtualMachineFindAll returns all found virtual machines
func (c Client) VirtualMachineFindAll(limit int64, offset int64) (int64, []*models.VirtualMachine, error) {
	params := virtualization.NewVirtualizationVirtualMachinesListParams()

	if limit > 0 {
		params.WithLimit(&limit)
	}

	if offset > 0 {
		params.WithOffset(&offset)
	}

	res, err := c.client.Virtualization.VirtualizationVirtualMachinesList(params, nil)

	if err != nil {
		return 0, nil, err
	}

	return *res.Payload.Count, res.Payload.Results, nil
}

//VMCreate creates a new VM object in Netbox.
func (c Client) VMCreate(clusterID int64, vm types.VirtualServer) (*types.VirtualServer, error) {
	var netboxVM models.WritableVirtualMachine
	netboxVM.Name = &vm.Hostname
	netboxVM.Tags = []string{}

	netboxVM.Cluster = &clusterID

	netboxVM.Vcpus = swag.Int64(int64(vm.Resources.Cores))
	netboxVM.Memory = swag.Int64(vm.Resources.Memory)

	if len(vm.Resources.Disks) > 0 {
		netboxVM.Disk = swag.Int64(vm.Resources.Disks[0].Size / 1024)
	}

	params := virtualization.NewVirtualizationVirtualMachinesCreateParams()
	params.WithData(&netboxVM)

	_, err := c.client.Virtualization.VirtualizationVirtualMachinesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.VMGet(vm.Hostname)
}

//VMGet retrieves an existing VM object from netbox by it's hostname.
func (c Client) VMGet(hostname string) (out *types.VirtualServer, err error) {
	params := virtualization.NewVirtualizationVirtualMachinesListParams()
	params.WithName(&hostname)

	res, err := c.client.Virtualization.VirtualizationVirtualMachinesList(params, nil)

	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("VM name %s not unique", hostname)
	}

	vm := *res.Payload.Results[0]

	interfaces, err := c.InterfaceGetAll(vm.ID)
	if err != nil {
		return nil, err
	}

	return c.VirtualMachineConvertFromNetbox(vm, interfaces)
}

//VMGetCreate is a convenience wrapper for retrieving an existing VM object or creating it instead.
func (c Client) VMGetCreate(clusterID int64, vm types.VirtualServer) (*types.VirtualServer, error) {
	vmOut, err := c.VMGet(vm.Hostname)

	if err != nil {
		return nil, err
	}

	if vmOut == nil {
		return c.VMCreate(clusterID, vm)
	}

	return vmOut, nil
}

//VMDelete deletes a virtual machine in Netbox
func (c Client) VMDelete(vmID int64) (err error) {
	params := virtualization.NewVirtualizationVirtualMachinesDeleteParams()
	params.SetID(vmID)

	_, err = c.client.Virtualization.VirtualizationVirtualMachinesDelete(params, nil)

	return err
}

//VMUpdate returns true if the vm was actually updated
func (c Client) VMUpdate(vm types.VirtualServer) (updated bool, err error) {
	if !vm.IsChanged() {
		return false, nil
	}

	dcimClient := netboxDcim.NewClient(c.client)

	hyp, err := dcimClient.HypervisorFindByHostname(vm.Hypervisor)
	if err != nil {
		return false, err
	}

	h := hyp.Metadata.NetboxEntity.(models.Device)

	//check if base data is equal
	if vm.IsEqual(vm.OriginalEntity.(types.VirtualServer), false) {
		_, err = c.updateInterfaces(vm, h)
		if err != nil {
			return false, err
		}
	}

	data := new(models.WritableVirtualMachine)

	data.Name = swag.String(vm.Hostname)
	data.Tags = vm.Tags
	data.Cluster = &h.Cluster.ID
	data.Comments = utils.GenerateVMComment(vm)

	//custom fields
	customFields := make(map[string]string)
	customFields["hypervisor_label"] = vm.Hypervisor

	data.CustomFields = customFields

	//Resources
	data.Vcpus = swag.Int64(int64(vm.Resources.Cores))
	data.Memory = swag.Int64(vm.Resources.Memory)

	if len(vm.Resources.Disks) > 0 {
		data.Disk = swag.Int64(vm.Resources.Disks[0].Size / 1024)
	}

	//Primary IPs
	if len(vm.PrimaryIPv4.Address) > 0 && vm.OriginalEntity.(types.VirtualServer).PrimaryIPv4.Address != vm.PrimaryIPv4.Address {
		IPID, err := c.preparePrimaryIPAddress(vm.PrimaryIPv4)
		if err != nil {
			return false, err
		}

		data.PrimaryIp4 = &IPID
	}

	if len(vm.PrimaryIPv6.Address) > 0 && vm.OriginalEntity.(types.VirtualServer).PrimaryIPv6.Address != vm.PrimaryIPv6.Address {
		IPID, err := c.preparePrimaryIPAddress(vm.PrimaryIPv6)
		if err != nil {
			return false, err
		}

		data.PrimaryIp6 = &IPID
	}

	//we need to update interfaces before we possibly assign new primary ip addresses
	//otherwise netbox might complain about ip addresses not being assigned to a virtual machine
	u, err := c.updateInterfaces(vm, h)
	if err != nil {
		return false, err
	}

	if u != false {
		updated = true
	}

	params := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams()
	params.WithID(vm.Metadata.NetboxEntity.(models.VirtualMachine).ID)
	params.WithData(data)

	_, err = c.client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(params, nil)
	if err != nil {
		return false, err
	}

	return updated, nil
}

//preparePrimaryIPAddress is a helper method to clear a possible primary ip address assignment before assigning an ip
//address to a different vm
func (c Client) preparePrimaryIPAddress(primaryIP types.IPAddress) (int64, error) {
	ipamClient := netboxIpam.NewClient(c.client)

	ip, err := ipamClient.IPAddressFind(primaryIP)
	if err != nil {
		return -1, err
	}

	if ip != nil {
		err = ipamClient.IPAddressDelete(ip.ID)
		if err != nil {
			return -1, err
		}
	}

	ip, err = ipamClient.IPAddressFindCreate(primaryIP)
	if err != nil {
		return -1, err
	}

	return ip.ID, nil
}

func (c Client) updateInterfaces(vm types.VirtualServer, hyp models.Device) (updated bool, err error) {
	ipamClient := netboxIpam.NewClient(c.client)

	//process network interfaces
	for _, netIf := range vm.NetworkInterfaces {
		vlan := new(models.VLAN)
		if netIf.UntaggedVlan != nil {
			vlan, err = ipamClient.VLANGet(netIf.UntaggedVlan.ID, &hyp.Site.ID)

			if err != nil {
				return false, err
			}
		}

		vmInterface, err := c.InterfaceGetCreate(vm, netIf.Name, vlan)
		if err != nil {
			return false, err
		}

		for _, network := range netIf.IPAddresses {
			_, err = ipamClient.IPAddressAssignInterface(network, vmInterface.ID)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
