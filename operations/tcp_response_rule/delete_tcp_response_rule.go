// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2019 HAProxy Technologies LLC
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

package tcp_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteTCPResponseRuleHandlerFunc turns a function with the right signature into a delete TCP response rule handler
type DeleteTCPResponseRuleHandlerFunc func(DeleteTCPResponseRuleParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteTCPResponseRuleHandlerFunc) Handle(params DeleteTCPResponseRuleParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteTCPResponseRuleHandler interface for that can handle valid delete TCP response rule params
type DeleteTCPResponseRuleHandler interface {
	Handle(DeleteTCPResponseRuleParams, interface{}) middleware.Responder
}

// NewDeleteTCPResponseRule creates a new http.Handler for the delete TCP response rule operation
func NewDeleteTCPResponseRule(ctx *middleware.Context, handler DeleteTCPResponseRuleHandler) *DeleteTCPResponseRule {
	return &DeleteTCPResponseRule{Context: ctx, Handler: handler}
}

/*DeleteTCPResponseRule swagger:route DELETE /services/haproxy/configuration/tcp_response_rules/{id} TCPResponseRule deleteTcpResponseRule

Delete a TCP Response Rule

Deletes a TCP Response Rule configuration by it's ID from the specified backend.

*/
type DeleteTCPResponseRule struct {
	Context *middleware.Context
	Handler DeleteTCPResponseRuleHandler
}

func (o *DeleteTCPResponseRule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteTCPResponseRuleParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
