package types_test

import (
	"strconv"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"
)

func TestVirtualServer_Copy(t *testing.T) {
	vm1 := types.VirtualServer{
		Hypervisor: "hypervisor1",
		Resources: types.VirtualServerResources{
			Cores: 4,
			Disks: []types.VirtualServerDisk{
				{Size: 10},
			},
		},
	}

	vm2 := vm1.Copy()

	assert.Equal(t, vm1, vm2)
}

func TestVirtualServer_IsChanged(t *testing.T) {
	vm := types.VirtualServer{
		Hypervisor: "hypervisor1",
		Resources: types.VirtualServerResources{
			Cores: 4,
			Disks: []types.VirtualServerDisk{
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
		vm1     types.VirtualServer
		vm2     types.VirtualServer
		isEqual bool
	}{
		{
			vm1: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			isEqual: true,
		},
		{
			vm1: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: types.VirtualServer{
				Hypervisor: "hypervisor2",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 20},
					},
				},
			},
			isEqual: false,
		},
		{
			vm1: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 10},
					},
				},
			},
			vm2: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 4,
					Disks: []types.VirtualServerDisk{
						{Size: 20},
					},
				},
			},
			isEqual: false,
		},
	}

	for key, testcase := range cases {
		if testcase.isEqual {
			assert.Equal(t, testcase.vm1, testcase.vm2, "Case ID: "+strconv.Itoa(key))
			assert.True(t, testcase.vm1.IsEqual(testcase.vm2, true), "Case ID: "+strconv.Itoa(key))
		} else {
			assert.NotEqual(t, testcase.vm1, testcase.vm2, "Case ID: "+strconv.Itoa(key))
			assert.False(t, testcase.vm1.IsEqual(testcase.vm2, true), "Case ID: "+strconv.Itoa(key))
		}
	}
}
