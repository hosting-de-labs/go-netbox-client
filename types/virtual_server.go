package types

import (
	"sort"
)

//VirtualServer represents a virtual server
type VirtualServer struct {
	OriginalHost *VirtualServer

	Host
	Hypervisor string

	Resources VirtualServerResources
}

//Copy creates a deep copy of a VirtualServer object
func (vm VirtualServer) Copy() *VirtualServer {
	var out VirtualServer

	out.Host = *vm.Host.Copy()

	out.Hypervisor = vm.Hypervisor
	out.Resources = VirtualServerResources{
		Cores:  vm.Resources.Cores,
		Memory: vm.Resources.Memory,
	}

	//copy disks
	out.Resources.Disks = make([]VirtualServerDisk, len(vm.Resources.Disks))
	copy(out.Resources.Disks, vm.Resources.Disks)

	return &out
}

//IsChanged compares the current object against the original object
func (vm VirtualServer) IsChanged() bool {
	return !vm.IsEqual(*vm.OriginalHost, true)
}

//IsEqual compares the current object with another VirtualServer object
func (vm VirtualServer) IsEqual(vm2 VirtualServer, deep bool) bool {
	//compare Host struct
	if !vm.Host.IsEqual(vm2.Host, deep) {
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

	sort.Sort(BySize(vm.Resources.Disks))
	sort.Sort(BySize(vm2.Resources.Disks))

	for key, disk := range vm.Resources.Disks {
		if !disk.IsEqual(vm2.Resources.Disks[key]) {
			return false
		}
	}

	return true
}

//VirtualServerResources represents the number of cores, memory and disks assigned to a VirtualServer object
type VirtualServerResources struct {
	Cores  int
	Memory int64

	Disks []VirtualServerDisk
}

//VirtualServerDisk represents a disk of a virtual server
type VirtualServerDisk struct {
	Size int64
}

type BySize []VirtualServerDisk

func (a BySize) Len() int      { return len(a) }
func (a BySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool {
	return a[i].Size < a[j].Size
}

//IsEqual compares the current object against another VirtualServerDisk object
func (d VirtualServerDisk) IsEqual(d2 VirtualServerDisk) bool {
	if d.Size != d2.Size {
		return false
	}

	return true
}
