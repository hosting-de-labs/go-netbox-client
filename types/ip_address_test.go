package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockIpv4Address() (out IPAddress) {
	return IPAddress{
		Type:        IPAddressTypeIPv4,
		Address:     "192.168.10.1",
		CIDR:        24,
		Status:      IPAddressStatusActive,
		Description: "An internal ip address",
		Tags:        []string{"internal", "netbox-sync"},
	}
}

func mockIpv6Address() (out IPAddress) {
	return IPAddress{
		Type:        IPAddressTypeIPv6,
		Address:     "fc00::1",
		CIDR:        128,
		Status:      IPAddressStatusActive,
		Description: "An internal ip address",
		Tags:        []string{"internal", "netbox-sync"},
	}
}

func TestIPAddress_String(t *testing.T) {
	ip1 := mockIpv4Address()
	assert.Equal(t, ip1.String(), "192.168.10.1/24")

	ip2 := mockIpv6Address()
	assert.Equal(t, ip2.String(), "fc00::1/128")
}

func TestIPAddress_IsEqual(t *testing.T) {
	ip1 := mockIpv4Address()

	ip2 := ip1.Clone()
	assert.Equal(t, ip1, ip2)
	assert.True(t, ip1.IsEqual(ip2))

	ip2.Address = "192.168.10.2"
	assert.NotEqual(t, ip1, ip2)
}