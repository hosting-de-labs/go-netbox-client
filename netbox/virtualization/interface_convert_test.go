package virtualization_test

import (
	"net"
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/test/mock"

	"github.com/hosting-de-labs/go-netbox-client/netbox/virtualization"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
	"github.com/jarcoal/httpmock"
)

func init() {
	mock.RunServer()
}

func mockNetboxVirtualMachineInterface() models.VirtualMachineInterface {
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

func TestConvertVirtualMachineInterface(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		`=~http://localhost:8000\/api\/ipam\/ip-addresses\/\?interface_id=10`,
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
				"value": 1,
				"label": "Active"
			},
			"role": null,
			"interface": {
				"id": 10,
                "url": "http://localhost:8000/api/virtualization/interfaces/10/",
                "device": null,
				"form_factor": {
					"value": 1000
				},
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
				"value": 1,
				"label": "Active"
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
}`,
		),
	)

	netboxClient := netbox.NewNetboxAt("localhost:8000")
	virtualizationClient := virtualization.NewClient(*netboxClient)

	netIf, err := virtualizationClient.InterfaceConvertFromNetbox(mockNetboxVirtualMachineInterface())

	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, netIf.Name, "eth0")

	mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	assert.Nil(t, err)

	assert.Equal(t, netIf.MACAddress, mac)

	assert.Equal(t, len(netIf.IPAddresses), 2)

	assert.Equal(t, "123.123.123.123", netIf.IPAddresses[0].Address)
	assert.Equal(t, uint16(24), netIf.IPAddresses[0].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv4, netIf.IPAddresses[0].Family)

	assert.Equal(t, "2001:db8:a::123", netIf.IPAddresses[1].Address)
	assert.Equal(t, uint16(64), netIf.IPAddresses[1].CIDR)
	assert.Equal(t, types.IPAddressFamilyIPv6, netIf.IPAddresses[1].Family)
}
