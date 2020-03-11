package netbox_types

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func MockNetboxVirtualMachineInterface() models.VirtualMachineInterface {
	o := models.VirtualMachineInterface{}

	o.ID = 10
	o.Name = swag.String("eth0")
	o.MacAddress = swag.String("aa:bb:cc:dd:ee:ff")

	o.UntaggedVlan = &models.NestedVLAN{
		ID:   10,
		Vid:  swag.Int64(400),
		Name: swag.String("Public VLAN"),
	}

	o.TaggedVlans = []*models.NestedVLAN{
		{
			ID:   20,
			Vid:  swag.Int64(600),
			Name: swag.String("Private VLAN"),
		},
	}

	return o
}
