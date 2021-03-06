package types_test

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/stretchr/testify/assert"
)

func MockInventoryItem() types.InventoryItem {
	return types.InventoryItem{
		Type:         types.InventoryItemTypeProcessor,
		Manufacturer: "Intel",
		Model:        "Xeon X5670",
		AssetTag:     "Asset Tag",
		PartNumber:   "Part Number",
		SerialNumber: "Serial Number",
	}
}

func TestInventoryItem_GetHashableString(t *testing.T) {
	item := MockInventoryItem()
	assert.Equal(t, item.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber")

	item.AddDetail("Cores", "2")
	assert.Equal(t, item.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber:details{Cores:2}")

	item.AddDetail("Threads", "4")
	assert.Equal(t, item.GetHashableString(), "Intel:XeonX5670:PartNumber:AssetTag:SerialNumber:details{Cores:2,Threads:4}")
}

func TestInventoryItem_AddDetail(t *testing.T) {
	item := MockInventoryItem()

	assert.Nil(t, item.Details)

	item.AddDetail("cores", "2")

	assert.NotNil(t, item.Details)
	assert.NotEmpty(t, item.Details)

	val, ok := item.Details["cores"]

	assert.True(t, ok)
	assert.Equal(t, val, "2")
}

func TestInventoryItem_Copy(t *testing.T) {
	item := MockInventoryItem()
	item2 := item.Copy()

	assert.Equal(t, item, item2)
	assert.True(t, item.IsEqual(item2))
}

func TestInventoryItem_IsEqual(t *testing.T) {
	item := MockInventoryItem()
	item.AddDetail("Cores", "2")
	item.AddDetail("Threads", "4")

	item2 := item.Copy()
	assert.Equal(t, item, item2)

	item2.AddDetail("L3 Cache", "12MB")
	assert.NotEqual(t, item, item2)
}
