package main

import (
	"os"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/dataplaneapi/client"
	"github.com/haproxytech/dataplaneapi/client/backend"
	"github.com/haproxytech/dataplaneapi/client/bind"
	"github.com/haproxytech/dataplaneapi/client/frontend"
	"github.com/haproxytech/dataplaneapi/client/server"
	"github.com/haproxytech/dataplaneapi/client/transactions"
	"github.com/haproxytech/models"
	log "github.com/sirupsen/logrus"
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
	var needCommitTxn bool
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
			Arguments: []string{},
		},
		Mode:          "tcp",
		Name:          backendName,
		ServerTimeout: &tcpTimeout,
	}

	srvsModExp := models.Servers{}
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

	txnClient = transactions.New(rt, strfmt.Default)
	var txnRsp *transactions.StartTransactionCreated
	if txnRsp, err = txnClient.StartTransaction(transactions.NewStartTransactionParams().WithVersion(version), authInfo); err != nil {
		err = errors.Wrap(err, "")
		return
	}
	txnID = txnRsp.Payload.ID

	if frtRsp, err = frtClient.GetFrontend(frontend.NewGetFrontendParams().WithTransactionID(&txnID).WithName(frontendName), authInfo); err != nil {
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
	if bndsRsp, err = bndClient.GetBinds(bind.NewGetBindsParams().WithTransactionID(&txnID).WithFrontend(frontendName), authInfo); err != nil {
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
	if bkdRsp, err = bkdClient.GetBackend(backend.NewGetBackendParams().WithTransactionID(&txnID).WithName(backendName), authInfo); err != nil {
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
	if srvsRsp, err = srvClient.GetServers(server.NewGetServersParams().WithTransactionID(&txnID).WithBackend(backendName), authInfo); err != nil {
		if _, ok := err.(*server.GetServersDefault); ok {
			err = nil
		} else {
			err = errors.Wrap(err, "")
			return
		}
	} else {
		srvsMod = srvsRsp.Payload.Data
	}

	// ensure frontend
	diff := cmp.Diff(frtMod, frtModExp)
	if diff != "" {
		needCommitTxn = true
		log.Infof("frontend difference: %s", diff)
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
	diff = cmp.Diff(bndsMod, bndsModExp)
	if diff != "" {
		needCommitTxn = true
		log.Infof("binds difference: %s", diff)
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
	diff = cmp.Diff(bkdMod, bkdModExp)
	if diff != "" {
		needCommitTxn = true
		log.Infof("backend difference: %s", diff)
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
	diff = cmp.Diff(srvsMod, srvsModExp)
	if diff != "" {
		needCommitTxn = true
		log.Infof("servers difference: %s", diff)
		for _, srvMod := range srvsMod {
			var srvModExp *models.Server
			var found bool
			for _, srvModExp = range srvsModExp {
				if srvModExp.Name == srvMod.Name {
					found = true
					break
				}
			}
			if found {
				diff = cmp.Diff(srvMod, srvModExp)
				if diff != "" {
					if _, _, err = srvClient.ReplaceServer(server.NewReplaceServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithName(srvMod.Name).WithData(srvModExp), authInfo); err != nil {
						err = errors.Wrap(err, "")
						return
					}
				}
			} else {
				if srvMod.Disabled {
					if _, _, err = srvClient.DeleteServer(server.NewDeleteServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithName(srvMod.Name), authInfo); err != nil {
						err = errors.Wrap(err, "")
						return
					}
				} else {
					srvMod.Disabled = true
					if _, _, err = srvClient.ReplaceServer(server.NewReplaceServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithName(srvMod.Name).WithData(srvMod), authInfo); err != nil {
						err = errors.Wrap(err, "")
						return
					}
				}
			}
		}
		for _, srvModExp := range srvsModExp {
			var found bool
			for _, srvMod := range srvsMod {
				if srvModExp.Name == srvMod.Name {
					found = true
					break
				}
			}
			if !found {
				if _, _, err = srvClient.CreateServer(server.NewCreateServerParams().WithTransactionID(&txnID).WithBackend(backendName).WithData(srvModExp), authInfo); err != nil {
					err = errors.Wrap(err, "")
					return
				}
			}
		}
	}

	if needCommitTxn {
		if _, _, err = txnClient.CommitTransaction(transactions.NewCommitTransactionParams().WithID(txnID), authInfo); err != nil {
			err = errors.Wrap(err, "")
			return
		}
		log.Infof("DbCluster %s didn't match the expectation, has been fixed in transaction %s(version=%d)", dbCluster.Name, txnID, version+1)
	} else {
		if _, err = txnClient.DeleteTransaction(transactions.NewDeleteTransactionParams().WithID(txnID), authInfo); err != nil {
			err = errors.Wrap(err, "")
			return
		}
		log.Infof("DbCluster %s matchs the expectation", dbCluster.Name)
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
				Name: "pg0",
				Host: "127.0.0.1",
				Port: 10000,
			},
			DbNode{
				Name: "pg1",
				Host: "127.0.0.1",
				Port: 10001,
			},
			DbNode{
				Name: "pg2",
				Host: "127.0.0.1",
				Port: 10002,
			},
		},
	}
	if err = EnsureDbCluster(rt, authInfo, dbCluster); err != nil {
		log.Fatalf("got error %s\n%+v", spew.Sdump(err), err)
	}
}
