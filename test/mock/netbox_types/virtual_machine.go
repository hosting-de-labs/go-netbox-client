package netbox_types

import (
	"github.com/go-openapi/swag"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

func MockNetboxVirtualMachine(addResources bool, addIPAddresses bool, addTags bool, addCustomFields bool) (out models.VirtualMachineWithConfigContext) {
	out.ID = 10
	out.Name = swag.String("VM1")

	if addResources {
		out.Vcpus = swag.Int64(1)
		out.Memory = swag.Int64(4096)
		out.Disk = swag.Int64(10240)
	}

	if addIPAddresses {
		//TODO: add interfaces when adding ip addresses
		out.PrimaryIp4 = &models.NestedIPAddress{Address: swag.String("127.0.0.1/32")}
		out.PrimaryIp6 = &models.NestedIPAddress{Address: swag.String("::1/128")}
	}

	if addTags {
		out.Tags = append(out.Tags, "Tag1")
		out.Tags = append(out.Tags, "managed")
	}

	if addCustomFields {
		customFields := make(map[string]interface{})
		customFields["hypervisor_label"] = "Hypervisor1"

		out.CustomFields = customFields
	}

	return out
}
