package types

import (
	"fmt"
	"net"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

//IPAddressFamily represents a type an ip-address can have
type IPAddressFamily string

const (
	//IPAddressFamilyIPv6 represents an ipv6 ip
	IPAddressFamilyIPv6 IPAddressFamily = "IPv6"

	//IPAddressFamilyIPv4 represents an ipv4 ip
	IPAddressFamilyIPv4 IPAddressFamily = "IPv4"
)

type IPAddressStatus string

const (
	IPAddressStatusActive     IPAddressStatus = "active"
	IPAddressStatusReserved   IPAddressStatus = "reserved"
	IPAddressStatusDeprecated IPAddressStatus = "deprecated"
	IPAddressStatusDHCP       IPAddressStatus = "dhcp"
)

type IPAddressRole int

const (
	_ IPAddressRole = iota
	IPAddressRoleLoopback
	IPAddressRoleSecondary
	IPAddressRoleAnycast
	IPAddressRoleVIP
	IPAddressRoleVRRP
	IPAddressRoleHSRP
	IPAddressRoleGLBP
	IPAddressRoleCARP
)

//IPAddress represents an ip address
type IPAddress struct {
	//TODO: inherit from net.IPNet, get rid of "Type"
	Family  IPAddressFamily
	Address string
	CIDR    uint16

	Status      IPAddressStatus
	Role        *IPAddressRole
	Tags        []string
	Description string
}

func (ip IPAddress) String() string {
	return fmt.Sprintf("%s/%d", ip.Address, ip.CIDR)
}

func (ip IPAddress) Clone() (out IPAddress) {
	out = IPAddress{
		Family:      ip.Family,
		Address:     ip.Address,
		CIDR:        ip.CIDR,
		Status:      ip.Status,
		Description: ip.Description,
	}

	if ip.Role != nil {
		*out.Role = *ip.Role
	}

	if len(ip.Tags) > 0 {
		out.Tags = make([]string, len(ip.Tags))
		copy(out.Tags, ip.Tags)
	}

	return out
}

//IsEqual compares the current IPAddress object against another IPAddress object
func (ip IPAddress) IsEqual(ip2 IPAddress) bool {
	return utils.CompareStruct(ip, ip2, []string{}, []string{})
}

func (ip IPAddress) IsNetwork() (bool, error) {
	ipTmp, network, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.Address, ip.CIDR))
	if err != nil {
		return false, err
	}

	return ipTmp.Equal(network.IP), nil
}

//TODO: isBroadcast
