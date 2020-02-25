package mock_ipam

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

func IpAddressFindResponder() (string, string, httpmock.Responder) {
	return http.MethodGet,
		`=~http://localhost:8000\/api\/ipam\/ip-addresses\/`,
		httpmock.NewStringResponder(
			200,
			`{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
		{
			"id": 1234,
			"family": {
				"value": 4,
				"label": "IPv4"
			},
			"address": "123.123.123.123/24",
			"vrf": null,
			"tenant": null,
			"status": {
				"value": "active",
				"label": "Active",
				"id": 1
			},
			"role": null,
			"interface": {
				"id": 10,
                "url": "http://localhost:8000/api/virtualization/interfaces/10/",
                "device": null,
                "virtual_machine": {
                    "id": 5651,
                    "url": "http://localhost:8000/api/virtualization/virtual-machines/15/",
                    "name": "Virtual Machine 1"
                },
                "name": "eth0"
			},
			"description": "",
			"nat_inside": null,
			"nat_outside": null,
			"tags": [],
			"custom_fields": {},
			"created": "2019-01-01",
			"last_updated": "2019-01-01T12:30:00.000000Z"
		},
        {
			"id": 1235,
			"family": {
				"value": 6,
				"label": "IPv6"
			},
			"address": "2001:db8:a::123/64",
			"vrf": null,
			"tenant": null,
			"status": {
				"value": "active",
				"label": "Active",
				"id": 1
			},
			"role": null,
			"interface": {
				"id": 10,
                "url": "http://localhost:8000/api/virtualization/interfaces/10/",
                "device": null,
                "virtual_machine": {
                    "id": 5651,
                    "url": "http://localhost:8000/api/virtualization/virtual-machines/15/",
                    "name": "Virtual Machine 1"
                },
                "name": "eth0"
			},
			"description": "",
			"nat_inside": null,
			"nat_outside": null,
			"tags": [],
			"custom_fields": {},
			"created": "2019-01-01",
			"last_updated": "2019-01-01T12:30:00.000000Z"
		}
    ]
	}`)
}
