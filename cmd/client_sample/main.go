package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/dataplaneapi/client"
	frontend "github.com/haproxytech/dataplaneapi/client/frontend"
)

//https://github.com/haproxy/haproxy/blob/master/doc/configuration.txt, 2.4. Time format
type PostgresNode struct {
	Name string
	Host string
	Port int
}

type PostgresCluster struct {
	Name          string
	BindPort      int
	TimeoutClient int
	TimeoutServer int
	Servers       []PostgresNode
}

func EnsurePostgresCluster(pgCluster *PostgresCluster) (err error) {

	err = errors.Wrap(err, "")
	return
}

func main() {
	// create the transport
	rt := runtimeClient.New(os.Getenv("HAPROXY_DATAPLANE_ADDR"), client.DefaultBasePath, []string{"http"})
	writer := runtimeClient.BasicAuth("admin", "mypassword")

	// create the API client, with the transport
	client := frontend.New(rt, strfmt.Default)

	// make the request to get all frontends
	resp, err := client.GetFrontends(nil, writer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response: %#v\n", resp.Payload)
}
