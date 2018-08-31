package dcim

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
)

//InterfaceGet retrieves an existing device interface object.
func InterfaceGet(netboxClient *client.NetBox, interfaceName string, device *models.Device) (*models.Interface, error) {
	params := dcim.NewDcimInterfacesListParams()
	params.Name = &interfaceName

	params.DeviceID = &device.ID

	res, err := netboxClient.Dcim.DcimInterfacesList(params, nil)
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
func InterfaceCreate(netboxClient *client.NetBox, interfaceName string, device *models.Device, vlanTag *int64, interfaceFormFactor *int64) (*models.Interface, error) {
	data := new(models.Interface)
	data.Device.ID = device.ID
	data.Name = &interfaceName
	data.Mode.Value = swag.Int64(100)
	data.TaggedVlans = []int64{}

	if vlanTag != nil {
		vlan, err := netboxIpam.VlanGet(netboxClient, *vlanTag, device.Site.ID)
		if err != nil {
			return nil, err
		}

		data.UntaggedVlan.ID = vlan.ID
	}

	if interfaceFormFactor != nil {
		data.FormFactor.Value = interfaceFormFactor
	}

	params := dcim.NewDcimInterfacesCreateParams()
	params.WithData(data)

	_, err := netboxClient.Dcim.DcimInterfacesCreate(params, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create interface with vlan Tag %d. Original error was %s", vlanTag, err)
	}

	return InterfaceGet(netboxClient, interfaceName, device)
}

//InterfaceGetCreate is a convenience method to retrieve an existing device interface or otherwise to create it.
func InterfaceGetCreate(netboxClient *client.NetBox, interfaceName string, device *models.Device, vlanTag *int64) (*models.Interface, error) {
	res, err := InterfaceGet(netboxClient, interfaceName, device)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return InterfaceCreate(netboxClient, interfaceName, device, vlanTag, nil)
	}

	return res, nil
}
