package ipam

import "github.com/hosting-de-labs/go-netbox/netbox/client"

type Client struct {
	client client.NetBoxAPI
}

func NewClient(netboxClient client.NetBoxAPI) Client {
	return Client{
		client: netboxClient,
	}
}
