package types_test

import (
	"strconv"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/stretchr/testify/assert"
)

func TestVirtualServer_Copy(t *testing.T) {
	vm1 := client_types.MockVirtualServer()
	vm1.Hypervisor = "hypervisor1"
	vm1.Resources.Cores = 4
	vm1.Resources.Disks = []types.VirtualServerDisk{{Size: 10}}

	vm2 := vm1.Copy()

	assert.Equal(t, vm1, vm2)
}

func TestVirtualServer_IsChanged(t *testing.T) {
	vm := client_types.MockVirtualServer()
	vm.Hypervisor = "hypervisor2"

	assert.True(t, vm.IsChanged())
}

func TestVirtualServer_IsChangedWithEmptyMetadata(t *testing.T) {
	vm := client_types.MockVirtualServer()

	vm.Hypervisor = "hypervisor2"
	vm.Meta = nil

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
				},
			},
			vm2: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Cores: 2,
				},
			},
			isEqual: false,
		},
		{
			vm1: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Memory: 1024,
				},
			},
			vm2: types.VirtualServer{
				Hypervisor: "hypervisor1",
				Resources: types.VirtualServerResources{
					Memory: 2048,
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
