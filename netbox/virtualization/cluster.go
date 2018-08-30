package virtualization

import (
	"github.com/hosting-de-labs/go-netbox/netbox/client"
	"github.com/hosting-de-labs/go-netbox/netbox/client/virtualization"
	"github.com/hosting-de-labs/go-netbox/netbox/models"
)

//ClusterGetByID retrieves a cluster by it's id
func ClusterGetByID(netboxClient *client.NetBox, clusterID int64) (*models.Cluster, error) {
	params := virtualization.NewVirtualizationClustersReadParams()
	params.SetID(clusterID)

	res, err := netboxClient.Virtualization.VirtualizationClustersRead(params, nil)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}
