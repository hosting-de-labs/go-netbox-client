package mock_dcim

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

func DeviceGetResponder() (string, string, httpmock.Responder) {
	return http.MethodGet,
		`=~http://localhost:8000\/api\/dcim\/devices\/10\/`,
		httpmock.NewStringResponder(
			200,
			`{
    "id": 10,
    "name": "Dummy Device 10",
    "display_name": "Dummy Device 10",
    "device_type": {
        "id": 1,
        "url": "https://netbox.mgmt.masterlogin.de/api/dcim/device-types/1/",
        "manufacturer": {
            "id": 1,
            "url": "https://netbox.mgmt.masterlogin.de/api/dcim/manufacturers/1/",
            "name": "Generic",
            "slug": "generic"
        },
        "model": "device",
        "slug": "device",
        "display_name": "Generic Device"
    },
    "device_role": {
        "id": 1,
        "url": "https://netbox.mgmt.masterlogin.de/api/dcim/device-roles/1/",
        "name": "Server",
        "slug": "server"
    },
    "tenant": null,
    "platform": null,
    "serial": "",
    "asset_tag": "",
    "site": {
        "id": 1,
        "url": "https://netbox.mgmt.masterlogin.de/api/dcim/sites/1/",
        "name": "SITE1",
        "slug": "site1"
    },
    "rack": {
        "id": 1,
        "url": "https://netbox.mgmt.masterlogin.de/api/dcim/racks/1/",
        "name": "R1",
        "display_name": "Rack 1"
    },
    "position": null,
    "face": null,
    "parent_device": null,
    "status": {
        "value": 1,
        "label": "Active"
    },
    "primary_ip": null,
    "primary_ip4": null,
    "primary_ip6": null,
    "cluster": null,
    "virtual_chassis": null,
    "vc_position": null,
    "vc_priority": null,
    "comments": "",
    "local_context_data": null,
    "tags": [],
    "custom_fields": null,
    "config_context": {},
    "created": "2012-01-01",
    "last_updated": "2012-01-01T00:00:00.000000Z"
}`)
}
