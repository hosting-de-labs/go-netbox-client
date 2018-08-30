package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxCache "internal.keenlogics.com/di/netbox-sync/cache"
)

var cache netboxCache.Cache

//DeviceGet retrieves a model.Device object from Netbox by looking up the given hostname
func DeviceGet(netboxClient *client.NetBox, hostname string) (*models.Device, error) {
	it, ok := cache.Get(fmt.Sprintf("DCIM_DEVICE_BY_HOSTNAME_%s", hostname))
	if ok {
		return it.(*models.Device), nil
	}

	params := dcim.NewDcimDevicesListParams()
	params.WithName(&hostname)

	res, err := netboxClient.Dcim.DcimDevicesList(params, nil)

	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, fmt.Errorf("Hostname with name %s not found", hostname)
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("Hostname %s is not unique", hostname)
	}

	//store device in cache
	cache.Set(fmt.Sprintf("DCIM_DEVICE_BY_HOSTNAME_%s", hostname), res.Payload.Results[0])

	return res.Payload.Results[0], nil
}

//HypervisorGet is like NetboxDeviceGet but checks if the device has a cluster assigned
func HypervisorGet(netboxClient *client.NetBox, hostname string) (*models.Device, error) {
	res, err := DeviceGet(netboxClient, hostname)

	if err != nil {
		return nil, err
	}

	if res.Cluster == nil {
		return nil, fmt.Errorf("Device %s not assigned to a Virtualization Cluster", hostname)
	}

	return res, nil
}
