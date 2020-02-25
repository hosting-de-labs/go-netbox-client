package ipam

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func VlanConvertFromNetbox(netboxVlan interface{}) (*types.VLAN, error) {
	vlan := types.VLAN{}
	var vlanStatus string

	switch v := netboxVlan.(type) {
	case models.VLAN:
		vlan.ID = uint16(*v.Vid)
		vlan.Name = *v.Name
		vlan.Description = v.Description
		vlan.Tags = v.Tags
		vlanStatus = *v.Status.Value
	case models.NestedVLAN:
		vlan.ID = uint16(*v.Vid)
		vlan.Name = *v.Name
	default:
		return nil, fmt.Errorf("vlan has to be of type VLAN oder NestedVLAN, type is %T", netboxVlan)
	}

	switch vlanStatus {
	case "":
		fallthrough
	case "unknown":
		vlan.Status = types.VLANStatusUnknown
	case "active":
		vlan.Status = types.VLANStatusActive
	case "reserved":
		vlan.Status = types.VLANStatusReserved
	case "deprecated":
		vlan.Status = types.VLANStatusDeprecated

	default:
		return nil, fmt.Errorf("unknown vlan status %s", vlanStatus)
	}

	return &vlan, nil
}
