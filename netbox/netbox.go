package netbox

import (
	"time"

	"github.com/go-openapi/strfmt"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxClient "github.com/hosting-de-labs/go-netbox/netbox/client"

	dcimClient "github.com/hosting-de-labs/go-netbox-client/netbox/dcim"
	ipamClient "github.com/hosting-de-labs/go-netbox-client/netbox/ipam"
	virtualizationClient "github.com/hosting-de-labs/go-netbox-client/netbox/virtualization"
)

type NetBox struct {
	NetboxClient netboxClient.NetBox

	DCIM           dcimClient.Client
	IPAM           ipamClient.Client
	Virtualization virtualizationClient.Client
}

func NewAPIClient(url string, token string, defaultTimeout time.Duration) NetBox {
	//@TODO: timeout
	// timeout := 10 * time.Second
	// if defaultTimeout > 0 {
	// 	timeout = defaultTimeout
	// }

	t := runtimeclient.New(url, "/api", []string{"https"})
	t.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", "Token "+token)

	c := netboxClient.New(t, strfmt.Default)

	return NetBox{
		NetboxClient:   *c,
		DCIM:           dcimClient.NewClient(*c),
		IPAM:           ipamClient.NewClient(*c),
		Virtualization: virtualizationClient.NewClient(*c),
	}
}
