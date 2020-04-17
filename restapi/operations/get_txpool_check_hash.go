// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTxpoolCheckHashHandlerFunc turns a function with the right signature into a get txpool check hash handler
type GetTxpoolCheckHashHandlerFunc func(GetTxpoolCheckHashParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTxpoolCheckHashHandlerFunc) Handle(params GetTxpoolCheckHashParams) middleware.Responder {
	return fn(params)
}

// GetTxpoolCheckHashHandler interface for that can handle valid get txpool check hash params
type GetTxpoolCheckHashHandler interface {
	Handle(GetTxpoolCheckHashParams) middleware.Responder
}

// NewGetTxpoolCheckHash creates a new http.Handler for the get txpool check hash operation
func NewGetTxpoolCheckHash(ctx *middleware.Context, handler GetTxpoolCheckHashHandler) *GetTxpoolCheckHash {
	return &GetTxpoolCheckHash{Context: ctx, Handler: handler}
}

/*GetTxpoolCheckHash swagger:route GET /txpool/check/{hash} getTxpoolCheckHash

check the tx in hash

*/
type GetTxpoolCheckHash struct {
	Context *middleware.Context
	Handler GetTxpoolCheckHashHandler
}

func (o *GetTxpoolCheckHash) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTxpoolCheckHashParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
