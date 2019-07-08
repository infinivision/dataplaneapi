package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/dataplaneapi/client"
	"github.com/haproxytech/dataplaneapi/client/bind"
	frontend "github.com/haproxytech/dataplaneapi/client/frontend"
	"github.com/haproxytech/models"
)

//https://github.com/haproxy/haproxy/blob/master/doc/configuration.txt, 2.4. Time format
type DbNode struct {
	Name string
	Host string
	Port int
}

type DbCluster struct {
	Name     string
	BindPort int
	Nodes    []DbNode
}

func EnsureDbCluster(rt *runtimeClient.Runtime, authInfo runtime.ClientAuthInfoWriter, dbCluster *DbCluster) (err error) {
	var frtRsp *frontend.GetFrontendOK
	var bndsRsp *bind.GetBindsOK
	frontendName := dbCluster.Name + "-front"
	backendName := dbCluster.Name + "-back"
	tcpTimeout := int64(600 * 1000) //10 minutes
	maxConn := int64(3000)
	bindPort := int64(dbCluster.BindPort)
	frtModExp := &models.Frontend{
		ClientTimeout:  &tcpTimeout,
		DefaultBackend: backendName,
		Maxconn:        &maxConn,
		Mode:           "tcp",
		Name:           frontendName,
	}
	bndsModExp := models.Binds{
		&models.Bind{
			Address: "*",
			Name:    "bind1",
			Port:    &bindPort,
		},
	}
	/*
		bkdModExp := &models.Backend{
			Balance: models.Balance{
				Algorithm: "roundrobin",
			},
			Mode:          "tcp",
			Name:          backendName,
			ServerTimeout: &tcpTimeout,
		}
	*/

	srvsModExp := make([]models.Server, 0)
	for _, node := range dbCluster.Nodes {
		port := int64(node.Port)
		srv := models.Server{
			Name:         node.Name,
			Address:      node.Host,
			Port:         &port,
			OnMarkedDown: "shutdown-sessions",
		}
		srvsModExp = append(srvsModExp, srv)
	}
	sort.Slice(srvsModExp, func(i, j int) bool {
		return srvsModExp[i].Address < srvsModExp[j].Address || (srvsModExp[i].Address == srvsModExp[j].Address && *srvsModExp[i].Port < *srvsModExp[j].Port)
	})

	var frtMod *models.Frontend
	var bndsMod models.Binds
	frtClient := frontend.New(rt, strfmt.Default)
	if frtRsp, err = frtClient.GetFrontend(frontend.NewGetFrontendParams().WithName(frontendName), authInfo); err != nil {
		if _, ok := err.(*frontend.GetFrontendNotFound); ok {
			err = nil
		} else {
			err = errors.Wrap(err, "")
			return
		}
	} else {
		frtMod = frtRsp.Payload.Data
	}
	bndClient := bind.New(rt, strfmt.Default)
	if bndsRsp, err = bndClient.GetBinds(bind.NewGetBindsParams().WithFrontend(frontendName), authInfo); err != nil {
		if _, ok := err.(*bind.GetBindsDefault); ok {
			err = nil
		} else {
			err = errors.Wrap(err, "")
			return
		}
	} else {
		bndsMod = bndsRsp.Payload.Data
	}

	/*
		bkdClient := backend.New(rt, strfmt.Default)
		bkdRsp, err := bkdClient.GetBackend(frontend.NewGetFrontendParams().WithName(backendName), authInfo)
		if err != nil {
			err = errors.Wrap(err, "")
			return
		}
		bkdMod := bkdRsp.Payload.Data

		srvClient := server.New(rt, strfmt.Default)
		srvsRsp, err := srvClient.GetServers(server.NewGetServersParams().WithBackend(backendName), authInfo)
		if err != nil {
			err = errors.Wrap(err, "")
			return
		}
		srvsMod := srvsRsp.Payload.Data
	*/

	if !reflect.DeepEqual(frtModExp, frtMod) {
		if frtMod == nil {
			if _, _, err = frtClient.CreateFrontend(frontend.NewCreateFrontendParams().WithData(frtModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		} else {
			if _, _, err = frtClient.ReplaceFrontend(frontend.NewReplaceFrontendParams().WithData(frtModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		}
	}
	// binds
	if !reflect.DeepEqual(bndsModExp, bndsMod) {
		// TODO: delete all binds
		if frtMod == nil {
			if _, _, err = bndClient.CreateBind(bind.NewCreateBindParams().WithData(bndsModExp[0]), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		} else {
			if _, _, err = bndClient.ReplaceBind(bind.NewReplaceBindParams().WithData(bndsModExp[0]).WithFrontend(frontendName).WithName("bind1"), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		}
	}
	// backend
	// servers

	return
}

func main() {
	// create the transport
	rt := runtimeClient.New(os.Getenv("HAPROXY_DATAPLANE_ADDR"), client.DefaultBasePath, []string{"http"})
	authInfo := runtimeClient.BasicAuth("admin", "mypassword")

	// create the API client, with the transport
	client := frontend.New(rt, strfmt.Default)

	// make the request to get all frontends
	resp, err := client.GetFrontends(nil, authInfo)
	if err != nil {
		log.Fatalf("got error %+v", err)
	}
	fmt.Printf("response: %#v\n", resp.Payload)

	dbCluster := &DbCluster{
		Name:     "pg",
		BindPort: 12345,
		Nodes: []DbNode{
			DbNode{
				Name: "keeper0",
				Host: "127.0.0.1",
				Port: 10000,
			},
			DbNode{
				Name: "keeper1",
				Host: "127.0.0.1",
				Port: 10001,
			},
			DbNode{
				Name: "keeper2",
				Host: "127.0.0.1",
				Port: 10002,
			},
		},
	}
	if err = EnsureDbCluster(rt, authInfo, dbCluster); err != nil {
		log.Fatalf("got error %+v", err)
	}
}
