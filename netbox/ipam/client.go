package ipam

import "github.com/hosting-de-labs/go-netbox/netbox/client"

type Client struct {
	client client.NetBox
}

func NewClient(netboxClient client.NetBox) Client {
	return Client{
		client: netboxClient,
	}
}
