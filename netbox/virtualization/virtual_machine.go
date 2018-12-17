package virtualization

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/sirupsen/logrus"

	netboxDcim "github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
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

	return c.VMGet(vm.Hostname, false)
}

//VMGet retrieves an existing VM object from netbox by it's hostname.
func (c Client) VMGet(hostname string, deep bool) (*types.VirtualServer, error) {
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
		panic(err)
	}
	out := c.VirtualMachineConvertFromNetbox(vm, interfaces)

	return &out, nil
}

//VMGetCreate is a convenience wrapper for retrieving an existing VM object or creating it instead.
func (c Client) VMGetCreate(vm types.VirtualServer, hyp *models.Device) (*types.VirtualServer, error) {
	vmOut, err := c.VMGet(vm.Hostname, false)

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
func (c Client) VMUpdate(vm *types.VirtualServer, logger *logrus.Entry) bool {
	if !vm.IsChanged() {
		return false
	}

	dcimClient := netboxDcim.NewClient(c.client)
	ipamClient := netboxIpam.NewClient(c.client)

	hyp, err := dcimClient.HypervisorFindByHostname(vm.Hypervisor)
	if err != nil {
		panic(err)
	}

	//process network interfaces
	for _, netIf := range vm.NetworkInterfaces {
		vlan := new(models.VLAN)
		if netIf.VlanTag != 0 {
			vlan, err = ipamClient.VlanGet(netIf.VlanTag, hyp.Site.ID)
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

	//check if base data are not equal
	if !vm.IsEqual(vm.OriginalEntity.(types.VirtualServer), false) {
		data := new(models.WritableVirtualMachine)

		data.Name = swag.String(vm.Hostname)
		data.Tags = vm.Tags
		data.Cluster = &hyp.Cluster.ID
		data.Comments = utils.GenerateVMComment(vm)

		//custom fields
		customFields := make(map[string]string)
		customFields["hypervisor_label"] = vm.Hypervisor
		//FIXME: set hypervisor_url again
		//customFields["hypervisor_url"] = fmt.Sprintf("https://%s/dcim/devices/%d", config.netboxURL, hyp.ID)

		data.CustomFields = customFields

		//Resources
		data.Vcpus = swag.Int64(int64(vm.Resources.Cores))
		data.Memory = swag.Int64(vm.Resources.Memory)

		if len(vm.Resources.Disks) > 0 {
			data.Disk = swag.Int64(vm.Resources.Disks[0].Size / 1024)
		}

		ipamClient := netboxIpam.NewClient(c.client)

		//Primary IPs
		// we need this additional check due to a bug in netbox API
		if len(vm.PrimaryIPv4.Address) > 0 && vm.OriginalEntity.(types.VirtualServer).PrimaryIPv4.Address != vm.PrimaryIPv4.Address {
			ip4, err := ipamClient.IPAddressGet(vm.PrimaryIPv4)
			if err != nil {
				logger.WithError(err).Error("Cannot get ip address")
				return false
			}

			//set primary ip
			data.PrimaryIp4 = &ip4.ID
		}

		//we need this additional check due to a bug in netbox API
		if len(vm.PrimaryIPv6.Address) > 0 && vm.OriginalEntity.(types.VirtualServer).PrimaryIPv6.Address != vm.PrimaryIPv6.Address {
			ip6, err := ipamClient.IPAddressGet(vm.PrimaryIPv6)
			if err != nil {
				logger.WithError(err).Error("Cannot get ip address")
				return false
			}

			//set primary ip
			data.PrimaryIp6 = &ip6.ID
		}

		params := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams()
		params.WithID(vm.ID)
		params.WithData(data)

		_, err = c.client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(params, nil)
		if err != nil {
			logger.WithError(err).Error("Cannot update virtual machine")
			return false
		}
	}

	return true
}
