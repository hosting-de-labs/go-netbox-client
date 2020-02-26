package types

import (
	"sort"

	"github.com/hosting-de-labs/go-netbox-client/utils"
)

//Host represents a host
type Host struct {
	CommonEntity

	Hostname    string
	PrimaryIPv4 IPAddress
	PrimaryIPv6 IPAddress
	IsManaged   bool
	Tags        []string
	Comments    []string

	NetworkInterfaces []NetworkInterface
}

//HasTag checks for a specific tag being assigned to this host
func (h *Host) HasTag(tag string) bool {
	for _, existingTag := range h.Tags {
		if existingTag == tag {
			return true
		}
	}

	return false
}

//AddTag is a helper method to allow adding a number of tags to a host.
func (h *Host) AddTag(tags ...string) {
	for _, newTag := range tags {
		for _, existingTag := range h.Tags {
			if existingTag == newTag {
				break
			}
		}

		h.Tags = append(h.Tags, newTag)
	}
}

//Copy creates a deep copy of the given host
func (h Host) Copy() Host {
	out := Host{
		Hostname:    h.Hostname,
		IsManaged:   h.IsManaged,
		PrimaryIPv4: h.PrimaryIPv4,
		PrimaryIPv6: h.PrimaryIPv6,
	}

	//copy comments
	if len(h.Comments) > 0 {
		out.Comments = make([]string, len(h.Comments))
		copy(out.Comments, h.Comments)
	}

	//copy interfaces
	if len(h.NetworkInterfaces) > 0 {
		out.NetworkInterfaces = make([]NetworkInterface, len(h.NetworkInterfaces))
		copy(out.NetworkInterfaces, h.NetworkInterfaces)
	}

	//copy tags
	if len(h.Tags) > 0 {
		out.Tags = make([]string, len(h.Tags))
		copy(out.Tags, h.Tags)
	}

	return out
}

//IsChanged returns true if the current and the original object differ
func (h Host) IsChanged() bool {
	return !h.IsEqual(h.OriginalEntity.(Host), true)
}

//IsEqual compares the current object against another Host object
func (h Host) IsEqual(h2 Host, deep bool) bool {
	h.OriginalEntity = nil
	h2.OriginalEntity = nil

	if !utils.CompareStruct(h, h2, []string{}, []string{"Metadata", "NetworkInterfaces", "OriginalEntity", "Tags"}) {
		return false
	}

	//tags
	if len(h.Tags) != len(h2.Tags) {
		return false
	}

	sort.Strings(h.Tags)
	sort.Strings(h2.Tags)

	for i := 0; i < len(h.Tags); i++ {
		if h.Tags[i] != h2.Tags[i] {
			return false
		}
	}

	if deep {
		//network interfaces
		if len(h.NetworkInterfaces) != len(h2.NetworkInterfaces) {
			return false
		}

		//sort interfaces by name
		sort.Slice(h.NetworkInterfaces, func(i, j int) bool { return h.NetworkInterfaces[i].Name < h.NetworkInterfaces[j].Name })
		sort.Slice(h2.NetworkInterfaces, func(i, j int) bool { return h2.NetworkInterfaces[i].Name < h2.NetworkInterfaces[j].Name })

		for i := 0; i < len(h.NetworkInterfaces); i++ {
			if !h.NetworkInterfaces[i].IsEqual(h2.NetworkInterfaces[i]) {
				return false
			}
		}
	}

	return true
}
