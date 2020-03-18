package virtualization_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"

	"github.com/hosting-de-labs/go-netbox-client/netbox/virtualization"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/stretchr/testify/assert"
)

func TestVirtualMachineConvertFromNetbox(t *testing.T) {
	c := virtualization.NewClient(client.NetBox{})

	vm := netbox_types.MockNetboxVirtualMachine(false, false, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.Hostname, "VM1")
	assert.Equal(t, res.Resources.Cores, 0)
	assert.Equal(t, res.Resources.Memory, int64(0))
	assert.Equal(t, len(res.Resources.Disks), 0)

	assert.Equal(t, &types.IPAddress{}, res.PrimaryIPv4)
	assert.Equal(t, &types.IPAddress{}, res.PrimaryIPv6)

	assert.Equal(t, len(res.Tags), 0)
	assert.Equal(t, res.IsManaged, false)

	assert.Equal(t, res.Hypervisor, "")
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
