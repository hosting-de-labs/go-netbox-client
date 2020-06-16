package virtualization

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//InterfaceGet fetches an interface from netbox
func (c Client) InterfaceGet(interfaceID int64) (out *types.NetworkInterface, err error) {
	params := virtualization.NewVirtualizationInterfacesReadParams()
	params.WithID(interfaceID)

	res, err := c.client.Virtualization.VirtualizationInterfacesRead(params, nil)
	if err != nil {
		return nil, err
	}

	return c.InterfaceConvertFromNetbox(*res.Payload)
}

//InterfaceFindAll returns all interfaces of a virtual machine identified by it's id
func (c Client) InterfaceFindAll(vmID int64) ([]*models.VirtualMachineInterface, error) {
	params := virtualization.NewVirtualizationInterfacesListParams()
	params.VirtualMachineID = swag.Int64(vmID)

	res, err := c.client.Virtualization.VirtualizationInterfacesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	return res.Payload.Results, nil
}

//InterfaceFind retrieves an existing VM interface object.
func (c Client) InterfaceFind(vmID int64, interfaceName string) (out *types.NetworkInterface, err error) {
	params := virtualization.NewVirtualizationInterfacesListParams()
	params.Name = swag.String(interfaceName)

	params.VirtualMachineID = &vmID

	res, err := c.client.Virtualization.VirtualizationInterfacesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("interface %s is not unique", interfaceName)
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	out, err = c.InterfaceConvertFromNetbox(*res.Payload.Results[0])
	if err != nil {
		return nil, err
	}

	return out, nil
}

//InterfaceCreate creates a VM interface in Netbox.
func (c Client) InterfaceCreate(vmID int64, intf types.NetworkInterface) (machineInterface *types.NetworkInterface, err error) {
	data, err := c.InterfaceConvertToNetbox(vmID, intf)

	if err != nil {
		return nil, err
	}

	params := virtualization.NewVirtualizationInterfacesCreateParams()
	params.WithData(data)

	_, err = c.client.Virtualization.VirtualizationInterfacesCreate(params, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create interface: %s", err)
	}

	return c.InterfaceFind(vmID, intf.Name)
}

//InterfaceGetCreate is a convenience method to retrieve an existing VM interface or otherwise to create it.
func (c Client) InterfaceGetCreate(vmID int64, intf types.NetworkInterface) (machineInterface *types.NetworkInterface, err error) {
	res, err := c.InterfaceFind(vmID, intf.Name)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return c.InterfaceCreate(vmID, intf)
	}

	return res, nil
}
