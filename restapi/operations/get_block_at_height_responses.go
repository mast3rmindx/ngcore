// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetBlockAtHeightOKCode is the HTTP code returned for type GetBlockAtHeightOK
const GetBlockAtHeightOKCode int = 200

/*GetBlockAtHeightOK OK

swagger:response getBlockAtHeightOK
*/
type GetBlockAtHeightOK struct {
}

// NewGetBlockAtHeightOK creates GetBlockAtHeightOK with default headers values
func NewGetBlockAtHeightOK() *GetBlockAtHeightOK {

	return &GetBlockAtHeightOK{}
}

// WriteResponse to the client
func (o *GetBlockAtHeightOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// GetBlockAtHeightBadRequestCode is the HTTP code returned for type GetBlockAtHeightBadRequest
const GetBlockAtHeightBadRequestCode int = 400

/*GetBlockAtHeightBadRequest Error

swagger:response getBlockAtHeightBadRequest
*/
type GetBlockAtHeightBadRequest struct {
}

// NewGetBlockAtHeightBadRequest creates GetBlockAtHeightBadRequest with default headers values
func NewGetBlockAtHeightBadRequest() *GetBlockAtHeightBadRequest {

	return &GetBlockAtHeightBadRequest{}
}

// WriteResponse to the client
func (o *GetBlockAtHeightBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}
