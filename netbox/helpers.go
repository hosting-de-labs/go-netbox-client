package netbox

import (
	"fmt"

	"internal.keenlogics.com/di/netbox-sync/types"
)

func generateVMComment(host *types.VirtualServer) string {
	comment := "--- NETBOX SYNC: DO NOT MODIFY ---\n"
	//TODO: Add lines for additional disks
	if len(host.Resources.Disks) > 1 {
		comment += "\n"
		comment += "Additional disks:\n"

		for index, disk := range host.Resources.Disks {
			if index == 0 {
				continue
			}

			comment += fmt.Sprintf("Size: %d MBytes", disk.Size)
		}
	}

	if len(host.Comments) > 0 {
		comment += "\n"
		for _, line := range host.Comments {
			comment += "\n"
			comment += line
			comment += "\n"
		}
	}

	comment += "--- NETBOX SYNC: DO NOT MODIFY ---"

	return comment
}
