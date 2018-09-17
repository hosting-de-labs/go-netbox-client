package types

//DedicatedServer represents a dedicated server
type DedicatedServer struct {
	Host
	OriginalHost *DedicatedServer

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

	//TODO: sort inventory items

	//iterate through inventory items and compare each item using IsEqual
	for key, item1 := range d.Inventory {
		if !item1.IsEqual(*d2.Inventory[key]) {
			return false
		}
	}

	return true
}
