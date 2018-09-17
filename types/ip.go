package types

import "fmt"

//IPAddressType represents a type an ip-address can have
type IPAddressType string

const (
	//IPAddressTypeIPv6 represents an ipv6 ip
	IPAddressTypeIPv6 IPAddressType = "IPv6"

	//IPAddressTypeIPv4 represents an ipv4 ip
	IPAddressTypeIPv4 IPAddressType = "IPv4"
)

//IPAddress represents an ip address
type IPAddress struct {
	Type    IPAddressType
	Address string
	CIDR    uint16
}

func (i IPAddress) String() string {
	return fmt.Sprintf("%s/%d", i.Address, i.CIDR)
}

type ByAddress []IPAddress

func (a ByAddress) Len() int      { return len(a) }
func (a ByAddress) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAddress) Less(i, j int) bool {
	return a[i].Address < a[j].Address
}

//IsEqual compares the current IPAddress object against another IPAddress object
func (n IPAddress) IsEqual(n2 IPAddress) bool {
	if n.Type != n2.Type {
		return false
	}

	if n.Address != n2.Address {
		return false
	}

	if n.CIDR != n2.CIDR {
		return false
	}

	return true
}
