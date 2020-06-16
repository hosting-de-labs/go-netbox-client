package dcim_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/stretchr/testify/assert"
)

func TestInterfaceGet(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")

	c := dcim.NewClient(*netboxClient)
	netIf, err := c.InterfaceGet(1)
	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, "eth0", netIf.Name)
	assert.Equal(t, types.InterfaceTypeEthernetFixed1000BaseT1G, netIf.Type)
	assert.Equal(t, net.HardwareAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}), netIf.MACAddress)

	assert.Nil(t, netIf.UntaggedVlan)
	assert.Empty(t, netIf.TaggedVlans)
	assert.Empty(t, netIf.Tags)
}

func TestInterfaceGet_WithVlan(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")

	c := dcim.NewClient(*netboxClient)
	netIf, err := c.InterfaceGet(3)
	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, "eth0", netIf.Name)
	assert.Equal(t, types.InterfaceTypeEthernetFixed1000BaseT1G, netIf.Type)
	assert.Equal(t, net.HardwareAddr([]byte{0xbb, 0xcc, 0xdd, 0xee, 0xff, 0xaa}), netIf.MACAddress)

	assert.NotNil(t, netIf.UntaggedVlan)
	assert.Equal(t, uint16(5), netIf.UntaggedVlan.ID)
	assert.Equal(t, "vlan1", netIf.UntaggedVlan.Name)

	assert.Empty(t, netIf.TaggedVlans)
	assert.Empty(t, netIf.Tags)
}
