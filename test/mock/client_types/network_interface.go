package client_types

import (
	"github.com/hosting-de-labs/go-netbox-client/types"
)

func MockNetworkInterface() types.NetworkInterface {
	netIf := types.NewNetworkInterface()
	netIf.Name = "eth0"
	netIf.MACAddress = []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	netIf.Type = types.InterfaceTypeEthernetFixed1000BaseT1G

	return *netIf
}
