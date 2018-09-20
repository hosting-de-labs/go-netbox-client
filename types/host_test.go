package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostIsEqual(t *testing.T) {
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
