package ipam

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/stretchr/testify/assert"
)

func mockNetboxVlan() models.VLAN {
	return models.VLAN{
		ID:   10,
		Vid:  swag.Int64(400),
		Name: swag.String("Public VLAN"),
		Status: &models.VLANStatus{
			Value: swag.String("active"),
			Label: swag.String("Enabled"),
		},
		Description: "This is Public VLAN description",
	}
}

func mockNetboxNestedVlan() models.NestedVLAN {
	return models.NestedVLAN{
		ID:   20,
		Vid:  swag.Int64(600),
		Name: swag.String("Private VLAN"),
	}
}

func TestVlanConvertFromNetbox(t *testing.T) {
	vlan400, err := VlanConvertFromNetbox(mockNetboxVlan())
	assert.Nil(t, err)

	assert.Equal(t, uint16(400), vlan400.ID)
	assert.Equal(t, "Public VLAN", vlan400.Name)

	assert.Equal(t, types.VLANStatus(types.VLANStatusActive), vlan400.Status)

	assert.Equal(t, "This is Public VLAN description", vlan400.Description)
}

func TestVlanConvertFromNetboxWithNestedVlan(t *testing.T) {
	vlan600, err := VlanConvertFromNetbox(mockNetboxNestedVlan())
	assert.Nil(t, err)

	assert.Equal(t, uint16(600), vlan600.ID)
	assert.Equal(t, "Private VLAN", vlan600.Name)
}
