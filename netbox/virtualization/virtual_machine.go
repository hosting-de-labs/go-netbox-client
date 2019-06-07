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
	"github.com/sirupsen/logrus"
)

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
func (c Client) VMCreate(vm types.VirtualServer, hyp *models.Device) (*types.VirtualServer, error) {
	var netboxVM models.WritableVirtualMachine
	netboxVM.Name = &vm.Hostname
	netboxVM.Tags = []string{}

	netboxVM.Cluster = &hyp.Cluster.ID

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
func (c Client) VMGet(hostname string) (*types.VirtualServer, error) {
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
	out, err := c.VirtualMachineConvertFromNetbox(vm, interfaces)

	return out, nil
}

//VMGetCreate is a convenience wrapper for retrieving an existing VM object or creating it instead.
func (c Client) VMGetCreate(vm types.VirtualServer, hyp *models.Device) (*types.VirtualServer, error) {
	vmOut, err := c.VMGet(vm.Hostname)

	if err != nil {
		return nil, err
	}

	if vmOut == nil {
		return c.VMCreate(vm, hyp)
	}

	return vmOut, nil
}

//VMDelete deletes a virtual machine in Netbox
func (c Client) VMDelete(vmID int64) error {
	params := virtualization.NewVirtualizationVirtualMachinesDeleteParams()
	params.SetID(vmID)

	_, err := c.client.Virtualization.VirtualizationVirtualMachinesDelete(params, nil)
	if err != nil {
		return err
	}

	return nil
}

//VMUpdate returns true if the vm was actually updated
func (c Client) VMUpdate(vm *types.VirtualServer, logger *logrus.Entry) (bool, error) {
	if !vm.IsChanged() {
		return false, nil
	}

	dcimClient := netboxDcim.NewClient(c.client)

	hyp, err := dcimClient.HypervisorFindByHostname(vm.Hypervisor)
	if err != nil {
		return false, err
	}

	//check if base data is equal
	if vm.IsEqual(vm.OriginalEntity.(types.VirtualServer), false) {
		c.updateInterfaces(vm, hyp, logger)
		return true, nil
	}

	data := new(models.WritableVirtualMachine)

	data.Name = swag.String(vm.Hostname)
	data.Tags = vm.Tags
	data.Cluster = &hyp.Cluster.ID
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
		IPID, err := c.preparePrimaryIpAddress(vm.PrimaryIPv4)
		if err != nil {
			return false, err
		}

		data.PrimaryIp4 = &IPID
	}

	if len(vm.PrimaryIPv6.Address) > 0 && vm.OriginalEntity.(types.VirtualServer).PrimaryIPv6.Address != vm.PrimaryIPv6.Address {
		IPID, err := c.preparePrimaryIpAddress(vm.PrimaryIPv6)
		if err != nil {
			return false, err
		}

		data.PrimaryIp6 = &IPID
	}

	//we need to update interfaces before we possibly assign new primary ip addresses
	//otherwise netbox might complain about ip addresses not being assigned to a virtual machine
	c.updateInterfaces(vm, hyp, logger)

	params := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams()
	params.WithID(vm.ID)
	params.WithData(data)

	_, err = c.client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(params, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

//preparePrimaryIpAddress is a helper method to clear a possible primary ip address assignment before assigning an ip
//address to a different vm
func (c Client) preparePrimaryIpAddress(primaryIP types.IPAddress) (int64, error) {
	ipamClient := netboxIpam.NewClient(c.client)

	ip, err := ipamClient.IPAddressGet(primaryIP)
	if err != nil {
		return -1, err
	}

	if ip != nil {
		err = ipamClient.IPAddressDelete(ip.ID)
		if err != nil {
			return -1, err
		}
	}

	ip, err = ipamClient.IPAddressGetCreate(primaryIP)
	if err != nil {
		return -1, err
	}

	return ip.ID, nil
}

func (c Client) updateInterfaces(vm *types.VirtualServer, hyp *models.Device, logger *logrus.Entry) {
	ipamClient := netboxIpam.NewClient(c.client)

	//process network interfaces
	for _, netIf := range vm.NetworkInterfaces {
		vlan := new(models.VLAN)
		if netIf.UntaggedVlan != nil {
			var err error
			vlan, err = ipamClient.VLANGet(netIf.UntaggedVlan.ID, &hyp.Site.ID)

			if err != nil {
				logger.WithError(err).Errorln("Cannot get VLAN")
				continue
			}
		}

		vmInterface, err := c.InterfaceGetCreate(vm, netIf.Name, vlan)
		if err != nil {
			logger.WithError(err).Error("Cannot create interface")
			continue
		}

		for _, network := range netIf.IPAddresses {
			_, err = ipamClient.IPAddressAssignInterface(network, vmInterface.ID)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"host":  vm.Hostname,
					"error": err,
					"ip":    fmt.Sprintf("%s/%d", network.Address, network.CIDR),
				}).Error("Cannot assign ip to interface")
				continue
			}
		}
	}
}
