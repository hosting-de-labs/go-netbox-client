package dcim

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	netboxIpam "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"

	"github.com/go-openapi/strfmt"

	"github.com/hosting-de-labs/go-netbox-client/netbox/utils"
	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func (c Client) DeviceConvertFromNetbox(device interface{}) (out *types.DedicatedServer, err error) {
	out = types.NewDedicatedServer()

	primaryIPv4 := &models.NestedIPAddress{}
	primaryIPv6 := &models.NestedIPAddress{}
	switch device.(type) {
	case models.Device:
		d := device.(models.Device)
		out.SetNetboxEntity(d.ID, device)
		err = out.SetCustomFields(d.CustomFields)
		if err != nil {
			return nil, err
		}

		if d.Name != nil {
			out.Hostname = *d.Name
		}

		if d.AssetTag != nil {
			out.AssetTag = *d.AssetTag
		}

		out.SerialNumber = d.Serial
		out.Tags = d.Tags
		out.Comments = strings.Split(d.Comments, "\n") //TODO: use utils.ParseVMComment

		primaryIPv4 = d.PrimaryIp4
		primaryIPv6 = d.PrimaryIp6
	case models.DeviceWithConfigContext:
		d := device.(models.DeviceWithConfigContext)
		out.SetNetboxEntity(d.ID, device)
		err = out.SetCustomFields(d.CustomFields)
		if err != nil {
			return nil, err
		}

		if d.Name != nil {
			out.Hostname = *d.Name
		}

		if d.AssetTag != nil {
			out.AssetTag = *d.AssetTag
		}

		out.SerialNumber = d.Serial
		out.Tags = d.Tags
		out.Comments = strings.Split(d.Comments, "\n") //TODO: use utils.ParseVMComment

		primaryIPv4 = d.PrimaryIp4
		primaryIPv6 = d.PrimaryIp6
	default:
		return nil, fmt.Errorf("unsupported type for device: %s", reflect.TypeOf(device))
	}

	for _, tag := range out.Tags {
		if tag == "managed" {
			out.IsManaged = true
			break
		}
	}

	if primaryIPv4 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*primaryIPv4.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv4 = &types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv4,
		}
	}

	if primaryIPv6 != nil {
		//split cidr
		address, cidr, err := utils.SplitCidrFromIP(*primaryIPv6.Address)
		if err != nil {
			return nil, err
		}

		out.PrimaryIPv6 = &types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv6,
		}
	}

	out.SetOriginalEntity(out.Copy())

	return out, nil
}

func (c Client) DeviceConvertToNetbox(server types.DedicatedServer) (out *models.WritableDeviceWithConfigContext, intf []*models.WritableDeviceInterface, err error) {
	out = &models.WritableDeviceWithConfigContext{
		Name:        &server.Hostname,
		Tags:        server.Tags,
		AssetTag:    &server.AssetTag,
		Serial:      server.SerialNumber,
		LastUpdated: strfmt.DateTime(time.Now()),
	}

	if server.Created != nil {
		out.Created = *server.Created
	}

	if !server.HasNetboxEntity() {
		err = c.PopulateDevice(&server)
		if err != nil {
			return nil, nil, err
		}
	}

	out.ID = server.Meta.ID

	//Interfaces
	for _, objIntf := range server.NetworkInterfaces {
		netIntf, err := c.InterfaceConvertToNetbox(out.ID, objIntf)
		if err != nil {
			return nil, nil, err
		}

		intf = append(intf, netIntf)
	}

	//Primary IPs
	ipamClient := netboxIpam.NewClient(c.client)
	if server.PrimaryIPv4 != nil {
		res, err := ipamClient.IPAddressFind(*server.PrimaryIPv4)
		if err != nil {
			return nil, nil, err
		}

		out.PrimaryIp4 = &res.ID
	}

	if server.PrimaryIPv6 != nil {
		res, err := ipamClient.IPAddressFind(*server.PrimaryIPv6)
		if err != nil {
			return nil, nil, err
		}

		out.PrimaryIp4 = &res.ID
	}

	return out, intf, nil
}
