package netbox_types

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func MockNetboxDevice(addIPAddresses bool, addTags bool) (out models.Device) {
	out = models.Device{
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
	}

	if addIPAddresses {
		//TODO: add interfaces when adding ip addresses
		out.PrimaryIp4 = &models.NestedIPAddress{Address: swag.String("127.0.0.1/32")}
		out.PrimaryIp6 = &models.NestedIPAddress{Address: swag.String("::1/128")}
	}

	if addTags {
		out.Tags = append(out.Tags, "Tag1")
		out.Tags = append(out.Tags, "managed")
	}

	return out
}
