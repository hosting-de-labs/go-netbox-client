package types

//VirtualServer represents a virtual server
type VirtualServer struct {
	Host

	Hypervisor string
	Resources  VirtualServerResources
}

//NewVirtualServer returns a new instance of VirtualServer
func NewVirtualServer() *VirtualServer {
	return &VirtualServer{
		Host: Host{
			CommonEntity: CommonEntity{
				Meta: &Metadata{},
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
	if vm.Meta == nil || vm.Meta.OriginalEntity == nil {
		return true
	}

	return !vm.IsEqual(*vm.Meta.OriginalEntity.(*VirtualServer), true)
}

//IsEqual compares the current object with another VirtualServer object
func (vm VirtualServer) IsEqual(vm2 VirtualServer, deep bool) bool {
	vm.Meta = nil
	vm2.Meta = nil

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

	for key, disk := range vm.Resources.Disks {
		if !disk.IsEqual(vm2.Resources.Disks[key]) {
			return false
		}
	}

	return true
}
