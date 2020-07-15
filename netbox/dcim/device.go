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
	if !host.HasNetboxEntity() {
		err := c.PopulateDevice(host)

		if err != nil {
			return err
		}
	}

	if !host.IsChanged() {
		return nil
	}

	data := new(models.WritableDeviceWithConfigContext)

	params := dcim.NewDcimDevicesUpdateParams()
	params.WithID(host.Meta.ID).WithData(data)

	//TODO: Go through every item and check if it must be updated

	//TODO: Iterate over Inventory Items

	return nil
}

func (c Client) PopulateDevice(device *types.DedicatedServer) (err error) {
	if device.HasNetboxEntity() {
		return nil
	}

	oh, err := c.DeviceFind(device.Hostname)
	if err != nil {
		return fmt.Errorf("cannot update DedicatedServer %s.: %s", device.Hostname, err)
	}

	device.SetNetboxEntity(oh.Meta.ID, oh.Meta.NetboxEntity)

	return nil
}
