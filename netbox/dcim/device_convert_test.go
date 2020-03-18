package dcim_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConvertFromNetbox(t *testing.T) {
	c := dcim.NewClient(client.NetBox{})

	device := netbox_types.MockNetboxDevice()

	res, err := c.DeviceConvertFromNetbox(device)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.False(t, res.IsChanged())

	assert.Equal(t, "Host 10", res.Hostname)
	assert.Equal(t, "123-456", res.AssetTag)
	assert.Equal(t, "1234567890", res.SerialNumber)
}
