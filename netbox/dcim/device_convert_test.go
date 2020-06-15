package dcim_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox/netbox"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConvertFromNetbox(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	device := netbox_types.MockNetboxDevice(false, false)

	res, err := c.DeviceConvertFromNetbox(device)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.False(t, res.IsChanged())

	assert.Equal(t, "Host 10", res.Hostname)
	assert.Equal(t, "123-456", res.AssetTag)
	assert.Equal(t, "1234567890", res.SerialNumber)
}

func TestDeviceConvertFromNetbox_WithWrongType(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	var device interface{}
	res, err := c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDeviceConvertFromNetbox_WithIPAddresses(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	device := netbox_types.MockNetboxDevice(true, false)

	res, err := c.DeviceConvertFromNetbox(device)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.False(t, res.IsChanged())

	assert.NotNil(t, res.PrimaryIPv4)
	assert.NotNil(t, res.PrimaryIPv6)
}

func TestDeviceConvertFromNetbox_WithInvalidIPAddresses(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	device := netbox_types.MockNetboxDevice(false, false)
	device.PrimaryIp4 = &models.NestedIPAddress{Address: swag.String("123.456.789.101112")}
	res, err := c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)

	device.PrimaryIp4 = nil
	device.PrimaryIp6 = &models.NestedIPAddress{Address: swag.String("::827")}
	res, err = c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDeviceConvertFromNetbox_WithTags(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	device := netbox_types.MockNetboxDevice(false, true)

	res, err := c.DeviceConvertFromNetbox(device)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.False(t, res.IsChanged())

	assert.NotEmpty(t, res.Tags)
	assert.True(t, res.IsManaged)
	assert.Equal(t, "Tag1", res.Tags[0])
	assert.Equal(t, "managed", res.Tags[1])
}

func TestDeviceConvertToNetbox(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	device := client_types.MockDedicatedServer()
	nbDevice, nbInterfaces, err := c.DeviceConvertToNetbox(device)

	assert.Nil(t, err)
	assert.NotNil(t, nbDevice)
	assert.Nil(t, nbInterfaces)
}
