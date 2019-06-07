package types

import (
	"fmt"
	"net"

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

type IPAddressStatus int

const (
	_ IPAddressStatus = iota
	IPAddressStatusActive
	IPAddressStatusReserved
	IPAddressStatusDeprecated
	IPAddressStatusDHCP
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
	Type    IPAddressType
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
		Type:        ip.Type,
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
