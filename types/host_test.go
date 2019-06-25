package types

import (
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockHost() Host {
	return Host{
		ID:       1,
		Hostname: "host1",
		PrimaryIPv4: IPAddress{
			Address: "192.168.1.1",
			CIDR:    24,
			Type:    IPAddressTypeIPv4,
		},
		PrimaryIPv6: IPAddress{
			Address: "::1",
			CIDR:    64,
			Type:    IPAddressTypeIPv6,
		},
	}
}

func TestHost_HasTag(t *testing.T) {
	host := mockHost()
	assert.False(t, host.HasTag("tag1"))

	host.AddTag("tag2")
	assert.True(t, host.HasTag("tag2"))
}

func TestHost_AddTag(t *testing.T) {
	host := mockHost()

	host.AddTag("tag1")
	assert.True(t, host.HasTag("tag1"))
	assert.False(t, host.HasTag("tag2"))

	host.AddTag("tag2")
	assert.True(t, host.HasTag("tag1"))
	assert.True(t, host.HasTag("tag2"))

	host.AddTag("tag1")
	assert.True(t, host.HasTag("tag1"))
	assert.True(t, host.HasTag("tag2"))
}

func TestHost_Copy(t *testing.T) {
	host1 := mockHost()
	host2 := host1.Copy()
	assert.Equal(t, host1, host2)

	host2.Hostname = "host2"
	assert.NotEqual(t, host1, host2)

	host3 := mockHost()

	mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	assert.Nil(t, err)

	host3.NetworkInterfaces = append(host3.NetworkInterfaces, HostNetworkInterface{
		Name:       "vlan.1",
		MACAddress: mac,
		IPAddresses: []IPAddress{
			{
				Address: "192.168.10.1",
				CIDR:    24,
				Type:    IPAddressTypeIPv4,
			},
		},
	})

	host4 := host3.Copy()
	assert.Equal(t, host3, host4)

	host4.NetworkInterfaces[0].Name = "vlan.2"
	assert.NotEqual(t, host3, host4)

	host5 := mockHost()
	host5.AddTag("tag1")
	host6 := host5.Copy()
	assert.Equal(t, host5, host6)

	host6.AddTag("tag2")
	assert.NotEqual(t, host5, host6)
}

func TestHost_IsChanged(t *testing.T) {
	host := mockHost()
	host.CommonEntity.OriginalEntity = host.Copy()
	assert.False(t, host.IsChanged())

	host.Hostname = "host2"
	assert.True(t, host.IsChanged())
}

func TestHost_IsEqual(t *testing.T) {
	cases := []struct {
		host1   Host
		host2   Host
		isEqual bool
	}{
		{
			host1:   Host{},
			host2:   Host{},
			isEqual: true,
		},
		{
			host1: Host{
				ID:        10,
				Hostname:  "Server",
				IsManaged: true,
				PrimaryIPv4: IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Type:    IPAddressTypeIPv4,
				},
				PrimaryIPv6: IPAddress{
					Address: "::1",
					CIDR:    128,
					Type:    IPAddressTypeIPv6,
				},
				Comments: []string{
					"Comment1",
					"Comment2",
				},
				Tags: []string{
					"Tag1",
					"Tag2",
				},
			},
			host2: Host{
				ID:        10,
				Hostname:  "Server",
				IsManaged: true,
				PrimaryIPv4: IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Type:    IPAddressTypeIPv4,
				},
				PrimaryIPv6: IPAddress{
					Address: "::1",
					CIDR:    128,
					Type:    IPAddressTypeIPv6,
				},
				Comments: []string{
					"Comment1",
					"Comment2",
				},
				Tags: []string{
					"Tag1",
					"Tag2",
				},
			},
			isEqual: true,
		},
		{
			host1: Host{
				ID: 10,
			},
			host2: Host{
				ID: 20,
			},
			isEqual: false,
		},
		{
			host1: Host{
				Hostname: "Server1",
			},
			host2: Host{
				Hostname: "Server2",
			},
			isEqual: false,
		},
		{
			host1: Host{
				IsManaged: true,
			},
			host2: Host{
				IsManaged: false,
			},
			isEqual: false,
		},
		{
			host1: Host{
				PrimaryIPv4: IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Type:    IPAddressTypeIPv4,
				},
			},
			host2: Host{
				PrimaryIPv4: IPAddress{
					Address: "10.10.10.2",
					CIDR:    24,
					Type:    IPAddressTypeIPv4,
				},
			},
			isEqual: false,
		},
		{
			host1: Host{
				PrimaryIPv6: IPAddress{
					Address: "::1",
					CIDR:    128,
					Type:    IPAddressTypeIPv6,
				},
			},
			host2: Host{
				PrimaryIPv6: IPAddress{
					Address: "::2",
					CIDR:    128,
					Type:    IPAddressTypeIPv6,
				},
			},
			isEqual: false,
		},
		{
			host1: Host{
				Tags: []string{"Tag1"},
			},
			host2: Host{
				Tags: []string{"Tag2"},
			},
			isEqual: false,
		},
		{
			host1: Host{
				Comments: []string{"Comment1"},
			},
			host2: Host{
				Comments: []string{"Comment2"},
			},
			isEqual: false,
		},
		{
			host1: Host{
				NetworkInterfaces: []HostNetworkInterface{
					{Name: "eth0"},
					{Name: "eth1"},
				},
			},
			host2: Host{
				NetworkInterfaces: []HostNetworkInterface{
					{Name: "eth1"},
					{Name: "eth0"},
				},
			},
			isEqual: true,
		},
	}

	for key, testcase := range cases {
		if testcase.isEqual {
			assert.True(t, testcase.host1.IsEqual(testcase.host2, true), "Case ID: "+strconv.Itoa(key))
		} else {
			assert.False(t, testcase.host1.IsEqual(testcase.host2, true), "Case ID: "+strconv.Itoa(key))
		}
	}
}
