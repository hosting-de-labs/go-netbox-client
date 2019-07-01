package dcim_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func init() {
	mock.RunServer()
}

func mockNetboxDeviceInterface() models.DeviceInterface {
	o := models.DeviceInterface{}

	o.ID = 10
	o.Name = swag.String("eth0")
	o.MacAddress = swag.String("aa:bb:cc:dd:ee:ff")

	o.UntaggedVlan = &models.NestedVLAN{
		ID:   10,
		Vid:  swag.Int64(400),
		Name: swag.String("Public VLAN"),
	}

	o.TaggedVlans = []*models.NestedVLAN{
		{
			ID:   20,
			Vid:  swag.Int64(600),
			Name: swag.String("Private VLAN"),
		},
	}

	return o
}

func TestConvertDeviceInterface(t *testing.T) {

	netboxClient := netbox.NewNetboxAt("localhost:8000")
	dcimClient := dcim.NewClient(*netboxClient)

	netIf, err := dcimClient.InterfaceConvertFromNetbox(mockNetboxDeviceInterface())

	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, netIf.Name, "eth0")

	mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	assert.Nil(t, err)

	assert.Equal(t, netIf.MACAddress, mac)

	assert.Equal(t, len(netIf.IPAddresses), 2)

	assert.Equal(t, "123.123.123.123", netIf.IPAddresses[0].Address)
	assert.Equal(t, uint16(24), netIf.IPAddresses[0].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv4, netIf.IPAddresses[0].Family)

	assert.Equal(t, "2001:db8:a::123", netIf.IPAddresses[1].Address)
	assert.Equal(t, uint16(64), netIf.IPAddresses[1].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv6, netIf.IPAddresses[1].Family)
}
