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

func NewHost() *Host {
	return &Host{
		CommonEntity: CommonEntity{
			Meta: &Metadata{},
		},
		Hostname: "",
	}
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
		tagFound := false
		for _, existingTag := range h.Tags {
			if existingTag == newTag {
				tagFound = true
				break
			}
		}

		if !tagFound {
			h.Tags = append(h.Tags, newTag)
		}
	}
}

//Copy creates a deep copy of the given host
func (h Host) Copy() Host {
	out := NewHost()
	out.Hostname = h.Hostname
	out.IsManaged = h.IsManaged
	out.PrimaryIPv4 = h.PrimaryIPv4
	out.PrimaryIPv6 = h.PrimaryIPv6

	if h.Meta != nil {
		out.Meta = &Metadata{
			ID:             h.Meta.ID,
			OriginalEntity: h.Meta.OriginalEntity,
			NetboxEntity:   h.Meta.NetboxEntity,
		}
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

	return *out
}

//IsChanged returns true if the current and the original object differ
func (h Host) IsChanged() bool {
	if h.Meta == nil || h.Meta.OriginalEntity == nil {
		return true
	}

	return !h.IsEqual(h.Meta.OriginalEntity.(Host), true)
}

//IsEqual compares the current object against another Host object
func (h Host) IsEqual(h2 Host, deep bool) bool {
	h.Meta = nil
	h2.Meta = nil

	if !utils.CompareStruct(h, h2, []string{}, []string{"Meta", "NetworkInterfaces", "Tags"}) {
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
