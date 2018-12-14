package types

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

//IsEqual compares the current object against another VirtualServerDisk object
func (d VirtualServerDisk) IsEqual(d2 VirtualServerDisk) bool {
	if d.Size != d2.Size {
		return false
	}

	return true
}
