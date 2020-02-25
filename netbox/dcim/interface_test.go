package dcim_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock"

	"github.com/hosting-de-labs/go-netbox-client/netbox/dcim"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/stretchr/testify/assert"
)

func init() {
	mock.RunServer()
}

func TestInterfaceGet(t *testing.T) {
	netboxClient := netbox.NewNetboxAt("localhost:8000")

	dcimClient := dcim.NewClient(*netboxClient)
	netIf, err := dcimClient.InterfaceGet(10)

	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, "eth0", netIf.Name)
	assert.Equal(t, types.InterfaceTypeEthernetFixed1000BaseT1G, *netIf.Type)

	assert.Equal(t, net.HardwareAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}), netIf.MACAddress)

	assert.NotNil(t, netIf.UntaggedVlan)
	assert.Equal(t, uint16(400), netIf.UntaggedVlan.ID)
	assert.Equal(t, "vlan-400", netIf.UntaggedVlan.Name)

	assert.Empty(t, netIf.TaggedVlans)
	assert.Empty(t, netIf.Tags)
}
