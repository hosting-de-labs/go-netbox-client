package types

import "sort"

//HostNetworkInterface represents a network interface assigned to a host
type HostNetworkInterface struct {
	Name       string
	VlanTag    int64
	MACAddress string

	IPAddresses []IPAddress
}

//IsEqual compares the current HostNetworkInterface object against another HostNetworkInterface
func (netIf HostNetworkInterface) IsEqual(netIf2 HostNetworkInterface) bool {
	if netIf.Name != netIf2.Name {
		return false
	}

	if netIf.VlanTag != netIf2.VlanTag {
		return false
	}

	if netIf.MACAddress != netIf2.MACAddress {
		return false
	}

	if len(netIf.IPAddresses) != len(netIf2.IPAddresses) {
		return false
	}

	sort.Sort(ByAddress(netIf.IPAddresses))
	sort.Sort(ByAddress(netIf2.IPAddresses))

	for i := 0; i < len(netIf.IPAddresses); i++ {
		if !netIf.IPAddresses[i].IsEqual(netIf2.IPAddresses[i]) {

			return false
		}
	}

	return true
}
