package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockIpv4Address() (out types.IPAddress) {
	return types.IPAddress{
		Family:      types.IPAddressFamilyIPv4,
		Address:     "192.168.10.1",
		CIDR:        24,
		Status:      types.IPAddressStatusActive,
		Description: "An internal ip address",
		Tags:        []string{"internal", "netbox-sync"},
	}
}

func MockIpv6Address() (out types.IPAddress) {
	return types.IPAddress{
		Family:      types.IPAddressFamilyIPv6,
		Address:     "fc00::1",
		CIDR:        128,
		Status:      types.IPAddressStatusActive,
		Description: "An internal ip address",
		Tags:        []string{"internal", "netbox-sync"},
	}
}
