package types

import (
	"net"
	"sort"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

//InterfaceType represents the type of an interface
type InterfaceType string

const (
	InterfaceTypeVirtualInterfacesVirtual         InterfaceType = "virtual"
	InterfaceTypeVirtualInterfacesLinkAggregation InterfaceType = "lag"
	InterfaceTypeEthernetFixed100BaseTx100M       InterfaceType = "100base-tx"
	InterfaceTypeEthernetFixed1000BaseT1G         InterfaceType = "1000base-t"
	InterfaceTypeEthernetFixedGbic1G              InterfaceType = "2.5gbase-t"
	InterfaceTypeEthernetFixedSfp1G               InterfaceType = "5gbase-t"
	InterfaceTypeEthernetFixed10GbaseT10G         InterfaceType = "10gbase-t"
	InterfaceTypeEthernetFixed10GbaseCx410G       InterfaceType = "10gbase-cx4"
	InterfaceTypeEthernetModularSfpPlus10G        InterfaceType = "10gbase-x-sfpp"
	InterfaceTypeEthernetModularXfp10G            InterfaceType = "10gbase-x-xfp"
	InterfaceTypeEthernetModularXenpak10G         InterfaceType = "10gbase-x-xenpak"
	InterfaceTypeEthernetModularX210G             InterfaceType = "10gbase-x-x2"
	InterfaceTypeEthernetModularSfp2825G          InterfaceType = "25gbase-x-sfp28"
	InterfaceTypeEthernetModularQsfpPlus40G       InterfaceType = "40gbase-x-qsfpp"
	InterfaceTypeEthernetModularQsfp2850G         InterfaceType = "50gbase-x-sfp28"
	InterfaceTypeEthernetModularCfp100G           InterfaceType = "100gbase-x-cfp"
	InterfaceTypeEthernetModularCfp2100G          InterfaceType = "100gbase-x-cfp2"
	InterfaceTypeEthernetModularCfp2200G          InterfaceType = "200gbase-x-cfp2"
	InterfaceTypeEthernetModularCfp4100G          InterfaceType = "100gbase-x-cfp4"
	InterfaceTypeEthernetModularCiscoCpak100G     InterfaceType = "100gbase-x-cpak"
	InterfaceTypeEthernetModularQsfp28100G        InterfaceType = "100gbase-x-qsfp28"
	InterfaceTypeEthernetModularQsfp56200G        InterfaceType = "200gbase-x-qsfp56"
	InterfaceTypeEthernetModularQsfpDD400G        InterfaceType = "400gbase-x-qsfpdd"
	InterfaceTypeEthernetModularOsfp400G          InterfaceType = "400gbase-x-osfp"
	InterfaceTypeWirelessIEEE80211a               InterfaceType = "ieee802.11a"
	InterfaceTypeWirelessIEEE80211bg              InterfaceType = "ieee802.11g"
	InterfaceTypeWirelessIEEE80211n               InterfaceType = "ieee802.11n"
	InterfaceTypeWirelessIEEE80211ac              InterfaceType = "ieee802.11ac"
	InterfaceTypeWirelessIEEE80211ad              InterfaceType = "ieee802.11ad"
	InterfaceTypeWirelessIEEE80211ax              InterfaceType = "ieee802.11ax"
	InterfaceTypeCellularGsm                      InterfaceType = "gsm"
	InterfaceTypeCellularCdma                     InterfaceType = "cdma"
	InterfaceTypeCellularLte                      InterfaceType = "lte"
	InterfaceTypeSonetOc3Stm1                     InterfaceType = "sonet-oc3"
	InterfaceTypeSonetOc12Stm4                    InterfaceType = "sonet-oc12"
	InterfaceTypeSonetOc48Stm16                   InterfaceType = "sonet-oc48"
	InterfaceTypeSonetOc192Stm64                  InterfaceType = "sonet-oc192"
	InterfaceTypeSonetOc768Stm256                 InterfaceType = "sonet-oc768"
	InterfaceTypeSonetOc1920Stm640                InterfaceType = "sonet-oc1920"
	InterfaceTypeSonetOc3840Stm1234               InterfaceType = "sonet-oc3840"
	InterfaceTypeFibreChannelSfp1G                InterfaceType = "1gfc-sfp"
	InterfaceTypeFibreChannelSfp2G                InterfaceType = "2gfc-sfp"
	InterfaceTypeFibreChannelSfp4G                InterfaceType = "4gfc-sfp"
	InterfaceTypeFibreChannelSfpPlus8G            InterfaceType = "8gfc-sfpp"
	InterfaceTypeFibreChannelSfpPlus16G           InterfaceType = "16gfc-sfpp"
	InterfaceTypeFibreChannelSfp2832G             InterfaceType = "32gfc-sfp28"
	InterfaceTypeFibreChannelQsfp28128G           InterfaceType = "128gfc-sfp28"
	InterfaceTypeInifinibandSdr2G                 InterfaceType = "inifiband-sdr"
	InterfaceTypeInifinibandDdr4G                 InterfaceType = "inifiband-ddr"
	InterfaceTypeInifinibandQdr8G                 InterfaceType = "inifiband-qdr"
	InterfaceTypeInifinibandFdr1010G              InterfaceType = "inifiband-fdr10"
	InterfaceTypeInifinibandFdr10135G             InterfaceType = "inifiband-fdr"
	InterfaceTypeInifinibandEdr25G                InterfaceType = "inifiband-edr"
	InterfaceTypeInifinibandHdr50G                InterfaceType = "inifiband-hdr"
	InterfaceTypeInifinibandNdr100G               InterfaceType = "inifiband-ndr"
	InterfaceTypeInifinibandXdr250G               InterfaceType = "inifiband-xdr"
	InterfaceTypeSerialT115M                      InterfaceType = "t1"
	InterfaceTypeSerialE12M                       InterfaceType = "e1"
	InterfaceTypeSerialT3455M                     InterfaceType = "t3"
	InterfaceTypeSerialE334M                      InterfaceType = "e3"
	InterfaceTypeStackingCiscoStackWise           InterfaceType = "cisco-stackwise"
	InterfaceTypeStackingCiscoStackWisePlus       InterfaceType = "cisco-stackwise-plus"
	InterfaceTypeStackingCiscoFlexStack           InterfaceType = "cisco-flexstack"
	InterfaceTypeStackingCiscoFlexStackPlus       InterfaceType = "cisco-flexstack-plus"
	InterfaceTypeStackingJuniperVcp               InterfaceType = "juniper-vcp"
	InterfaceTypeStackingExtremeSummitStack       InterfaceType = "extreme-summitstack"
	InterfaceTypeStackingExtremeSummitStack128    InterfaceType = "extreme-summitstack-128"
	InterfaceTypeStackingExtremeSummitStack256    InterfaceType = "extreme-summitstack-256"
	InterfaceTypeStackingExtremeSummitStack512    InterfaceType = "extreme-summitstack-512"
	InterfaceTypeOther                            InterfaceType = "other"
)

//NetworkInterface represents a network interface assigned to a host
type NetworkInterface struct {
	CommonEntity

	Type         InterfaceType
	Enabled      bool
	IPAddresses  []IPAddress
	IsManagement bool
	MACAddress   net.HardwareAddr
	Name         string
	TaggedVlans  []VLAN
	UntaggedVlan *VLAN
	Tags         []string
	Children     []NetworkInterface
}

func NewNetworkInterface() *NetworkInterface {
	return &NetworkInterface{
		CommonEntity: CommonEntity{
			Meta: &Metadata{},
		},
		Type:         "",
		Enabled:      true,
		IPAddresses:  []IPAddress{},
		IsManagement: false,
		MACAddress:   net.HardwareAddr{},
		Name:         "",
		TaggedVlans:  []VLAN{},
		UntaggedVlan: nil,
		Tags:         []string{},
		Children:     []NetworkInterface{},
	}
}

func (netIf NetworkInterface) String() string {
	return netIf.Name
}

//IsEqual compares the current NetworkInterface object against another NetworkInterface
func (netIf NetworkInterface) IsEqual(netIf2 NetworkInterface) bool {
	if !utils.CompareStruct(netIf, netIf2, []string{}, []string{"CommonEntity", "IPAddresses"}) {
		return false
	}

	//compare length of ip-address slices
	if len(netIf.IPAddresses) != len(netIf2.IPAddresses) {
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
