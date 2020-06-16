package dcim_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"

	"github.com/hosting-de-labs/go-netbox/netbox"
)

func TestInterfaceConvertFromNetboxDeviceInterface(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	netIf, err := c.InterfaceFind(2, "eth0")
	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, netIf.Name, "eth0")
	assert.Equal(t, netIf.Type, types.InterfaceTypeEthernetFixed1000BaseT1G)
	assert.Equal(t, net.HardwareAddr{0xab, 0xbc, 0xcd, 0xde, 0xef, 0xfa}, netIf.MACAddress)

	assert.Len(t, netIf.IPAddresses, 2)

	assert.Equal(t, "10.123.123.123", netIf.IPAddresses[0].Address)
	assert.Equal(t, uint16(24), netIf.IPAddresses[0].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv4, netIf.IPAddresses[0].Family)

	assert.Equal(t, "2001:db8:fefe::123", netIf.IPAddresses[1].Address)
	assert.Equal(t, uint16(64), netIf.IPAddresses[1].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv6, netIf.IPAddresses[1].Family)
}

func TestInterfaceConvertToNetboxDeviceInterface(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	c := dcim.NewClient(*netboxClient)

	intf, err := c.InterfaceFind(1, "eth0")
	assert.Nil(t, err)
	assert.NotNil(t, intf)

	nbIntf, err := c.InterfaceConvertToNetbox(1, *intf)
	assert.Nil(t, err)
	assert.NotNil(t, nbIntf)

	assert.Equal(t, "eth0", nbIntf.Name)
	assert.Equal(t, "aa:bb:cc:dd:ee:ff", *nbIntf.MacAddress)
	assert.False(t, nbIntf.MgmtOnly)
	assert.Equal(t, types.InterfaceTypeEthernetFixed1000BaseT1G, types.InterfaceType(nbIntf.Type))
}
