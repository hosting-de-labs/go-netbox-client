package netbox

import (
	"fmt"

	"github.com/hosting-de-labs/go-netbox-client/types"
)

func generateVMComment(host *types.VirtualServer) string {
	comment := "--- NETBOX SYNC: DO NOT MODIFY ---"

	//regular comments
	if len(host.Comments) > 0 {
		comment += "\nComments:"
		for _, line := range host.Comments {
			comment += "\n" + line
		}
	}

	//additional disks
	if len(host.Resources.Disks) > 1 {
		comment += "\nAdditional disks:"

		for index, disk := range host.Resources.Disks {
			if index == 0 {
				continue
			}

			comment += fmt.Sprintf("\nSize: %d MBytes", disk.Size)
		}
	}

	comment += "\n--- NETBOX SYNC: DO NOT MODIFY ---"

	return comment
}
