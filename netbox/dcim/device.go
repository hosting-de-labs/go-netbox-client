package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c Client) DeviceFindAll(limit int64, offset int64) (out []types.DedicatedServer, err error) {
	params := dcim.NewDcimDevicesListParams()

	if limit > 0 {
		params.WithLimit(&limit)
	}

	if offset > 0 {
		params.WithOffset(&offset)
	}

	res, err := c.client.Dcim.DcimDevicesList(params, nil)
	if err != nil {
		return nil, err
	}

	for _, nbDevice := range res.Payload.Results {
		intf, err := c.DeviceConvertFromNetbox(*nbDevice)
		if err != nil {
			return nil, err
		}

		out = append(out, *intf)
	}

	return out, nil
}

//DeviceFind retrieves a model.Device object from Netbox by looking up the given hostname
func (c Client) DeviceFind(hostname string) (out *types.DedicatedServer, err error) {
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

	return c.DeviceConvertFromNetbox(*res.Payload.Results[0])
}

func (c Client) DeviceGet(deviceID int64) (*types.DedicatedServer, error) {
	params := dcim.NewDcimDevicesReadParams()
	params.WithID(deviceID)

	res, err := c.client.Dcim.DcimDevicesRead(params, nil)
	if err != nil {
		return nil, err
	}

	return c.DeviceConvertFromNetbox(*res.Payload)
}

//TODO: Don't forget to cache
func (c Client) DeviceUpdate(host *types.DedicatedServer) error {
	if host.Meta.NetboxEntity == nil {
		oh, err := c.DeviceFind(host.Hostname)
		if err != nil {
			return fmt.Errorf("cannot update DedicatedServer %s. No OriginalHost assigned and no way to find a matching one", host.Hostname)
		}

		res, err := c.DeviceConvertFromNetbox(oh)
		if err != nil {
			return fmt.Errorf("cannot convert to DedicatedServer")
		}

		host.Meta.ID = res.Meta.ID
		host.Meta.NetboxEntity = res.Meta.NetboxEntity
	}

	data := new(models.WritableDeviceWithConfigContext)

	//Go through every item and check if it must be updated

	params := dcim.NewDcimDevicesUpdateParams()
	params.WithID(host.Meta.ID).WithData(data)

	//TODO: Iterate over Inventory Items

	return nil
}

//HypervisorFindByHostname is like NetboxDeviceGet but checks if the device has a cluster assigned
func (c Client) HypervisorFindByHostname(hostname string) (*types.DedicatedServer, error) {
	res, err := c.DeviceFind(hostname)

	if err != nil {
		return nil, err
	}

	d := res.Meta.NetboxEntity.(models.DeviceWithConfigContext)
	if d.Cluster == nil {
		return nil, fmt.Errorf("device %s not assigned to a Virtualization Cluster", hostname)
	}

	return res, nil
}
