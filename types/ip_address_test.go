package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPAddress_IsEqual(t *testing.T) {
	ip1 := IPAddress{
		Type:    IPAddressTypeIPv4,
		Address: "192.168.10.1",
		CIDR:    24,
	}

	ip2 := ip1.Clone()
	assert.Equal(t, ip1, ip2)
	assert.True(t, ip1.IsEqual(ip2))

	ip2.Address = "192.168.10.2"
	assert.NotEqual(t, ip1, ip2)
}
