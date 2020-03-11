package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockVirtualServer() types.VirtualServer {
	vm := types.NewVirtualServer()
	vm.Hypervisor = "hypervisor1"
	vm.Resources.Cores = 4
	vm.Resources.Disks = []types.VirtualServerDisk{{Size: 10}}

	vm.Meta.OriginalEntity = vm.Copy()
	return *vm
}
