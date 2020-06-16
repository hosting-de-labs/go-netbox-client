package virtualization_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox/netbox"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"

	"github.com/hosting-de-labs/go-netbox-client/netbox/virtualization"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/stretchr/testify/assert"
)

func TestVirtualMachineConvertFromNetbox(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := virtualization.NewClient(*netboxClient)

	vm, err := c.VirtualMachineFind("virtual machine 1")
	assert.Nil(t, err)
	assert.NotNil(t, vm)

	assert.Equal(t, "virtual machine 1", vm.Hostname)
	assert.Equal(t, 8, vm.Resources.Cores)
	assert.Equal(t, int64(4096), vm.Resources.Memory)
	assert.NotEmpty(t, vm.Resources.Disks)
	assert.Len(t, vm.Resources.Disks, 1)
	assert.Equal(t, int64(204800), vm.Resources.Disks[0].Size)

	assert.Nil(t, vm.PrimaryIPv4)
	assert.Nil(t, vm.PrimaryIPv6)

	assert.Empty(t, vm.Tags)
	assert.False(t, vm.IsManaged)

	assert.Empty(t, vm.Hypervisor)
}

func TestVirtualMachineConvertFromNetbox_WithUnknownType(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	type UnknownVMType interface{}
	var vm UnknownVMType

	_, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Error(t, err)
}

func TestVirtualMachineConvertFromNetbox_WithResources(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	vm := netbox_types.MockNetboxVirtualMachine(true, false, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.Resources.Cores, 1)
	assert.Equal(t, res.Resources.Memory, int64(4096))
	assert.Equal(t, len(res.Resources.Disks), 1)
	assert.Equal(t, res.Resources.Disks[0].Size, int64(10240*1024))
}

func TestVirtualMachineConvertFromNetbox_WithIPAddresses(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	vm := netbox_types.MockNetboxVirtualMachine(false, true, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.PrimaryIPv4, &types.IPAddress{Address: "127.0.0.1", CIDR: 32, Family: types.IPAddressFamilyIPv4})
	assert.Equal(t, res.PrimaryIPv6, &types.IPAddress{Address: "::1", CIDR: 128, Family: types.IPAddressFamilyIPv6})
}

func TestVirtualMachineConvertFromNetbox_WithTags(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	vm := netbox_types.MockNetboxVirtualMachine(false, false, true, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.NotNil(t, res)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(res.Tags))
	assert.Equal(t, "Tag1", res.Tags[0])
	assert.Equal(t, "managed", res.Tags[1])
	assert.Equal(t, true, res.IsManaged)
}

func TestVirtualMachineConvertFromNetbox_WithCustomFields(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	vm := netbox_types.MockNetboxVirtualMachine(false, false, false, true)
	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.NotNil(t, res)
	assert.Nil(t, err)

	assert.Equal(t, "Hypervisor1", res.Hypervisor)
}
