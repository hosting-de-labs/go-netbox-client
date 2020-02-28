package types_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/client_types"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"
)

func TestHost_HasTag(t *testing.T) {
	host := client_types.MockHost()
	assert.False(t, host.HasTag("tag1"))

	host.AddTag("tag2")
	assert.True(t, host.HasTag("tag2"))
}

func TestHost_AddTag(t *testing.T) {
	host := client_types.MockHost()

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
	host1 := client_types.MockHost()
	host2 := host1.Copy()
	assert.Equal(t, host1, host2)

	host2.Hostname = "host2"
	assert.NotEqual(t, host1, host2)

	host3 := client_types.MockHost()

	mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	assert.Nil(t, err)

	host3.NetworkInterfaces = append(host3.NetworkInterfaces, types.NetworkInterface{
		Name:       "vlan.1",
		MACAddress: mac,
		IPAddresses: []types.IPAddress{
			{
				Address: "192.168.10.1",
				CIDR:    24,
				Family:  types.IPAddressFamilyIPv4,
			},
		},
	})

	host4 := host3.Copy()
	assert.Equal(t, host3, host4)

	host4.NetworkInterfaces[0].Name = "vlan.2"
	assert.NotEqual(t, host3, host4)

	host5 := client_types.MockHost()
	host5.AddTag("tag1")
	host6 := host5.Copy()
	assert.Equal(t, host5, host6)

	host5.AddTag("tag2")
	assert.NotEqual(t, host5, host6)
}

func TestHost_IsChanged(t *testing.T) {
	host := client_types.MockHost()

	assert.False(t, host.IsChanged())

	host.Hostname = "host2"
	assert.True(t, host.IsChanged())
}

func TestHost_IsEqual(t *testing.T) {
	cases := []struct {
		host1   types.Host
		host2   types.Host
		isEqual bool
	}{
		{
			host1:   types.Host{},
			host2:   types.Host{},
			isEqual: true,
		},
		{
			host1: types.Host{
				Hostname:  "Server",
				IsManaged: true,
				PrimaryIPv4: types.IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Family:  types.IPAddressFamilyIPv4,
				},
				PrimaryIPv6: types.IPAddress{
					Address: "::1",
					CIDR:    128,
					Family:  types.IPAddressFamilyIPv6,
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
			host2: types.Host{
				Hostname:  "Server",
				IsManaged: true,
				PrimaryIPv4: types.IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Family:  types.IPAddressFamilyIPv4,
				},
				PrimaryIPv6: types.IPAddress{
					Address: "::1",
					CIDR:    128,
					Family:  types.IPAddressFamilyIPv6,
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
			host1: types.Host{
				Hostname: "Server1",
			},
			host2: types.Host{
				Hostname: "Server2",
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				IsManaged: true,
			},
			host2: types.Host{
				IsManaged: false,
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				PrimaryIPv4: types.IPAddress{
					Address: "10.10.10.1",
					CIDR:    24,
					Family:  types.IPAddressFamilyIPv4,
				},
			},
			host2: types.Host{
				PrimaryIPv4: types.IPAddress{
					Address: "10.10.10.2",
					CIDR:    24,
					Family:  types.IPAddressFamilyIPv4,
				},
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				PrimaryIPv6: types.IPAddress{
					Address: "::1",
					CIDR:    128,
					Family:  types.IPAddressFamilyIPv6,
				},
			},
			host2: types.Host{
				PrimaryIPv6: types.IPAddress{
					Address: "::2",
					CIDR:    128,
					Family:  types.IPAddressFamilyIPv6,
				},
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				Tags: []string{"Tag1"},
			},
			host2: types.Host{
				Tags: []string{"Tag2"},
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				Comments: []string{"Comment1"},
			},
			host2: types.Host{
				Comments: []string{"Comment2"},
			},
			isEqual: false,
		},
		{
			host1: types.Host{
				NetworkInterfaces: []types.NetworkInterface{
					{Name: "eth0"},
					{Name: "eth1"},
				},
			},
			host2: types.Host{
				NetworkInterfaces: []types.NetworkInterface{
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
