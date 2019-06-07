package types

import "github.com/hosting-de-labs/go-netbox-client/utils"

type VLANStatus int

const (
	VLANStatusUnknown = iota
	VLANStatusActive
	VLANStatusReserved
	VLANStatusDeprecated
)

type VLAN struct {
	ID          uint16
	Name        string
	Status      VLANStatus
	Description string
	Tags        []string
}

func (v VLAN) Clone() (out VLAN) {
	out = VLAN{
		ID:          v.ID,
		Name:        v.Name,
		Status:      v.Status,
		Description: v.Description,
	}

	if len(v.Tags) > 0 {
		out.Tags = make([]string, len(v.Tags))
		copy(out.Tags, v.Tags)
	}

	return out
}

//IsEqual compares the current IPAddress object against another IPAddress object
func (v VLAN) IsEqual(v2 VLAN) bool {
	return utils.CompareStruct(v, v2, []string{}, []string{})
}
