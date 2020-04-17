// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetAccountAtNumBalanceParams creates a new GetAccountAtNumBalanceParams object
// no default values defined in spec.
func NewGetAccountAtNumBalanceParams() GetAccountAtNumBalanceParams {

	return GetAccountAtNumBalanceParams{}
}

// GetAccountAtNumBalanceParams contains all the bound params for the get account at num balance operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetAccountAtNumBalance
type GetAccountAtNumBalanceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*number of account
	  Required: true
	  In: path
	*/
	Num int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetAccountAtNumBalanceParams() beforehand.
func (o *GetAccountAtNumBalanceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rNum, rhkNum, _ := route.Params.GetOK("num")
	if err := o.bindNum(rNum, rhkNum, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindNum binds and validates parameter Num from path.
func (o *GetAccountAtNumBalanceParams) bindNum(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("num", "path", "int64", raw)
	}
	o.Num = value

	return nil
}
