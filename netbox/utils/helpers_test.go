package utils

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/stretchr/testify/assert"
)

func TestSplitCidrFromIP(t *testing.T) {
	cases := []struct {
		ipWithCIDR string
		resultIP   string
		resultCIDR uint16
		isError    bool
	}{
		{
			ipWithCIDR: "192.168.10.70/24",
			resultIP:   "192.168.10.70",
			resultCIDR: 24,
			isError:    false,
		},
		{
			ipWithCIDR: "192.168.10.70/cc",
			resultIP:   "",
			resultCIDR: 0,
			isError:    true,
		},
		{
			ipWithCIDR: "192.168.10.70",
			resultIP:   "",
			resultCIDR: 0,
			isError:    true,
		},
		{
			ipWithCIDR: "random-string-that-is-never-an-ipaddress",
			resultIP:   "",
			resultCIDR: 0,
			isError:    true,
		},
		{
			ipWithCIDR: "random-string-that-is-never-an-ipaddress/random-invalid-cidr",
			resultIP:   "",
			resultCIDR: 0,
			isError:    true,
		},
	}

	for _, testcase := range cases {
		ip, cidr, err := SplitCidrFromIP(testcase.ipWithCIDR)

		if testcase.isError {
			assert.Error(t, err, "error message: %q")
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, testcase.resultIP, ip)
		assert.Equal(t, testcase.resultCIDR, cidr)
	}
}

func TestGenerateSlug(t *testing.T) {
	cases := []struct {
		input  string
		result string
	}{
		{
			input:  "Intel",
			result: "intel",
		},
		{
			input:  "AMD",
			result: "amd",
		},
		{
			input:  "Intel Xeon X5770",
			result: "intel-xeon-x5770",
		},
	}

	for _, testcase := range cases {
		assert.Equal(t, testcase.result, GenerateSlug(testcase.input))
	}
}

func TestGenerateVMComment(t *testing.T) {
	cases := []struct {
		host   types.VirtualServer
		result string
	}{
		{
			types.VirtualServer{
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
		{
			types.VirtualServer{
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
						types.VirtualServerDisk{
							Size: 10,
						},
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
Additional disks:
Size: 10 MBytes
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
		{
			types.VirtualServer{
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
						types.VirtualServerDisk{
							Size: 10,
						},
						types.VirtualServerDisk{
							Size: 30,
						},
						types.VirtualServerDisk{
							Size: 40,
						},
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
Additional disks:
Size: 10 MBytes
Size: 30 MBytes
Size: 40 MBytes
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
		{
			types.VirtualServer{
				Host: types.Host{
					Comments: []string{
						"Foo",
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
Comments:
Foo
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
		{
			types.VirtualServer{
				Host: types.Host{
					Comments: []string{
						"Foo",
					},
				},
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
Comments:
Foo
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
		{
			types.VirtualServer{
				Host: types.Host{
					Comments: []string{
						"Foo",
					},
				},
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
						types.VirtualServerDisk{
							Size: 10,
						},
					},
				},
			},
			`--- NETBOX SYNC: DO NOT MODIFY ---
Comments:
Foo
Additional disks:
Size: 10 MBytes
--- NETBOX SYNC: DO NOT MODIFY ---`,
		},
	}

	for _, testcase := range cases {
		assert.Equal(t, testcase.result, GenerateVMComment(&testcase.host))
	}
}

func TestParseVMComment(t *testing.T) {
	cases := []struct {
		comment string
		host    types.VirtualServer
	}{
		{
			`--- NETBOX SYNC: DO NOT MODIFY ---
Comments:
Foo
Additional disks:
Size: 10 MBytes
--- NETBOX SYNC: DO NOT MODIFY ---`,
			types.VirtualServer{
				Host: types.Host{
					Comments: []string{
						"Foo",
					},
				},
				Resources: types.VirtualServerResources{
					Disks: []types.VirtualServerDisk{
						types.VirtualServerDisk{
							Size: 20,
						},
						types.VirtualServerDisk{
							Size: 10,
						},
					},
				},
			},
		},
	}

	for _, testcase := range cases {
		vm := types.VirtualServer{}
		vm.Resources.Disks = append(vm.Resources.Disks, types.VirtualServerDisk{Size: 20})

		ParseVMComment(testcase.comment, &vm)

		assert.Equal(t, testcase.host, vm)
	}
}
