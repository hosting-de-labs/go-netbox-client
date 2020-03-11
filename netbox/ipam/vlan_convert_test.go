package ipam

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_types"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"
)

func TestVlanConvertFromNetbox(t *testing.T) {
	vlan400, err := VlanConvertFromNetbox(netbox_types.MockNetboxVlan())
	assert.Nil(t, err)

	assert.Equal(t, uint16(400), vlan400.ID)
	assert.Equal(t, "Public VLAN", vlan400.Name)

	assert.Equal(t, types.VLANStatus(types.VLANStatusActive), vlan400.Status)

	assert.Equal(t, "This is Public VLAN description", vlan400.Description)
	assert.Equal(t, []string{"public"}, vlan400.Tags)
}

func TestVlanConvertFromNetboxWithNestedVlan(t *testing.T) {
	vlan600, err := VlanConvertFromNetbox(netbox_types.MockNetboxNestedVlan())
	assert.Nil(t, err)

	assert.Equal(t, uint16(600), vlan600.ID)
	assert.Equal(t, "Private VLAN", vlan600.Name)
}
