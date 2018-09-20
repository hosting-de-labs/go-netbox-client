package types

import "sort"

//DedicatedServer represents a dedicated server
type DedicatedServer struct {
	Host

	Inventory []*InventoryItem
}

//Copy creates a deep copy of the given host
func (d DedicatedServer) Copy() *DedicatedServer {
	out := new(DedicatedServer)

	//Copy Host
	out.Host = *d.Host.Copy()

	//Copy Inventory
	out.Inventory = make([]*InventoryItem, len(d.Inventory), len(d.Inventory))
	for _, item := range d.Inventory {
		out.Inventory = append(out.Inventory, item.Copy())
	}

	return out
}

//IsEqual compares the current object with another VirtualServer object
func (d DedicatedServer) IsEqual(d2 DedicatedServer, deep bool) bool {
	//compare Host struct
	if !d.Host.IsEqual(d2.Host, deep) {
		return false
	}

	//compare length of inventory items
	if len(d.Inventory) != len(d2.Inventory) {
		return false
	}

	//sort inventory items
	sort.Slice(d.Inventory, func(i, j int) bool { return d.Inventory[i].GetHashableString() < d.Inventory[j].GetHashableString() })
	sort.Slice(d2.Inventory, func(i, j int) bool { return d2.Inventory[i].GetHashableString() < d2.Inventory[j].GetHashableString() })

	//iterate through inventory items and compare each item using IsEqual
	for key, item1 := range d.Inventory {
		if !item1.IsEqual(*d2.Inventory[key]) {
			return false
		}
	}

	return true
}
