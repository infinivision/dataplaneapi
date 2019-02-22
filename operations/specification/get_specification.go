// Code generated by go-swagger; DO NOT EDIT.

package specification

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetSpecificationHandlerFunc turns a function with the right signature into a get specification handler
type GetSpecificationHandlerFunc func(GetSpecificationParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetSpecificationHandlerFunc) Handle(params GetSpecificationParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetSpecificationHandler interface for that can handle valid get specification params
type GetSpecificationHandler interface {
	Handle(GetSpecificationParams, interface{}) middleware.Responder
}

// NewGetSpecification creates a new http.Handler for the get specification operation
func NewGetSpecification(ctx *middleware.Context, handler GetSpecificationHandler) *GetSpecification {
	return &GetSpecification{Context: ctx, Handler: handler}
}

/*GetSpecification swagger:route GET /specification Specification getSpecification

Dataplane API Specification

Return Dataplane API OpenAPI specification

*/
type GetSpecification struct {
	Context *middleware.Context
	Handler GetSpecificationHandler
}

func (o *GetSpecification) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetSpecificationParams()

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

// GetSpecificationOKBody get specification o k body
// swagger:model GetSpecificationOKBody
type GetSpecificationOKBody interface{}