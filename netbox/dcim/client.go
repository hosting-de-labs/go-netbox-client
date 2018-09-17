package dcim

import "github.com/hosting-de-labs/go-netbox/netbox/client"

type Client struct {
	netboxClient *client.NetBox
}

func (c *Client) getClient() *client.NetBox {
	return c.netboxClient
}

func (c *Client) setClient(client *client.NetBox) {
	c.netboxClient = client
}
