package netbox

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxClient "github.com/hosting-de-labs/go-netbox/netbox/client"
)

func NewAPIClient(url string, token string, defaultTimeout time.Duration) *netboxClient.NetBox {
	timeout := 10 * time.Second
	if defaultTimeout > 0 {
		timeout = defaultTimeout
	}

	t := runtimeclient.New(url, "/api", []string{"https"})
	t.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", "Token "+token)

	//TODO: default timeout setzen
	fmt.Println(timeout)

	c := netboxClient.New(t, strfmt.Default)

	return c
}
