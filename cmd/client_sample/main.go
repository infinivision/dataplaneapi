package main

import (
	"log"
	"os"
	"reflect"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/dataplaneapi/client"
	"github.com/haproxytech/dataplaneapi/client/backend"
	"github.com/haproxytech/dataplaneapi/client/bind"
	"github.com/haproxytech/dataplaneapi/client/frontend"
	"github.com/haproxytech/dataplaneapi/client/server"
	"github.com/haproxytech/dataplaneapi/client/transactions"
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
	var version int64
	var txnID string
	var txnClient *transactions.Client
	var frtsRsp *frontend.GetFrontendsOK
	var frtRsp *frontend.GetFrontendOK
	var bndsRsp *bind.GetBindsOK
	var bkdRsp *backend.GetBackendOK
	var srvsRsp *server.GetServersOK
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

	bkdModExp := &models.Backend{
		Balance: &models.Balance{
			Algorithm: "roundrobin",
		},
		Mode:          "tcp",
		Name:          backendName,
		ServerTimeout: &tcpTimeout,
	}

	srvsModExp := make([]*models.Server, 0)
	for _, node := range dbCluster.Nodes {
		port := int64(node.Port)
		srv := &models.Server{
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
	if frtsRsp, err = frtClient.GetFrontends(nil, authInfo); err != nil {
		err = errors.Wrap(err, "")
		return
	}
	version = frtsRsp.Payload.Version

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

	var bkdMod *models.Backend
	var srvsMod models.Servers
	bkdClient := backend.New(rt, strfmt.Default)
	if bkdRsp, err = bkdClient.GetBackend(backend.NewGetBackendParams().WithName(backendName), authInfo); err != nil {
		if _, ok := err.(*backend.GetBackendNotFound); ok {
			err = nil
		} else {
			err = errors.Wrap(err, "")
			return
		}
	} else {
		bkdMod = bkdRsp.Payload.Data
	}

	srvClient := server.New(rt, strfmt.Default)
	if srvsRsp, err = srvClient.GetServers(server.NewGetServersParams().WithBackend(backendName), authInfo); err != nil {
		if _, ok := err.(*server.GetServersDefault); ok {
			err = nil
		} else {
			err = errors.Wrap(err, "")
			return
		}
	} else {
		srvsMod = srvsRsp.Payload.Data
	}

	if !reflect.DeepEqual(frtModExp, frtMod) || !reflect.DeepEqual(bndsModExp, bndsMod) || !reflect.DeepEqual(bkdModExp, bkdMod) || !reflect.DeepEqual(srvsModExp, srvsMod) {
		txnClient = transactions.New(rt, strfmt.Default)
		var txnRsp *transactions.StartTransactionCreated
		if txnRsp, err = txnClient.StartTransaction(transactions.NewStartTransactionParams().WithVersion(version), authInfo); err != nil {
			err = errors.Wrap(err, "")
			return
		}
		txnID = txnRsp.Payload.ID
	}

	// ensure frontend
	if !reflect.DeepEqual(frtModExp, frtMod) {
		if frtMod == nil {
			if _, _, err = frtClient.CreateFrontend(frontend.NewCreateFrontendParams().WithTransactionID(&txnID).WithData(frtModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		} else {
			if _, _, err = frtClient.ReplaceFrontend(frontend.NewReplaceFrontendParams().WithTransactionID(&txnID).WithName(frontendName).WithData(frtModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		}
	}
	// ensure binds
	if !reflect.DeepEqual(bndsModExp, bndsMod) {
		if bndsMod != nil {
			for _, bndMod := range bndsMod {
				if _, _, err = bndClient.DeleteBind(bind.NewDeleteBindParams().WithTransactionID(&txnID).WithFrontend(frontendName).WithName(bndMod.Name), authInfo); err != nil {
					err = errors.Wrap(err, "")
					return
				}
			}
		}
		if _, _, err = bndClient.CreateBind(bind.NewCreateBindParams().WithTransactionID(&txnID).WithFrontend(frontendName).WithData(bndsModExp[0]), authInfo); err != nil {
			err = errors.Wrap(err, "")
			return
		}
	}
	// ensure backend
	if !reflect.DeepEqual(bkdModExp, bkdMod) {
		if bkdMod == nil {
			if _, _, err = bkdClient.CreateBackend(backend.NewCreateBackendParams().WithTransactionID(&txnID).WithData(bkdModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		} else {
			if _, _, err = bkdClient.ReplaceBackend(backend.NewReplaceBackendParams().WithTransactionID(&txnID).WithName(backendName).WithData(bkdModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		}
	}
	// ensure servers
	if !reflect.DeepEqual(srvsModExp, srvsMod) {
		if srvsMod != nil {
			for _, srvMod := range srvsMod {
				if _, _, err = srvClient.DeleteServer(server.NewDeleteServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithName(srvMod.Name), authInfo); err != nil {
					err = errors.Wrap(err, "")
					return
				}
			}
		}
		for _, srvModExp := range srvsModExp {
			if _, _, err = srvClient.CreateServer(server.NewCreateServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithData(srvModExp), authInfo); err != nil {
				err = errors.Wrap(err, "")
				return
			}
		}
	}

	if txnID != "" {
		if _, _, err = txnClient.CommitTransaction(transactions.NewCommitTransactionParams().WithID(txnID), authInfo); err != nil {
			err = errors.Wrap(err, "")
			return
		}
	}
	return
}

func main() {
	spew.Config = spew.ConfigState{ContinueOnMethod: true}
	// create the transport
	rt := runtimeClient.New(os.Getenv("HAPROXY_DATAPLANE_ADDR"), client.DefaultBasePath, []string{"http"})
	authInfo := runtimeClient.BasicAuth("admin", "mypassword")

	var err error
	/*
		// create the API client, with the transport
		cfgCli := configuration.New(rt, strfmt.Default)

		var cfgRsp *configuration.GetHAProxyConfigurationOK
		if cfgRsp, err = cfgCli.GetHAProxyConfiguration(nil, authInfo); err != nil {
			err = errors.Wrap(err, "")
			log.Fatalf("got error %s\n%+v", spew.Sdump(err), err)
		}
		fmt.Printf("haproxy configuration: %s\n", spew.Sdump(cfgRsp.Payload))
	*/
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
		log.Fatalf("got error %s\n%+v", spew.Sdump(err), err)
	}
}
