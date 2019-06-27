package virtualization

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/magiconair/properties/assert"
)

func mockNetboxVirtualMachine(addResources bool, addIPAddresses bool, addTags bool, addCustomFields bool) (out models.VirtualMachine) {
	out = models.VirtualMachine{
		Name: swag.String("VM1"),
	}

	if addResources {
		out.Vcpus = swag.Int64(1)
		out.Memory = swag.Int64(4096)
		out.Disk = swag.Int64(10240)
	}

	if addIPAddresses {
		out.PrimaryIp4 = &models.NestedIPAddress{Address: swag.String("127.0.0.1/32")}
		out.PrimaryIp6 = &models.NestedIPAddress{Address: swag.String("::1/128")}
	}

	if addTags {
		out.Tags = append(out.Tags, "Tag1")
		out.Tags = append(out.Tags, "managed")
	}

	if addCustomFields {
		customFields := make(map[string]interface{})
		customFields["hypervisor_label"] = "Hypervisor1"

		out.CustomFields = customFields
	}

	return out
}

func TestVirtualMachineConvertFromNetbox(t *testing.T) {
	c := NewClient(client.NetBox{})

	vm := mockNetboxVirtualMachine(false, false, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.Hostname, "VM1")
	assert.Equal(t, res.Resources.Cores, 0)
	assert.Equal(t, res.Resources.Memory, int64(0))
	assert.Equal(t, len(res.Resources.Disks), 0)

	assert.Equal(t, res.PrimaryIPv4, types.IPAddress{})
	assert.Equal(t, res.PrimaryIPv6, types.IPAddress{})

	assert.Equal(t, len(res.Tags), 0)
	assert.Equal(t, res.IsManaged, false)

	assert.Equal(t, res.Hypervisor, "")
}

func TestVirtualMachineConvertFromNetbox_WithResources(t *testing.T) {
	c := NewClient(client.NetBox{})

	vm := mockNetboxVirtualMachine(true, false, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.Resources.Cores, 1)
	assert.Equal(t, res.Resources.Memory, int64(4096))
	assert.Equal(t, len(res.Resources.Disks), 1)
	assert.Equal(t, res.Resources.Disks[0].Size, int64(10240*1024))
}

func TestVirtualMachineConvertFromNetbox_WithIPAddresses(t *testing.T) {
	c := NewClient(client.NetBox{})

	vm := mockNetboxVirtualMachine(false, true, false, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.PrimaryIPv4, types.IPAddress{Address: "127.0.0.1", CIDR: 32, Family: types.IPAddressFamilyIPv4})
	assert.Equal(t, res.PrimaryIPv6, types.IPAddress{Address: "::1", CIDR: 128, Family: types.IPAddressFamilyIPv6})
}

func TestVirtualMachineConvertFromNetbox_WithTags(t *testing.T) {
	c := NewClient(client.NetBox{})

	vm := mockNetboxVirtualMachine(false, false, true, false)

	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, len(res.Tags), 2)
	assert.Equal(t, res.Tags[0], "Tag1")
	assert.Equal(t, res.Tags[1], "managed")
	assert.Equal(t, res.IsManaged, true)
}

func TestVirtualMachineConvertFromNetbox_WithCustomFields(t *testing.T) {
	c := NewClient(client.NetBox{})

	vm := mockNetboxVirtualMachine(false, false, false, true)
	res, err := c.VirtualMachineConvertFromNetbox(vm, []*models.VirtualMachineInterface{})
	assert.Equal(t, err, nil)

	assert.Equal(t, res.Hypervisor, "Hypervisor1")
}
