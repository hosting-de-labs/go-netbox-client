package types

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

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

//IsEqual compares the current IPAddress object against another IPAddress object
func (i IPAddress) IsEqual(i2 IPAddress) bool {
	return utils.CompareStruct(i, i2, []string{}, []string{})
}
