package types_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"

	"github.com/stretchr/testify/assert"
)

func TestIPAddress_String(t *testing.T) {
	ip1 := client_types.MockIpv4Address()
	assert.Equal(t, ip1.String(), "192.168.10.1/24")

	ip2 := client_types.MockIpv6Address()
	assert.Equal(t, ip2.String(), "fc00::1/128")
}

func TestIPAddress_IsEqual(t *testing.T) {
	ip1 := client_types.MockIpv4Address()

	ip2 := ip1.Clone()
	assert.Equal(t, ip1, ip2)
	assert.True(t, ip1.IsEqual(ip2))

	ip2.Address = "192.168.10.2"
	assert.NotEqual(t, ip1, ip2)
}
