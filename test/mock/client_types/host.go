package client_types

import "github.com/hosting-de-labs/go-netbox-client/types"

func MockHost() types.Host {
	h := types.NewHost()
	h.Hostname = "host1"

	h.PrimaryIPv4 = &types.IPAddress{
		Address: "192.168.1.1",
		CIDR:    24,
		Family:  types.IPAddressFamilyIPv4,
	}
	h.PrimaryIPv6 = &types.IPAddress{
		Address: "::1",
		CIDR:    64,
		Family:  types.IPAddressFamilyIPv6,
	}

	h.Meta.OriginalEntity = h.Copy()

	return *h
}
