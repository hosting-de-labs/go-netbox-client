package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
)

func (c Client) InterfaceCreate(deviceID int64, networkInterface types.NetworkInterface) (*types.NetworkInterface, error) {
	data, err := c.InterfaceConvertToNetbox(deviceID, networkInterface)
	if err != nil {
		return nil, err
	}

	params := dcim.NewDcimInterfacesCreateParams()
	params.WithData(data)

	res, err := c.client.Dcim.DcimInterfacesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.InterfaceGet(res.Payload.ID)
}

//InterfaceGet retrieves an existing device interface object.
func (c Client) InterfaceFind(deviceID int64, interfaceName string) (*types.NetworkInterface, error) {
	params := dcim.NewDcimInterfacesListParams()
	params.Name = &interfaceName

	params.DeviceID = &deviceID

	res, err := c.client.Dcim.DcimInterfacesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("Interface %s is not unique", interfaceName)
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	out, err := c.InterfaceConvertFromNetbox(*res.Payload.Results[0])
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c Client) InterfaceFindAll(deviceID int64) (out []types.NetworkInterface, err error) {
	params := dcim.NewDcimInterfacesListParams()
	params.WithDeviceID(&deviceID)

	res, err := c.client.Dcim.DcimInterfacesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	for _, nbInterface := range res.Payload.Results {
		intf, err := c.InterfaceConvertFromNetbox(*nbInterface)
		if err != nil {
			return nil, err
		}

		out = append(out, *intf)
	}

	return out, nil
}

func (c Client) InterfaceGet(interfaceID int64) (*types.NetworkInterface, error) {
	params := dcim.NewDcimInterfacesReadParams()
	params.WithID(interfaceID)

	res, err := c.client.Dcim.DcimInterfacesRead(params, nil)
	if err != nil {
		return nil, err
	}

	out, err := c.InterfaceConvertFromNetbox(*res.Payload)
	if err != nil {
		return nil, err
	}

	return out, nil
}

//InterfaceGetCreate is a convenience method to retrieve an existing device interface or otherwise to create it.
func (c Client) InterfaceGetCreate(deviceID int64, networkInterface types.NetworkInterface) (*types.NetworkInterface, error) {
	res, err := c.InterfaceFind(deviceID, networkInterface.Name)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return c.InterfaceCreate(deviceID, networkInterface)
	}

	return res, nil
}
