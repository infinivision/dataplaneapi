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
	"github.com/haproxytech/dataplaneapi/operations/stick_rule"
	"github.com/haproxytech/models"
)

//CreateStickRuleHandlerImpl implementation of the CreateStickRuleHandler interface using client-native client
type CreateStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteStickRuleHandlerImpl implementation of the DeleteStickRuleHandler interface using client-native client
type DeleteStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetStickRuleHandlerImpl implementation of the GetStickRuleHandler interface using client-native client
type GetStickRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetStickRulesHandlerImpl implementation of the GetStickRulesHandler interface using client-native client
type GetStickRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceStickRuleHandlerImpl implementation of the ReplaceStickRuleHandler interface using client-native client
type ReplaceStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateStickRuleHandlerImpl) Handle(params stick_rule.CreateStickRuleParams, principal interface{}) middleware.Responder {
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
		return stick_rule.NewCreateStickRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.CreateStickRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewCreateStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return stick_rule.NewCreateStickRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return stick_rule.NewCreateStickRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return stick_rule.NewCreateStickRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return stick_rule.NewCreateStickRuleAccepted().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteStickRuleHandlerImpl) Handle(params stick_rule.DeleteStickRuleParams, principal interface{}) middleware.Responder {
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
		return stick_rule.NewDeleteStickRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.DeleteStickRule(params.ID, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewDeleteStickRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return stick_rule.NewDeleteStickRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return stick_rule.NewDeleteStickRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return stick_rule.NewDeleteStickRuleAccepted().WithReloadID(rID)
	}
	return stick_rule.NewDeleteStickRuleAccepted()
}

//Handle executing the request and returning a response
func (h *GetStickRuleHandlerImpl) Handle(params stick_rule.GetStickRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, rule, err := h.Client.Configuration.GetStickRule(params.ID, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewGetStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_rule.NewGetStickRuleOK().WithPayload(&stick_rule.GetStickRuleOKBody{Version: v, Data: rule})
}

//Handle executing the request and returning a response
func (h *GetStickRulesHandlerImpl) Handle(params stick_rule.GetStickRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, rules, err := h.Client.Configuration.GetStickRules(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewGetStickRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_rule.NewGetStickRulesOK().WithPayload(&stick_rule.GetStickRulesOKBody{Version: v, Data: rules})
}

//Handle executing the request and returning a response
func (h *ReplaceStickRuleHandlerImpl) Handle(params stick_rule.ReplaceStickRuleParams, principal interface{}) middleware.Responder {
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
		return stick_rule.NewReplaceStickRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.EditStickRule(params.ID, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewReplaceStickRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return stick_rule.NewReplaceStickRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return stick_rule.NewReplaceStickRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return stick_rule.NewReplaceStickRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return stick_rule.NewReplaceStickRuleAccepted().WithPayload(params.Data)
}
