package dcim_test

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConvertFromNetbox(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	device, err := c.DeviceFind("host1")
	assert.Nil(t, err)
	assert.NotNil(t, device)

	assert.False(t, device.IsChanged())

	assert.Equal(t, "host1", device.Hostname)
	assert.Equal(t, "123456", device.AssetTag)
	assert.Equal(t, "1234567890", device.SerialNumber)
}

func TestDeviceConvertFromNetbox_WithWrongType(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	var device interface{}
	res, err := c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDeviceConvertFromNetbox_WithIPAddresses(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	tmpDevice, _ := c.DeviceFind("host2")
	tmpNbDevice, _ := tmpDevice.GetMetaNetboxEntity()

	res, err := c.DeviceConvertFromNetbox(tmpNbDevice)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.False(t, res.IsChanged())

	assert.NotNil(t, res.PrimaryIPv4)
	assert.NotNil(t, res.PrimaryIPv6)
}

func TestDeviceConvertFromNetbox_WithInvalidIPAddresses(t *testing.T) {
	device := (func() (device models.Device) {
		return models.Device{
			AssetTag:    swag.String("123-456"),
			Created:     strfmt.Date(time.Now()),
			DisplayName: "",
			ID:          10,
			LastUpdated: strfmt.DateTime(time.Now()),
			Name:        swag.String("host1"),
			Serial:      "1234567890",
			Status: &models.DeviceStatus{
				Label: swag.String("Active"),
				Value: swag.String("active"),
			},
			PrimaryIp4: &models.NestedIPAddress{Address: swag.String("123.456.789.101112")},
			PrimaryIp6: &models.NestedIPAddress{Address: swag.String("::827")},
		}
	})()

	c := dcim.NewClient(client.NetBox{})
	res, err := c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)

	device.PrimaryIp4 = nil
	res, err = c.DeviceConvertFromNetbox(device)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDeviceConvertToNetbox(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	device, err := c.DeviceFind("host1")
	assert.Nil(t, err)
	assert.NotNil(t, device)

	nbDevice, nbInterfaces, err := c.DeviceConvertToNetbox(*device)
	assert.Nil(t, err)
	assert.NotNil(t, nbDevice)

	assert.Nil(t, nbInterfaces)
}
