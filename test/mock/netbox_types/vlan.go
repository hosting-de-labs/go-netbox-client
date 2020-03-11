package netbox_types

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func MockNetboxVlan() models.VLAN {
	return models.VLAN{
		ID:   10,
		Vid:  swag.Int64(400),
		Name: swag.String("Public VLAN"),
		Status: &models.VLANStatus{
			Value: swag.String("active"),
			Label: swag.String("Active"),
		},
		Description: "This is Public VLAN description",
		Tags:        []string{"public"},
	}
}

func MockNetboxNestedVlan() models.NestedVLAN {
	return models.NestedVLAN{
		ID:   20,
		Vid:  swag.Int64(600),
		Name: swag.String("Private VLAN"),
	}
}
