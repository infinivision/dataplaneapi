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

package server_switching_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ReplaceServerSwitchingRuleHandlerFunc turns a function with the right signature into a replace server switching rule handler
type ReplaceServerSwitchingRuleHandlerFunc func(ReplaceServerSwitchingRuleParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ReplaceServerSwitchingRuleHandlerFunc) Handle(params ReplaceServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ReplaceServerSwitchingRuleHandler interface for that can handle valid replace server switching rule params
type ReplaceServerSwitchingRuleHandler interface {
	Handle(ReplaceServerSwitchingRuleParams, interface{}) middleware.Responder
}

// NewReplaceServerSwitchingRule creates a new http.Handler for the replace server switching rule operation
func NewReplaceServerSwitchingRule(ctx *middleware.Context, handler ReplaceServerSwitchingRuleHandler) *ReplaceServerSwitchingRule {
	return &ReplaceServerSwitchingRule{Context: ctx, Handler: handler}
}

/*ReplaceServerSwitchingRule swagger:route PUT /services/haproxy/configuration/server_switching_rules/{id} ServerSwitchingRule replaceServerSwitchingRule

Replace a Server Switching Rule

Replaces a Server Switching Rule configuration by it's ID in the specified backend.

*/
type ReplaceServerSwitchingRule struct {
	Context *middleware.Context
	Handler ReplaceServerSwitchingRuleHandler
}

func (o *ReplaceServerSwitchingRule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewReplaceServerSwitchingRuleParams()

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
