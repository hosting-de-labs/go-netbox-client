package types_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"

	"github.com/stretchr/testify/assert"
)

func TestNetworkInterface_IsEqual(t *testing.T) {
	if1 := client_types.MockNetworkInterface()
	if2 := client_types.MockNetworkInterface()

	assert.Equal(t, if1, if2)
	assert.True(t, if1.IsEqual(if2))

	if2.Name = "eth1"

	assert.NotEqual(t, if1, if2)
	assert.False(t, if1.IsEqual(if2))
}

func TestNetworkInterface_IsEqualMultipleAddresses(t *testing.T) {
	if1 := client_types.MockNetworkInterface()
	if1.IPAddresses = append(if1.IPAddresses, client_types.MockIpv6Address())

	if2 := client_types.MockNetworkInterface()
	if2.IPAddresses = append(if2.IPAddresses, client_types.MockIpv6Address())

	assert.Equal(t, if1, if2)
	assert.True(t, if1.IsEqual(if2))

	if2.Name = "eth1"

	assert.NotEqual(t, if1, if2)
	assert.False(t, if1.IsEqual(if2))
}
