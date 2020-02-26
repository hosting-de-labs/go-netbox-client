package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockHost() types.Host {
	return types.Host{
		Hostname: "host1",
		PrimaryIPv4: types.IPAddress{
			Address: "192.168.1.1",
			CIDR:    24,
			Family:  types.IPAddressFamilyIPv4,
		},
		PrimaryIPv6: types.IPAddress{
			Address: "::1",
			CIDR:    64,
			Family:  types.IPAddressFamilyIPv6,
		},
	}
}
