package virtualization

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
)

//VirtualMachineCreate creates a new VM object in Netbox.
func (c Client) VirtualMachineCreate(clusterID int64, vm types.VirtualServer) (*types.VirtualServer, error) {
	var netboxVM models.WritableVirtualMachineWithConfigContext
	netboxVM.Name = vm.Hostname
	netboxVM.Tags = []string{}

	netboxVM.Cluster = clusterID

	netboxVM.Vcpus = swag.Int64(int64(vm.Resources.Cores))
	netboxVM.Memory = swag.Int64(vm.Resources.Memory)

	if len(vm.Resources.Disks) > 0 {
		netboxVM.Disk = swag.Int64(vm.Resources.Disks[0].Size / 1024)
	}

	//Custom fields
	customFields := make(map[string]string)
	customFields["hypervisor_label"] = vm.Hypervisor
	netboxVM.CustomFields = customFields

	params := virtualization.NewVirtualizationVirtualMachinesCreateParams()
	params.WithData(&netboxVM)

	res, err := c.client.Virtualization.VirtualizationVirtualMachinesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.VirtualMachineGet(res.Payload.ID)
}

//VirtualMachineDelete deletes a virtual machine in Netbox
func (c Client) VirtualMachineDelete(vmID int64) (err error) {
	params := virtualization.NewVirtualizationVirtualMachinesDeleteParams()
	params.SetID(vmID)

	_, err = c.client.Virtualization.VirtualizationVirtualMachinesDelete(params, nil)

	return err
}

//VirtualMachineFindAll returns all found virtual machines
func (c Client) VirtualMachineFindAll(limit int64, offset int64) (int64, []*models.VirtualMachineWithConfigContext, error) {
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

//VirtualMachineFind returns the first found virtual machines
func (c Client) VirtualMachineFind(hostname string) (out *types.VirtualServer, err error) {
	params := virtualization.NewVirtualizationVirtualMachinesListParams()
	params.WithName(&hostname)

	res, err := c.client.Virtualization.VirtualizationVirtualMachinesList(params, nil)
	if err != nil {
		return nil, err
	}

	if len(res.Payload.Results) == 0 {
		return nil, nil
	}

	return c.VirtualMachineConvertFromNetbox(*res.Payload.Results[0], nil)
}

//VirtualMachineGet retrieves an existing VM object from netbox by it's hostname.
func (c Client) VirtualMachineGet(vmID int64) (out *types.VirtualServer, err error) {
	params := virtualization.NewVirtualizationVirtualMachinesReadParams()
	params.WithID(vmID)

	res, err := c.client.Virtualization.VirtualizationVirtualMachinesRead(params, nil)

	if err != nil {
		return nil, err
	}

	interfaces, err := c.InterfaceFindAll(vmID)
	if err != nil {
		return nil, err
	}

	return c.VirtualMachineConvertFromNetbox(*res.Payload, interfaces)
}

//VirtualMachineGetCreate is a convenience wrapper for retrieving an existing VM object or creating it instead.
func (c Client) VirtualMachineGetCreate(clusterID int64, vm types.VirtualServer) (*types.VirtualServer, error) {
	vmOut, err := c.VirtualMachineFind(vm.Hostname)

	if err != nil {
		return nil, err
	}

	if vmOut == nil {
		return c.VirtualMachineCreate(clusterID, vm)
	}

	return vmOut, nil
}

//VirtualMachineUpdate returns true if the vm was actually updated
func (c Client) VirtualMachineUpdate(vm types.VirtualServer) (updated bool, err error) {
	if !vm.IsChanged() {
		return false, nil
	}

	res, err := c.VirtualMachineFind(vm.Hostname)
	if err != nil {
		return false, err
	}

	nbVM := res.Meta.NetboxEntity.(models.VirtualMachineWithConfigContext)

	//check if base data is equal and only update interfaces
	origVirtualServer, ok := vm.Meta.OriginalEntity.(types.VirtualServer)
	if ok {
		if vm.IsEqual(origVirtualServer, false) {
			return c.updateInterfaces(vm)
		}
	}

	data := new(models.WritableVirtualMachineWithConfigContext)

	data.Name = vm.Hostname
	data.Tags = vm.Tags
	data.Cluster = nbVM.Cluster.ID
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

	//compare ipv4 / v6 addresses and sets ids of writable object if not nil
	u, err := c.compareIPAddresses(vm, origVirtualServer, data)
	if err != nil {
		return false, err
	}

	if u != false {
		updated = true
	}

	//we need to update interfaces before we possibly assign new primary ip addresses
	//otherwise netbox might complain about ip addresses not being assigned to a virtual machine
	u, err = c.updateInterfaces(vm)
	if err != nil {
		return false, err
	}

	if u != false {
		updated = true
	}

	params := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams()
	params.WithID(vm.Meta.NetboxEntity.(models.VirtualMachineWithConfigContext).ID)
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

func (c Client) updateInterfaces(vm types.VirtualServer) (updated bool, err error) {
	ipamClient := netboxIpam.NewClient(c.client)

	//process network interfaces
	for _, netIf := range vm.NetworkInterfaces {
		vmInterface, err := c.InterfaceGetCreate(vm.Meta.ID, netIf)
		if err != nil {
			return false, err
		}

		for _, network := range netIf.IPAddresses {
			_, err = ipamClient.IPAddressAssignInterface(network, vmInterface.Meta.ID)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (c Client) compareIPAddresses(vm1 types.VirtualServer, vm2 types.VirtualServer, updateObject *models.WritableVirtualMachineWithConfigContext) (updated bool, err error) {
	//Primary IPs
	if vm1.PrimaryIPv4 != nil && vm1.PrimaryIPv4.Address != vm2.PrimaryIPv4.Address {
		ipv4Id, err := c.preparePrimaryIPAddress(*vm1.PrimaryIPv4)
		if err != nil {
			return updated, err
		}

		if updateObject != nil {
			updateObject.PrimaryIp4 = &ipv4Id
			updated = true
		}
	}

	if vm1.PrimaryIPv6 != nil && vm1.PrimaryIPv6.Address != vm2.PrimaryIPv6.Address {
		ipv6Id, err := c.preparePrimaryIPAddress(*vm1.PrimaryIPv6)
		if err != nil {
			return updated, err
		}

		if updateObject != nil {
			updateObject.PrimaryIp6 = &ipv6Id
			updated = true
		}
	}

	return updated, nil
}
