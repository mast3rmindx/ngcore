// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetBalanceMyHandlerFunc turns a function with the right signature into a get balance my handler
type GetBalanceMyHandlerFunc func(GetBalanceMyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBalanceMyHandlerFunc) Handle(params GetBalanceMyParams) middleware.Responder {
	return fn(params)
}

// GetBalanceMyHandler interface for that can handle valid get balance my params
type GetBalanceMyHandler interface {
	Handle(GetBalanceMyParams) middleware.Responder
}

// NewGetBalanceMy creates a new http.Handler for the get balance my operation
func NewGetBalanceMy(ctx *middleware.Context, handler GetBalanceMyHandler) *GetBalanceMy {
	return &GetBalanceMy{Context: ctx, Handler: handler}
}

/*GetBalanceMy swagger:route GET /balance/my getBalanceMy

get all my accounts

*/
type GetBalanceMy struct {
	Context *middleware.Context
	Handler GetBalanceMyHandler
}

func (o *GetBalanceMy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetBalanceMyParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
