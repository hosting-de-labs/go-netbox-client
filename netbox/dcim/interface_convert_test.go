package dcim_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_api"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"

	"github.com/hosting-de-labs/go-netbox/netbox"
)

func init() {
	netbox_api.RunServer()
}

func TestInterfaceConvertFromNetboxDeviceInterface(t *testing.T) {
	netboxClient := netbox.NewNetboxAt("localhost:8000")
	dcimClient := dcim.NewClient(*netboxClient)

	netIf, err := dcimClient.InterfaceConvertFromNetbox(netbox_types.MockNetboxDeviceInterface())

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

func TestInterfaceConvertToNetboxDeviceInterface(t *testing.T) {
	netboxClient := netbox.NewNetboxAt("localhost:8000")
	dcimClient := dcim.NewClient(*netboxClient)

	intf, err := dcimClient.InterfaceConvertToNetbox(10, client_types.MockNetworkInterface())
	assert.NotNil(t, intf)
	assert.Nil(t, err)

	assert.Equal(t, "eth0", intf.Name)
	assert.Equal(t, "aa:bb:cc:dd:ee:ff", *intf.MacAddress)
	assert.False(t, intf.MgmtOnly)
	assert.Equal(t, types.InterfaceTypeEthernetFixed1000BaseT1G, types.InterfaceType(intf.Type))
}
