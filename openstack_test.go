package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/stretchr/testify/require"
)

func TestListCluster(t *testing.T) {
	fmt.Println("Starting Testing")
	opts, err := openstack.AuthOptionsFromEnv()
	fmt.Println("Identity: ", opts.IdentityEndpoint)

	fmt.Println("Username ", opts.Username)
	fmt.Println("Domain ", opts.DomainID)
	provider, err := openstack.AuthenticatedClient(opts)
	require.NoError(t, err)

	if err != nil {
		panic(err)
	}

	fmt.Println("End auth")
	listClust := clusters.ListOpts{}
	client, err := openstack.NewCCE(provider, golangsdk.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})

	allClusts, err := clusters.List(client, listClust)
	require.NoError(t, err)
	for _, node := range allClusts {
		fmt.Printf("%+v\n", node.Status.Phase)
	}
}
