package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxCache "github.com/hosting-de-labs/go-netbox-client/cache"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

var cache netboxCache.Cache

//DeviceGet retrieves a model.Device object from Netbox by looking up the given hostname
func (c *Client) DeviceGet(hostname string) (*models.Device, error) {
	it, ok := cache.Get(fmt.Sprintf("DCIM_DEVICE_BY_HOSTNAME_%s", hostname))
	if ok {
		return it.(*models.Device), nil
	}

	params := dcim.NewDcimDevicesListParams()
	params.WithName(&hostname)

	res, err := c.getClient().Dcim.DcimDevicesList(params, nil)

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

//TODO: Don't forget to cache
func (c *Client) DeviceUpdate(host *types.DedicatedServer) error {
	if host.OriginalHost == nil {
		oh, err := c.DeviceGet(host.Hostname)
		if err != nil {
			return fmt.Errorf("cannot update DedicatedServer %s. No OriginalHost assigned and no way to find a matching one", host.Hostname)
		}

		res, err := c.DeviceConvertFromNetbox(oh)
		if err != nil {
			return fmt.Errorf("cannot convert to DedicatedServer")
		}

		host.OriginalHost = res
	}

	data := new(models.WritableDevice)

	//Go through every item and check if it must be updated

	params := dcim.NewDcimDevicesUpdateParams()
	params.WithID(host.OriginalHost.ID).WithData(data)

	//Iterate over Inventory Items

	return nil
}

//HypervisorGet is like NetboxDeviceGet but checks if the device has a cluster assigned
func (c *Client) HypervisorGet(hostname string) (*models.Device, error) {
	res, err := c.DeviceGet(hostname)

	if err != nil {
		return nil, err
	}

	if res.Cluster == nil {
		return nil, fmt.Errorf("Device %s not assigned to a Virtualization Cluster", hostname)
	}

	return res, nil
}
