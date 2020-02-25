package mock_dcim

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

func InterfaceGetResponder() (string, string, httpmock.Responder) {
	return http.MethodGet,
		`=~http://localhost:8000\/api\/dcim\/interfaces\/10\/`,
		httpmock.NewStringResponder(
			200,
			`{
    "id": 10,
    "device": {
        "id": 20,
        "url": "http://localhost:8000/api/dcim/devices/20/",
        "name": "device_20",
        "display_name": "Device 20"
    },
    "name": "eth0",
    "type": {
        "value": "1000base-t",
        "label": "1000BASE-T (1GE)",
		"id": 1000
    },
    "enabled": true,
    "lag": null,
    "mtu": null,
    "mac_address": "aa:bb:cc:dd:ee:ff",
    "mgmt_only": false,
    "description": "",
    "connected_endpoint_type": null,
    "connected_endpoint": null,
    "connection_status": null,
    "cable": null,
    "mode": null,
    "untagged_vlan": {
        "id": 30,
        "url": "http://localhost:8000/api/ipam/vlans/30/",
        "vid": 400,
        "name": "vlan-400",
        "display_name": "400 (vlan-400)"
    },
    "tagged_vlans": [],
    "tags": []
}`)
}
