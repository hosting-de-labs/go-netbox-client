package types

import (
	"sort"
)

//VirtualServer represents a virtual server
type VirtualServer struct {
	Host

	Hypervisor string
	Resources  VirtualServerResources
}

func NewVirtualServer() *VirtualServer {
	return &VirtualServer{
		Host: Host{
			CommonEntity: CommonEntity{
				Metadata: &Metadata{},
			},
		},
		Hypervisor: "",
		Resources:  VirtualServerResources{},
	}
}

//Copy creates a deep copy of a VirtualServer object
func (vm VirtualServer) Copy() (out VirtualServer) {
	out.Host = vm.Host.Copy()

	out.Hypervisor = vm.Hypervisor
	out.Resources = VirtualServerResources{
		Cores:  vm.Resources.Cores,
		Memory: vm.Resources.Memory,
	}

	//copy disks
	out.Resources.Disks = make([]VirtualServerDisk, len(vm.Resources.Disks))
	copy(out.Resources.Disks, vm.Resources.Disks)

	return out
}

//IsChanged compares the current object against the original object
func (vm VirtualServer) IsChanged() bool {
	return !vm.IsEqual(vm.Metadata.NetboxEntity.(VirtualServer), true)
}

//IsEqual compares the current object with another VirtualServer object
func (vm VirtualServer) IsEqual(vm2 VirtualServer, deep bool) bool {
	vm.Metadata = nil
	vm2.Metadata = nil

	//compare Host struct
	if !vm.Host.IsEqual(vm2.Host, deep) {
		return false
	}

	if vm.Hypervisor != vm2.Hypervisor {
		return false
	}

	//Resources
	if vm.Resources.Cores != vm2.Resources.Cores {
		return false
	}

	if vm.Resources.Memory != vm2.Resources.Memory {
		return false
	}

	if len(vm.Resources.Disks) != len(vm2.Resources.Disks) {
		return false
	}

	sort.Slice(vm.Resources.Disks, func(i int, j int) bool {
		return vm.Resources.Disks[i].Size < vm.Resources.Disks[j].Size
	})

	for key, disk := range vm.Resources.Disks {
		if !disk.IsEqual(vm2.Resources.Disks[key]) {
			return false
		}
	}

	return true
}
