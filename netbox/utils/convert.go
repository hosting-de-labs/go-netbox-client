package utils

import (
	"strings"

	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	"internal.keenlogics.com/di/netbox-sync/helper"
	"internal.keenlogics.com/di/netbox-sync/netbox/ipam"
	"internal.keenlogics.com/di/netbox-sync/types"
)

func ConvertVMToVirtualServer(netboxClient *client.NetBox, netboxVM models.VirtualMachine, interfaces []*models.Interface) types.VirtualServer {
	var out types.VirtualServer
	out.ID = netboxVM.ID
	out.Hostname = *netboxVM.Name

	if netboxVM.PrimaryIp4 != nil {
		address, cidr, err := helper.SplitCidrFromIP(*netboxVM.PrimaryIp4.Address)
		if err != nil {
			panic(err)
		}

		out.PrimaryIPv4.Address = address
		out.PrimaryIPv4.CIDR = cidr
		out.PrimaryIPv4.Type = types.IPAddressTypeIPv4
	}

	if netboxVM.PrimaryIp6 != nil {
		address, cidr, err := helper.SplitCidrFromIP(*netboxVM.PrimaryIp6.Address)
		if err != nil {
			panic(err)
		}

		out.PrimaryIPv6.Address = address
		out.PrimaryIPv6.CIDR = cidr
		out.PrimaryIPv6.Type = types.IPAddressTypeIPv6
	}

	out.Resources.Cores = int(*netboxVM.Vcpus)
	out.Resources.Memory = *netboxVM.Memory

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

	customFields := ConvertCustomFields(netboxVM.CustomFields)
	for key, val := range customFields {
		switch key {
		case "hypervisor_label":
			out.Hypervisor = val
		}
	}

	//TODO: import additional disks from comments section

	//interfaces / ips
	for _, netboxInterface := range interfaces {
		var netIf types.HostNetworkInterface
		netIf.Name = *netboxInterface.Name
		//TODO: netIf.VlanTag = netboxInterface.UntaggedVlan
		netIf.MACAddress = netboxInterface.MacAddress

		netboxAddresses, err := ipam.IPAddressGetByInterfaceID(netboxClient, netboxInterface.ID)
		if err != nil {
			panic(err)
		}

		for _, netboxAddress := range netboxAddresses {
			var addr types.IPAddress
			ip, cidr, err := helper.SplitCidrFromIP(*netboxAddress.Address)
			if err != nil {
				panic(err)
			}

			addr.Address = ip
			addr.CIDR = cidr

			addr.Type = types.IPAddressTypeIPv4
			if strings.Contains(ip, ":") {
				addr.Type = types.IPAddressTypeIPv6
			}

			netIf.IPAddresses = append(netIf.IPAddresses, addr)
		}

		out.NetworkInterfaces = append(out.NetworkInterfaces, netIf)
	}

	out.OriginalHost = out.Copy()

	return out
}

func ConvertCustomFields(customFields interface{}) map[string]string {
	tmp := customFields.(map[string]interface{})

	out := make(map[string]string)
	for key, val := range tmp {
		if val != nil {
			tmpVal, ok := val.(string)
			if ok {
				out[key] = tmpVal
				continue
			}

			//TODO: parse maps
		}
	}

	return out
}
