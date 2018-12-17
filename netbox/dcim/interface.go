package dcim

import (
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
)

//InterfaceGet retrieves an existing device interface object.
func (c Client) InterfaceGet(interfaceName string, device *models.Device) (*models.Interface, error) {
	params := dcim.NewDcimInterfacesListParams()
	params.Name = &interfaceName

	params.DeviceID = &device.ID

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

	return res.Payload.Results[0], nil
}

//InterfaceCreate creates a device interface in Netbox.
func (c Client) InterfaceCreate(interfaceName string, device *models.Device, vlanTag *int64, interfaceFormFactor *int64) (*models.Interface, error) {
	data := new(models.WritableInterface)
	data.Device = &device.ID
	data.Name = &interfaceName
	data.Mode = swag.Int64(100)
	data.TaggedVlans = []int64{}

	ipamClient := netboxIpam.NewClient(c.client)

	if vlanTag != nil {
		vlan, err := ipamClient.VlanGet(*vlanTag, device.Site.ID)
		if err != nil {
			return nil, err
		}

		data.UntaggedVlan = &vlan.ID
	}

	if interfaceFormFactor != nil {
		data.FormFactor = *interfaceFormFactor
	}

	params := dcim.NewDcimInterfacesCreateParams()
	params.WithData(data)

	_, err := c.client.Dcim.DcimInterfacesCreate(params, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create interface with vlan Tag %d. Original error was %s", vlanTag, err)
	}

	return c.InterfaceGet(interfaceName, device)
}

//InterfaceGetCreate is a convenience method to retrieve an existing device interface or otherwise to create it.
func (c Client) InterfaceGetCreate(interfaceName string, device *models.Device, vlanTag *int64) (*models.Interface, error) {
	res, err := c.InterfaceGet(interfaceName, device)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return c.InterfaceCreate(interfaceName, device, vlanTag, nil)
	}

	return res, nil
}
