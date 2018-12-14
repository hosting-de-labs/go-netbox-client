package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualServer_Copy(t *testing.T) {
	vm1 := VirtualServer{
		Hypervisor: "hypervisor1",
		Resources: VirtualServerResources{
			Cores: 4,
			Disks: []VirtualServerDisk{
				{Size: 10},
			},
		},
	}

	vm2 := vm1.Copy()

	assert.Equal(t, vm1, vm2)
}

func TestVirtualServer_IsChanged(t *testing.T) {
	vm := VirtualServer{
		Hypervisor: "hypervisor1",
		Resources: VirtualServerResources{
			Cores: 4,
			Disks: []VirtualServerDisk{
				{Size: 10},
			},
		},
	}

	vm.OriginalEntity = vm.Copy()
	vm.Hypervisor = "hypervisor2"

	assert.True(t, vm.IsChanged())
}

func TestVirtualServer_IsEqual(t *testing.T) {
	cases := []struct {
		vm1     VirtualServer
		vm2     VirtualServer
		isEqual bool
	}{
		{
			vm1: VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			isEqual: true,
		},
		{
			vm1: VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: VirtualServer{
				Hypervisor: "hypervisor2",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 20},
					},
				},
			},
			isEqual: false,
		},
		{
			vm1: VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: VirtualServerResources{
					Cores: 4,
					Disks: []VirtualServerDisk{
						{Size: 20},
					},
				},
			},
			isEqual: false,
		},
	}

	for key, testcase := range cases {
		if testcase.isEqual {
			assert.True(t, testcase.vm1.IsEqual(testcase.vm2, true), "Case ID: "+strconv.Itoa(key))
		} else {
			assert.False(t, testcase.vm1.IsEqual(testcase.vm2, true), "Case ID: "+strconv.Itoa(key))
		}
	}
}
