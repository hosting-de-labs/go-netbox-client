package dcim

import (
	"fmt"
	"reflect"
	"strings"

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

		out.Metadata.ID = d.ID
		out.Metadata.NetboxEntity = device
		out.Metadata.EntityType = reflect.TypeOf(device)

		out.Hostname = *d.Name
		out.Tags = d.Tags
		out.Comments = strings.Split(d.Comments, "\n") //TODO: use utils.ParseVMComment

		primaryIPv4 = d.PrimaryIp4
		primaryIPv6 = d.PrimaryIp6
	case models.DeviceWithConfigContext:
		d := device.(models.DeviceWithConfigContext)

		out.Metadata.ID = d.ID
		out.Metadata.NetboxEntity = device
		out.Metadata.EntityType = reflect.TypeOf(device)

		out.Hostname = *d.Name
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

		out.PrimaryIPv4 = types.IPAddress{
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

		out.PrimaryIPv6 = types.IPAddress{
			Address: address,
			CIDR:    cidr,
			Family:  types.IPAddressFamilyIPv6,
		}
	}

	return out, nil
}
