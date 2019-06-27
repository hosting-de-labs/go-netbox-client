package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c Client) DeviceFindAll(limit int64, offset int64) (int64, []*models.Device, error) {
	params := dcim.NewDcimDevicesListParams()

	if limit > 0 {
		params.WithLimit(&limit)
	}

	if offset > 0 {
		params.WithOffset(&offset)
	}

	res, err := c.client.Dcim.DcimDevicesList(params, nil)

	if err != nil {
		return 0, nil, err
	}

	return *res.Payload.Count, res.Payload.Results, nil
}

func (c Client) DeviceGet(deviceID int64) (*models.DeviceWithConfigContext, error) {
	params := dcim.NewDcimDevicesReadParams()
	params.WithID(deviceID)

	res, err := c.client.Dcim.DcimDevicesRead(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

//DeviceFindByHostname retrieves a model.Device object from Netbox by looking up the given hostname
//TODO: return *types.DedicatedServer
func (c Client) DeviceFindByHostname(hostname string) (*models.Device, error) {
	params := dcim.NewDcimDevicesListParams()
	params.WithName(&hostname)

	res, err := c.client.Dcim.DcimDevicesList(params, nil)

	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, fmt.Errorf("Hostname with name %s not found", hostname)
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("Hostname %s is not unique", hostname)
	}

	return res.Payload.Results[0], nil
}

//TODO: Don't forget to cache
func (c Client) DeviceUpdate(host *types.DedicatedServer) error {
	if host.OriginalEntity == nil {
		oh, err := c.DeviceFindByHostname(host.Hostname)
		if err != nil {
			return fmt.Errorf("cannot update DedicatedServer %s. No OriginalHost assigned and no way to find a matching one", host.Hostname)
		}

		res, err := c.DeviceConvertFromNetbox(oh)
		if err != nil {
			return fmt.Errorf("cannot convert to DedicatedServer")
		}

		host.OriginalEntity = res
	}

	data := new(models.WritableDevice)

	//Go through every item and check if it must be updated

	params := dcim.NewDcimDevicesUpdateParams()
	params.WithID(host.OriginalEntity.(types.DedicatedServer).ID).WithData(data)

	//TODO: Iterate over Inventory Items

	return nil
}

//HypervisorFindByHostname is like NetboxDeviceGet but checks if the device has a cluster assigned
func (c Client) HypervisorFindByHostname(hostname string) (*models.Device, error) {
	res, err := c.DeviceFindByHostname(hostname)

	if err != nil {
		return nil, err
	}

	if res.Cluster == nil {
		return nil, fmt.Errorf("device %s not assigned to a Virtualization Cluster", hostname)
	}

	return res, nil
}
