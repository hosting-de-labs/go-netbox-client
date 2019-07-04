package types

import (
	"net"
	"sort"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

type InterfaceFormFactor int64

const (
	InterfaceFormFactorVirtualInterfacesVirtual         InterfaceFormFactor = 0
	InterfaceFormFactorVirtualInterfacesLinkAggregation InterfaceFormFactor = 200
	InterfaceFormFactorEthernetFixed100BaseTx_100M      InterfaceFormFactor = 800
	InterfaceFormFactorEthernetFixed1000BaseT_1G        InterfaceFormFactor = 1000
	InterfaceFormFactorEthernetFixedGbic_1G             InterfaceFormFactor = 1050
	InterfaceFormFactorEthernetFixedSfp_1G              InterfaceFormFactor = 1100
	InterfaceFormFactorEthernetFixed10GbaseT_10G        InterfaceFormFactor = 1150
	InterfaceFormFactorEthernetFixed10GbaseCx4_10G      InterfaceFormFactor = 1170
	InterfaceFormFactorEthernetModularSfpPlus_10G       InterfaceFormFactor = 1200
	InterfaceFormFactorEthernetModularXfp_10G           InterfaceFormFactor = 1300
	InterfaceFormFactorEthernetModularXenpak_10G        InterfaceFormFactor = 1310
	InterfaceFormFactorEthernetModularX2_10G            InterfaceFormFactor = 1320
	InterfaceFormFactorEthernetModularSfp28_25G         InterfaceFormFactor = 1350
	InterfaceFormFactorEthernetModularQsfpPlus_40G      InterfaceFormFactor = 1400
	InterfaceFormFactorEthernetModularQsfp28_50G        InterfaceFormFactor = 1420
	InterfaceFormFactorEthernetModularCfp_100G          InterfaceFormFactor = 1500
	InterfaceFormFactorEthernetModularCfp2_100G         InterfaceFormFactor = 1510
	InterfaceFormFactorEthernetModularCfp4_100G         InterfaceFormFactor = 1520
	InterfaceFormFactorEthernetModularCiscoCpak_100G    InterfaceFormFactor = 1550
	InterfaceFormFactorEthernetModularCfp2_200G         InterfaceFormFactor = 1650
	InterfaceFormFactorEthernetModularQsfp28_100G       InterfaceFormFactor = 1600
	InterfaceFormFactorEthernetModularQsfp56_200G       InterfaceFormFactor = 1700
	InterfaceFormFactorEthernetModularQsfpDD_400G       InterfaceFormFactor = 1750
	InterfaceFormFactorWirelessIEEE80211a               InterfaceFormFactor = 2600
	InterfaceFormFactorWirelessIEEE80211bg              InterfaceFormFactor = 2610
	InterfaceFormFactorWirelessIEEE80211n               InterfaceFormFactor = 2620
	InterfaceFormFactorWirelessIEEE80211ac              InterfaceFormFactor = 2630
	InterfaceFormFactorWirelessIEEE80211ad              InterfaceFormFactor = 2640
	InterfaceFormFactorCellularGsm                      InterfaceFormFactor = 2810
	InterfaceFormFactorCellularCdma                     InterfaceFormFactor = 2820
	InterfaceFormFactorCellularLte                      InterfaceFormFactor = 2830
	InterfaceFormFactorSonetOc3Stm1                     InterfaceFormFactor = 6100
	InterfaceFormFactorSonetOc12Stm4                    InterfaceFormFactor = 6200
	InterfaceFormFactorSonetOc48Stm16                   InterfaceFormFactor = 6300
	InterfaceFormFactorSonetOc192Stm64                  InterfaceFormFactor = 6400
	InterfaceFormFactorSonetOc768Stm256                 InterfaceFormFactor = 6500
	InterfaceFormFactorSonetOc1920Stm640                InterfaceFormFactor = 6600
	InterfaceFormFactorSonetOc3840Stm1234               InterfaceFormFactor = 6700
	InterfaceFormFactorFibreChannelSfp_1G               InterfaceFormFactor = 3010
	InterfaceFormFactorFibreChannelSfp_2G               InterfaceFormFactor = 3020
	InterfaceFormFactorFibreChannelSfp_4G               InterfaceFormFactor = 3040
	InterfaceFormFactorFibreChannelSfpPlus_8G           InterfaceFormFactor = 3080
	InterfaceFormFactorFibreChannelSfpPlus_16G          InterfaceFormFactor = 3160
	InterfaceFormFactorFibreChannelSfp28_32G            InterfaceFormFactor = 3320
	InterfaceFormFactorFibreChannelQsfp28_128G          InterfaceFormFactor = 3400
	InterfaceFormFactorSerialT1_1_5M                    InterfaceFormFactor = 4000
	InterfaceFormFactorSerialE1_2M                      InterfaceFormFactor = 4010
	InterfaceFormFactorSerialT3_455M                    InterfaceFormFactor = 4040
	InterfaceFormFactorSerialE3_34M                     InterfaceFormFactor = 4050
	InterfaceFormFactorStackingCiscoStackWise           InterfaceFormFactor = 5000
	InterfaceFormFactorStackingCiscoStackWisePlus       InterfaceFormFactor = 5050
	InterfaceFormFactorStackingCiscoFlexStack           InterfaceFormFactor = 5100
	InterfaceFormFactorStackingCiscoFlexStackPlus       InterfaceFormFactor = 5150
	InterfaceFormFactorStackingJuniperVcp               InterfaceFormFactor = 5200
	InterfaceFormFactorStackingExtremeSummitStack       InterfaceFormFactor = 5300
	InterfaceFormFactorStackingExtremeSummitStack_128   InterfaceFormFactor = 5310
	InterfaceFormFactorStackingExtremeSummitStack_256   InterfaceFormFactor = 5320
	InterfaceFormFactorStackingExtremeSummitStack_512   InterfaceFormFactor = 5330
	InterfaceFormFactorOther                            InterfaceFormFactor = 32767
)

//NetworkInterface represents a network interface assigned to a host
type NetworkInterface struct {
	CommonEntity

	FormFactor   *InterfaceFormFactor
	IPAddresses  []IPAddress
	IsManagement bool
	MACAddress   net.HardwareAddr
	Name         string
	TaggedVlans  []VLAN
	UntaggedVlan *VLAN
	Tags         []string
	Children     []NetworkInterface
}

//IsEqual compares the current NetworkInterface object against another NetworkInterface
func (netIf NetworkInterface) IsEqual(netIf2 NetworkInterface) bool {
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
