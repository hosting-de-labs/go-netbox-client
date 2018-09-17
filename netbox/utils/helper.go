package utils

import (
	"strconv"
	"strings"
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
