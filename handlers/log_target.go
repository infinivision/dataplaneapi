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
	"github.com/haproxytech/dataplaneapi/operations/log_target"
	"github.com/haproxytech/models"
)

//CreateLogTargetHandlerImpl implementation of the CreateLogTargetHandler interface using client-native client
type CreateLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteLogTargetHandlerImpl implementation of the DeleteLogTargetHandler interface using client-native client
type DeleteLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetLogTargetHandlerImpl implementation of the GetLogTargetHandler interface using client-native client
type GetLogTargetHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetLogTargetsHandlerImpl implementation of the GetLogTargetsHandler interface using client-native client
type GetLogTargetsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceLogTargetHandlerImpl implementation of the ReplaceLogTargetHandler interface using client-native client
type ReplaceLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateLogTargetHandlerImpl) Handle(params log_target.CreateLogTargetParams, principal interface{}) middleware.Responder {
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
		return log_target.NewCreateLogTargetDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.CreateLogTarget(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewCreateLogTargetDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewCreateLogTargetDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewCreateLogTargetCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewCreateLogTargetAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_target.NewCreateLogTargetAccepted().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteLogTargetHandlerImpl) Handle(params log_target.DeleteLogTargetParams, principal interface{}) middleware.Responder {
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
		return log_target.NewDeleteLogTargetDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.DeleteLogTarget(params.ID, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewDeleteLogTargetDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewDeleteLogTargetDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewDeleteLogTargetNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewDeleteLogTargetAccepted().WithReloadID(rID)
	}
	return log_target.NewDeleteLogTargetAccepted()
}

//Handle executing the request and returning a response
func (h *GetLogTargetHandlerImpl) Handle(params log_target.GetLogTargetParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, logTarget, err := h.Client.Configuration.GetLogTarget(params.ID, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetLogTargetOK().WithPayload(&log_target.GetLogTargetOKBody{Version: v, Data: logTarget})
}

//Handle executing the request and returning a response
func (h *GetLogTargetsHandlerImpl) Handle(params log_target.GetLogTargetsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, logTargets, err := h.Client.Configuration.GetLogTargets(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetsDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetLogTargetsOK().WithPayload(&log_target.GetLogTargetsOKBody{Version: v, Data: logTargets})
}

//Handle executing the request and returning a response
func (h *ReplaceLogTargetHandlerImpl) Handle(params log_target.ReplaceLogTargetParams, principal interface{}) middleware.Responder {
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
		return log_target.NewReplaceLogTargetDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.EditLogTarget(params.ID, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewReplaceLogTargetDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewReplaceLogTargetOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewReplaceLogTargetAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_target.NewReplaceLogTargetAccepted().WithPayload(params.Data)
}
