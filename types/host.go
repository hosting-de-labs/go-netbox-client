package types

import "sort"

//Host represents a host
type Host struct {
	OriginalHost *Host

	ID          int64
	Hostname    string
	PrimaryIPv4 IPAddress
	PrimaryIPv6 IPAddress
	IsManaged   bool
	Tags        []string
	Comments    []string

	NetworkInterfaces []HostNetworkInterface
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
func (h Host) Copy() *Host {
	out := Host{
		ID:          h.ID,
		Hostname:    h.Hostname,
		IsManaged:   h.IsManaged,
		PrimaryIPv4: h.PrimaryIPv4,
		PrimaryIPv6: h.PrimaryIPv6,
	}

	//copy interfaces
	out.NetworkInterfaces = make([]HostNetworkInterface, len(h.NetworkInterfaces))
	copy(out.NetworkInterfaces, h.NetworkInterfaces)

	//copy tags
	out.Tags = make([]string, len(h.Tags))
	copy(out.Tags, h.Tags)

	return &out
}

//IsChanged returns true if the current and the original object differ
func (h Host) IsChanged() bool {
	return h.IsEqual(*h.OriginalHost, true)
}

//IsEqual compares the current object against another Host object
func (h Host) IsEqual(h2 Host, deep bool) bool {
	if h.Hostname != h2.Hostname {
		return false
	}

	if h.PrimaryIPv4 != h2.PrimaryIPv4 || h.PrimaryIPv6 != h2.PrimaryIPv6 {
		return false
	}

	if h.IsManaged != h2.IsManaged {
		return false
	}

	//comments
	if len(h.Comments) != len(h2.Comments) {
		return false
	}

	for i := 0; i < len(h.Comments); i++ {
		if h.Comments[i] != h2.Comments[i] {
			return false
		}
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

		//TODO: sort!

		for i := 0; i < len(h.NetworkInterfaces); i++ {
			if !h.NetworkInterfaces[i].IsEqual(h2.NetworkInterfaces[i]) {
				return false
			}
		}
	}

	return true
}
