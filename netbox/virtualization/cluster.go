package virtualization

import (
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//ClusterGetByID retrieves a cluster by it's id
func (c Client) ClusterGetByID(clusterID int64) (*models.Cluster, error) {
	params := virtualization.NewVirtualizationClustersReadParams()
	params.SetID(clusterID)

	res, err := c.client.Virtualization.VirtualizationClustersRead(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}
