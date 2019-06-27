package mock

import (
	"github.com/hosting-de-labs/go-netbox-client/test/mock/dcim"
	"github.com/hosting-de-labs/go-netbox-client/test/mock/ipam"
	"github.com/jarcoal/httpmock"
)

func RunServer() {
	//dcim
	httpmock.RegisterResponder(mock_dcim.InterfaceGetResponder())

	//ipam
	httpmock.RegisterResponder(mock_ipam.IpAddressFindResponder())

	httpmock.Activate()
}
