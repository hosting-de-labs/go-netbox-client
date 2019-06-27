package dcim

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestClient_InterfaceGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
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
    "form_factor": {
        "value": 1000,
        "label": "1000BASE-T (1GE)"
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
}`,
		),
	)

	netboxClient := netbox.NewNetboxAt("localhost:8000")

	dcimClient := NewClient(*netboxClient)
	netIf, err := dcimClient.InterfaceGet(10)

	assert.Nil(t, err)
	assert.NotNil(t, netIf)

	assert.Equal(t, "eth0", *netIf.Name)
	assert.Equal(t, types.InterfaceFormFactorEthernetFixed1000BaseT_1G, types.InterfaceFormFactor(*netIf.FormFactor.Value))
	assert.Equal(t, "aa:bb:cc:dd:ee:ff", *netIf.MacAddress)

	assert.NotNil(t, netIf.UntaggedVlan)
	assert.Equal(t, int64(30), netIf.UntaggedVlan.ID)
	assert.Equal(t, int64(400), *netIf.UntaggedVlan.Vid)
	assert.Equal(t, "vlan-400", *netIf.UntaggedVlan.Name)
	assert.Equal(t, "400 (vlan-400)", netIf.UntaggedVlan.DisplayName)

	assert.Empty(t, netIf.TaggedVlans)
	assert.Empty(t, netIf.Tags)
}
