package netbox_api

import (
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_api/dcim"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/netbox_api/ipam"
	"github.com/jarcoal/httpmock"
)

func RunServer() {
	//dcim
	httpmock.RegisterResponder(mock_dcim.DeviceGetResponder())
	httpmock.RegisterResponder(mock_dcim.InterfaceGetResponder())

	//ipam
	httpmock.RegisterResponder(mock_ipam.IpAddressFindResponder())

	httpmock.Activate()
}
