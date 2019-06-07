package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockInventoryItem() InventoryItem {
	return InventoryItem{
		Type:         InventoryItemTypeProcessor,
		Manufacturer: "Intel",
		Model:        "Xeon X5670",
		AssetTag:     "Asset Tag",
		PartNumber:   "Part Number",
		SerialNumber: "Serial Number",
	}
}

func TestInventoryItem_GetHashableString(t *testing.T) {
	item1 := mockInventoryItem()
	assert.Equal(t, item1.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber")

	item1.AddDetail("Cores", "2")
	assert.Equal(t, item1.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber:details{Cores:2}")

	item1.AddDetail("Threads", "4")
	assert.Equal(t, item1.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber:details{Cores:2,Threads:4}")
}

func TestInventoryItem_AddDetail(t *testing.T) {
	item := mockInventoryItem()

	assert.Nil(t, item.Details)

	item.AddDetail("cores", "2")

	assert.NotNil(t, item.Details)
	assert.NotEmpty(t, item.Details)

	val, ok := item.Details["cores"]

	assert.True(t, ok)
	assert.Equal(t, val, "2")
}

func TestInventoryItem_Copy(t *testing.T) {
	item := mockInventoryItem()
	item2 := item.Copy()

	assert.Equal(t, item, item2)
	assert.True(t, item.IsEqual(item2))
}

func TestInventoryItem_IsEqual(t *testing.T) {
	item := mockInventoryItem()
	item.AddDetail("Cores", "2")
	item.AddDetail("Threads", "4")

	item2 := item.Copy()
	assert.Equal(t, item, item2)

	item.AddDetail("L3 Cache", "12MB")
	assert.NotEqual(t, item, item2)
}
