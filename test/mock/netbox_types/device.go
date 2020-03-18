package netbox_types

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func MockNetboxDevice() (out models.Device) {
	return models.Device{
		AssetTag:    swag.String("123-456"),
		Created:     strfmt.Date(time.Now()),
		DisplayName: "",
		ID:          10,
		LastUpdated: strfmt.DateTime(time.Now()),
		Name:        swag.String("Host 10"),
		Serial:      "1234567890",
		Status: &models.DeviceStatus{
			Label: swag.String("Active"),
			Value: swag.String("active"),
		},
		Tags: []string{
			"Tag1",
			"Tag2",
		},
	}
}
