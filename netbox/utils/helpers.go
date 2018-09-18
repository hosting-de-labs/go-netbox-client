package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hosting-de-labs/go-netbox-client/types"
)

func SplitCidrFromIP(ipWithCidr string) (string, uint16, error) {
	res := strings.Split(ipWithCidr, "/")

	cidrInt, err := strconv.ParseInt(res[1], 10, 16)
	if err != nil {
		return "", 0, err
	}

	return res[0], uint16(cidrInt), nil
}

func ConvertCustomFields(customFields interface{}) map[string]string {
	tmp := customFields.(map[string]interface{})

	out := make(map[string]string)
	for key, val := range tmp {
		if val != nil {
			tmpVal, ok := val.(string)
			if ok {
				out[key] = tmpVal
				continue
			}

			//TODO: parse maps
		}
	}

	return out
}

//GenerateSlug generates a netbox compatible identifier for eg. manufacturers, devices
func GenerateSlug(s string) string {
	//to lower case
	slug := strings.ToLower(s)

	//replace spaces with hyphens
	slug = strings.Replace(slug, " ", "-", -1)

	//remove forward-slashes
	slug = strings.Replace(slug, "/", "", -1)

	return slug
}

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
