// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetAccountMyOKCode is the HTTP code returned for type GetAccountMyOK
const GetAccountMyOKCode int = 200

/*GetAccountMyOK OK

swagger:response getAccountMyOK
*/
type GetAccountMyOK struct {
}

// NewGetAccountMyOK creates GetAccountMyOK with default headers values
func NewGetAccountMyOK() *GetAccountMyOK {

	return &GetAccountMyOK{}
}

// WriteResponse to the client
func (o *GetAccountMyOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// GetAccountMyBadRequestCode is the HTTP code returned for type GetAccountMyBadRequest
const GetAccountMyBadRequestCode int = 400

/*GetAccountMyBadRequest Error

swagger:response getAccountMyBadRequest
*/
type GetAccountMyBadRequest struct {
}

// NewGetAccountMyBadRequest creates GetAccountMyBadRequest with default headers values
func NewGetAccountMyBadRequest() *GetAccountMyBadRequest {

	return &GetAccountMyBadRequest{}
}

// WriteResponse to the client
func (o *GetAccountMyBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}
