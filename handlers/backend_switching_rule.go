// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/backend_switching_rule"
	"github.com/haproxytech/models"
)

//CreateBackendSwitchingRuleHandlerImpl implementation of the CreateBackendSwitchingRuleHandler interface using client-native client
type CreateBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteBackendSwitchingRuleHandlerImpl implementation of the DeleteBackendSwitchingRuleHandler interface using client-native client
type DeleteBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetBackendSwitchingRuleHandlerImpl implementation of the GetBackendSwitchingRuleHandler interface using client-native client
type GetBackendSwitchingRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetBackendSwitchingRulesHandlerImpl implementation of the GetBackendSwitchingRulesHandler interface using client-native client
type GetBackendSwitchingRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceBackendSwitchingRuleHandlerImpl implementation of the ReplaceBackendSwitchingRuleHandler interface using client-native client
type ReplaceBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.CreateBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}
	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.CreateBackendSwitchingRule(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewCreateBackendSwitchingRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewCreateBackendSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend_switching_rule.NewCreateBackendSwitchingRuleAccepted().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.DeleteBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.DeleteBackendSwitchingRule(params.ID, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewDeleteBackendSwitchingRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewDeleteBackendSwitchingRuleAccepted().WithReloadID(rID)
	}
	return backend_switching_rule.NewDeleteBackendSwitchingRuleAccepted()
}

//Handle executing the request and returning a response
func (h *GetBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, bckRule, err := h.Client.Configuration.GetBackendSwitchingRule(params.ID, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRuleOK().WithPayload(&backend_switching_rule.GetBackendSwitchingRuleOKBody{Version: v, Data: bckRule})
}

//Handle executing the request and returning a response
func (h *GetBackendSwitchingRulesHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, bckRules, err := h.Client.Configuration.GetBackendSwitchingRules(params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRulesOK().WithPayload(&backend_switching_rule.GetBackendSwitchingRulesOKBody{Version: v, Data: bckRules})
}

//Handle executing the request and returning a response
func (h *ReplaceBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.ReplaceBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.EditBackendSwitchingRule(params.ID, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewReplaceBackendSwitchingRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewReplaceBackendSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend_switching_rule.NewReplaceBackendSwitchingRuleAccepted().WithPayload(params.Data)
}
