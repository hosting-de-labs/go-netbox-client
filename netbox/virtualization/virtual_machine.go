package virtualization

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/sirupsen/logrus"

	netboxDcim "github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	netboxUtils "github.com/hosting-de-labs/go-netbox-client/netbox/utils"
)

//VMCreate creates a new VM object in Netbox.
func VMCreate(netboxClient *client.NetBox, vm *types.VirtualServer, hyp *models.Device) (*types.VirtualServer, error) {
	var netboxVM models.WritableVirtualMachine
	netboxVM.Name = &vm.Hostname
	netboxVM.Tags = []string{}

	netboxVM.Cluster = hyp.Cluster.ID

	netboxVM.Vcpus = swag.Int64(int64(vm.Resources.Cores))
	netboxVM.Memory = swag.Int64(vm.Resources.Memory)

	if len(vm.Resources.Disks) > 0 {
		netboxVM.Disk = swag.Int64(vm.Resources.Disks[0].Size / 1024)
	}

	params := virtualization.NewVirtualizationVirtualMachinesCreateParams()
	params.WithData(&netboxVM)

	_, err := netboxClient.Virtualization.VirtualizationVirtualMachinesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return VMGet(netboxClient, vm.Hostname, false)
}

//VMGet retrieves an existing VM object from netbox by it's hostname.
func VMGet(netboxClient *client.NetBox, hostname string, deep bool) (*types.VirtualServer, error) {
	params := virtualization.NewVirtualizationVirtualMachinesListParams()
	params.WithName(&hostname)

	res, err := netboxClient.Virtualization.VirtualizationVirtualMachinesList(params, nil)

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

	interfaces, err := InterfaceGetAll(netboxClient, vm.ID)
	if err != nil {
		panic(err)
	}
	out := netboxUtils.ConvertVMToVirtualServer(netboxClient, vm, interfaces)

	return &out, nil
}

//VMGetCreate is a convenience wrapper for retrieving an existing VM object or creating it instead.
func VMGetCreate(netboxClient *client.NetBox, vm *types.VirtualServer, hyp *models.Device) (*types.VirtualServer, error) {
	vmOut, err := VMGet(netboxClient, vm.Hostname, false)

	if err != nil {
		return nil, err
	}

	if vmOut == nil {
		return VMCreate(netboxClient, vm, hyp)
	}

	return vmOut, nil
}

//VMDelete deletes a virtual machine in Netbox
func VMDelete(netboxClient *client.NetBox, vmID int64) error {
	params := virtualization.NewVirtualizationVirtualMachinesDeleteParams()
	params.SetID(vmID)

	_, err := netboxClient.Virtualization.VirtualizationVirtualMachinesDelete(params, nil)
	if err != nil {
		return err
	}

	return nil
}

//VMUpdate returns true if the vm was actually updated
func VMUpdate(netboxClient *client.NetBox, vm *types.VirtualServer, logger *logrus.Entry) bool {
	if !vm.IsChanged() {
		return false
	}

	hyp, err := netboxDcim.HypervisorGet(netboxClient, vm.Hypervisor)
	if err != nil {
		panic(err)
	}

	//process network interfaces
	for _, netIf := range vm.NetworkInterfaces {
		vlan := new(models.VLAN)
		if netIf.VlanTag != 0 {
			vlan, err = netboxIpam.VlanGet(netboxClient, netIf.VlanTag, hyp.Site.ID)
		}

		vmInterface, err := InterfaceGetCreate(netboxClient, vm, netIf.Name, vlan)
		if err != nil {
			logger.WithError(err).Error("Cannot create interface")
			continue
		}

		for _, network := range netIf.IPAddresses {
			ipWithCIDR := fmt.Sprintf("%s/%d", network.Address, network.CIDR)

			_, err = netboxIpam.IPAddressAssignInterface(netboxClient, ipWithCIDR, *vmInterface)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"host":  vm.Hostname,
					"error": err,
					"ip":    ipWithCIDR,
				}).Error("Cannot assign ip to interface")
				continue
			}
		}
	}

	//check if base data are not equal
	if !vm.IsEqual(*vm.OriginalHost, false) {
		data := new(models.WritableVirtualMachine)

		data.Name = swag.String(vm.Hostname)
		data.Tags = vm.Tags
		data.Cluster = hyp.Cluster.ID
		data.Comments = ""

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

		//Primary IPs
		// we need this additional check due to a bug in netbox API
		if len(vm.PrimaryIPv4.Address) > 0 && vm.OriginalHost.PrimaryIPv4.Address != vm.PrimaryIPv4.Address {
			ipWithCIDR := fmt.Sprintf("%s/%d", vm.PrimaryIPv4.Address, vm.PrimaryIPv4.CIDR)
			ip4, err := netboxIpam.IPAddressGet(netboxClient, ipWithCIDR)
			if err != nil {
				logger.WithError(err).Error("Cannot get ip address")
				return false
			}

			//set primary ip
			data.PrimaryIp4 = ip4.ID
		}

		//we need this additional check due to a bug in netbox API
		if len(vm.PrimaryIPv6.Address) > 0 && vm.OriginalHost.PrimaryIPv6.Address != vm.PrimaryIPv6.Address {
			ip6, err := netboxIpam.IPAddressGet(netboxClient, fmt.Sprintf("%s/%d", vm.PrimaryIPv6.Address, vm.PrimaryIPv6.CIDR))
			if err != nil {
				logger.WithError(err).Error("Cannot get ip address")
				return false
			}

			//set primary ip
			data.PrimaryIp6 = ip6.ID
		}

		params := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams()

		params.WithID(vm.ID)
		params.WithData(data)

		_, err = netboxClient.Virtualization.VirtualizationVirtualMachinesPartialUpdate(params, nil)
		if err != nil {
			logger.WithError(err).Error("Cannot update virtual machine")
			return false
		}
	}

	return true
}
