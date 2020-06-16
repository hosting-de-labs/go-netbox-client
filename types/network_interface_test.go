package types_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"
)

func MockNetworkInterface() types.NetworkInterface {
	netIf := types.NewNetworkInterface()
	netIf.Name = "eth0"
	netIf.MACAddress = []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	netIf.Type = types.InterfaceTypeEthernetFixed1000BaseT1G

	return *netIf
}

func TestNetworkInterface_IsEqual(t *testing.T) {
	if1 := MockNetworkInterface()
	if2 := MockNetworkInterface()

	assert.Equal(t, if1, if2)
	assert.True(t, if1.IsEqual(if2))

	if2.Name = "eth1"

	assert.NotEqual(t, if1, if2)
	assert.False(t, if1.IsEqual(if2))
}

func TestNetworkInterface_IsEqualMultipleAddresses(t *testing.T) {
	if1 := MockNetworkInterface()
	if1.IPAddresses = append(if1.IPAddresses, MockIpv6Address())

	if2 := MockNetworkInterface()
	if2.IPAddresses = append(if2.IPAddresses, MockIpv6Address())

	assert.Equal(t, if1, if2)
	assert.True(t, if1.IsEqual(if2))

	if2.Name = "eth1"

	assert.NotEqual(t, if1, if2)
	assert.False(t, if1.IsEqual(if2))
}
