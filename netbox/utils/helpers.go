package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/hosting-de-labs/go-netbox-client/types"
)

const (
	//NetboxSyncVMCommentStartingEndingLine
	NetboxSyncVMCommentStartingEndingLine string = "--- NETBOX SYNC: DO NOT MODIFY ---"
)

func SplitCidrFromIP(ipWithCidr string) (string, uint16, error) {
	res := strings.SplitN(ipWithCidr, "/", 2)

	if len(res) != 2 {
		return "", 0, fmt.Errorf("invalid input: %q", ipWithCidr)
	}

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

func GenerateVMComment(host types.VirtualServer) string {
	comment := NetboxSyncVMCommentStartingEndingLine

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

		for i := 1; i < len(host.Resources.Disks); i++ {
			comment += fmt.Sprintf("\nSize: %d MBytes", host.Resources.Disks[i].Size)
		}
	}

	comment += "\n"
	comment += NetboxSyncVMCommentStartingEndingLine

	return comment
}

func ParseVMComment(comments string, host *types.VirtualServer) {
	lines := strings.Split(comments, "\n")
	for i := 0; i < len(lines); i++ {
		//Parse Comments
		if lines[i] == "Comments:" {
			for i = i + 1; i < len(lines); i++ {
				if lines[i] == "" || lines[i] == NetboxSyncVMCommentStartingEndingLine || strings.HasSuffix(lines[i], ":") {
					break
				}

				host.Comments = append(host.Comments, lines[i])
			}
		}

		//Parse additional disks
		if lines[i] == "Additional disks:" {
			for i = i + 1; i < len(lines); i++ {
				if lines[i] == "" || lines[i] == NetboxSyncVMCommentStartingEndingLine || strings.HasSuffix(lines[i], ":") {
					break
				}

				//@TODO: Implement
				r := regexp.MustCompilePOSIX(`^Size: ([0-9]+) MBytes$`)
				res := r.FindStringSubmatch(lines[i])

				if len(res) != 2 {
					//FIXME: error handling?
					continue
				}

				size, err := strconv.ParseInt(res[1], 10, 64)
				if err != nil {
					//FIXME: error handling?
					continue
				}

				host.Resources.Disks = append(host.Resources.Disks, types.VirtualServerDisk{Size: size})
			}
		}
	}
}

func ExtractFromApiError(err error) error {
	apiError := err.(interface{}).(*runtime.APIError)
	body := apiError.Response.(interface{}).(runtime.ClientResponse)

	return fmt.Errorf("code: %d, operation: %s, ResponseBodyType: %s, ResponseBody: %+v", apiError.Code, apiError.OperationName, body, reflect.TypeOf(body.(interface{})).String())
}
