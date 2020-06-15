package virtualization_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/netbox/virtualization"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/stretchr/testify/assert"
)

func TestConvertVirtualMachineInterface(t *testing.T) {
	netboxClient := netbox.NewNetboxWithAPIKey("localhost:8080", "0123456789abcdef0123456789abcdef01234567")
	virtualizationClient := virtualization.NewClient(*netboxClient)

	netboxIf := netbox_types.MockNetboxVirtualMachineInterface()
	netIf, err := virtualizationClient.InterfaceConvertFromNetbox(netboxIf)

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
