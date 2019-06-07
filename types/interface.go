package types

import (
	"sort"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

//HostNetworkInterface represents a network interface assigned to a host
type HostNetworkInterface struct {
	Name         string
	UntaggedVlan *VLAN
	TaggedVlans  []VLAN
	MACAddress   string

	IPAddresses []IPAddress
}

//IsEqual compares the current HostNetworkInterface object against another HostNetworkInterface
func (netIf HostNetworkInterface) IsEqual(netIf2 HostNetworkInterface) bool {
	if !utils.CompareStruct(netIf, netIf2, []string{}, []string{"IPAddresses"}) {
		return false
	}

	//sort ip addresses
	sort.Slice(netIf.IPAddresses, func(i, j int) bool { return netIf.IPAddresses[i].Address < netIf.IPAddresses[j].Address })
	sort.Slice(netIf2.IPAddresses, func(i, j int) bool { return netIf2.IPAddresses[i].Address < netIf2.IPAddresses[j].Address })

	//compare each address
	for i := 0; i < len(netIf.IPAddresses); i++ {
		if !netIf.IPAddresses[i].IsEqual(netIf2.IPAddresses[i]) {
			return false
		}
	}

	return true
}
