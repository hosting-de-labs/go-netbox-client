package client_types

import (
	"net"

	"github.com/hosting-de-labs/go-netbox-client/types"
)

func MockNetworkInterface() types.NetworkInterface {
	return types.NetworkInterface{
		Name:         "eth0",
		MACAddress:   net.HardwareAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		IsManagement: false,
		Type:         types.InterfaceTypeEthernetFixed1000BaseT1G,
	}
}
