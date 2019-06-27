package dcim

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"

	"github.com/go-openapi/swag"

	"github.com/hosting-de-labs/go-netbox/netbox/client/dcim"
	"github.com/hosting-de-labs/go-netbox/netbox/models"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
)

func (c Client) InterfaceGet(interfaceID int64) (*models.DeviceInterface, error) {
	params := dcim.NewDcimInterfacesReadParams()
	params.WithID(interfaceID)

	res, err := c.client.Dcim.DcimInterfacesRead(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

//InterfaceGet retrieves an existing device interface object.
func (c Client) InterfaceFind(deviceID int64, interfaceName string) (*models.DeviceInterface, error) {
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

	return res.Payload.Results[0], nil
}

func (c Client) InterfaceCreate(deviceID int64, networkInterface *types.HostNetworkInterface) (*models.DeviceInterface, error) {
	data := &models.WritableDeviceInterface{}
	data.Device = &deviceID
	data.Name = &networkInterface.Name

	d, err := c.DeviceGet(deviceID)
	if err != nil {
		return nil, err
	}

	ipamClient := netboxIpam.NewClient(c.client)

	if networkInterface.UntaggedVlan != nil && len(networkInterface.TaggedVlans) > 0 {
		//Set mode to tagged
		data.Mode = swag.Int64(200)

		netboxVlan, err := ipamClient.VLANGet(networkInterface.UntaggedVlan.ID, &d.Site.ID)
		if err != nil {
			return nil, err
		}

		data.UntaggedVlan = netboxVlan.Vid

		var taggedVlans []int64
		for _, vlan := range networkInterface.TaggedVlans {
			netboxVlan, err := ipamClient.VLANGet(vlan.ID, &d.Site.ID)
			if err != nil {
				return nil, err
			}

			taggedVlans = append(taggedVlans, *netboxVlan.Vid)
		}

		data.TaggedVlans = taggedVlans
	} else if networkInterface.UntaggedVlan != nil {
		//Set mode to access
		data.Mode = swag.Int64(100)

		netboxVlan, err := ipamClient.VLANGet(networkInterface.UntaggedVlan.ID, &d.Site.ID)
		if err != nil {
			return nil, err
		}

		data.UntaggedVlan = netboxVlan.Vid
	} else if len(networkInterface.TaggedVlans) > 0 {
		//Set mode to tagged all
		data.Mode = swag.Int64(300)

		var taggedVlans []int64
		for _, vlan := range networkInterface.TaggedVlans {
			netboxVlan, err := ipamClient.VLANGet(vlan.ID, &d.Site.ID)
			if err != nil {
				return nil, err
			}

			taggedVlans = append(taggedVlans, *netboxVlan.Vid)
		}

		data.TaggedVlans = taggedVlans
	} else {
		data.Mode = nil
	}

	if networkInterface.FormFactor > 0 {
		data.FormFactor = int64(networkInterface.FormFactor)
	}

	params := dcim.NewDcimInterfacesCreateParams()
	params.WithData(data)

	res, err := c.client.Dcim.DcimInterfacesCreate(params, nil)
	if err != nil {
		return nil, err
	}

	return c.InterfaceGet(res.Payload.ID)
}

//InterfaceGetCreate is a convenience method to retrieve an existing device interface or otherwise to create it.
func (c Client) InterfaceGetCreate(deviceID int64, networkInterface *types.HostNetworkInterface) (*models.DeviceInterface, error) {
	res, err := c.InterfaceFind(deviceID, networkInterface.Name)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return c.InterfaceCreate(deviceID, networkInterface)
	}

	return res, nil
}
