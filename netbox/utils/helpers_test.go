package utils

import (
	"testing"

	"github.com/hosting-de-labs/go-netbox-client/types"
	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(t, testcase.result, generateVMComment(&testcase.host))
	}
}
