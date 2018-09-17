package virtualization

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//InterfaceGetAll returns all interfaces of a virtual machine identified by it's id
func (c Client) InterfaceGetAll(vmID int64) ([]*models.Interface, error) {
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

//InterfaceGet retrieves an existing VM interface object.
func (c Client) InterfaceGet(vm *types.VirtualServer, interfaceName string) (*models.Interface, error) {
	params := virtualization.NewVirtualizationInterfacesListParams()
	params.Name = swag.String(interfaceName)
	params.VirtualMachineID = &vm.ID

	res, err := c.client.Virtualization.VirtualizationInterfacesList(params, nil)
	if err != nil {
		return nil, err
	}

	if *res.Payload.Count > 1 {
		return nil, fmt.Errorf("Interface %s is not unique", interfaceName)
	}

	if *res.Payload.Count == 0 {
		return nil, nil
	}

	return res.Payload.Results[0], nil
}

//InterfaceCreate creates a VM interface in Netbox.
func (c Client) InterfaceCreate(vm *types.VirtualServer, interfaceName string, vlan *models.VLAN) (*models.Interface, error) {
	data := new(models.WritableVirtualizationInterface)
	data.VirtualMachine = &vm.ID
	data.Tags = []string{}
	data.Name = &interfaceName
	// data.Mode = &models.WritableVirtualizationInterfaceMode{Value: swag.Int64(100)}

	if vlan != nil {
		data.UntaggedVlan = vlan.ID
	}

	data.TaggedVlans = []int64{}

	params := virtualization.NewVirtualizationInterfacesCreateParams()
	params.WithData(data)

	_, err := c.client.Virtualization.VirtualizationInterfacesCreate(params, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create interface. Original error was %s", err)
	}

	return c.InterfaceGet(vm, interfaceName)
}

//InterfaceGetCreate is a convenience method to retrieve an existing VM interface or otherwise to create it.
func (c Client) InterfaceGetCreate(vm *types.VirtualServer, interfaceName string, vlan *models.VLAN) (*models.Interface, error) {
	res, err := c.InterfaceGet(vm, interfaceName)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return c.InterfaceCreate(vm, interfaceName, vlan)
	}

	return res, nil
}
