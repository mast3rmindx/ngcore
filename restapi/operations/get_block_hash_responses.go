// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetBlockHashOKCode is the HTTP code returned for type GetBlockHashOK
const GetBlockHashOKCode int = 200

/*GetBlockHashOK OK

swagger:response getBlockHashOK
*/
type GetBlockHashOK struct {
}

// NewGetBlockHashOK creates GetBlockHashOK with default headers values
func NewGetBlockHashOK() *GetBlockHashOK {

	return &GetBlockHashOK{}
}

// WriteResponse to the client
func (o *GetBlockHashOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// GetBlockHashBadRequestCode is the HTTP code returned for type GetBlockHashBadRequest
const GetBlockHashBadRequestCode int = 400

/*GetBlockHashBadRequest Error

swagger:response getBlockHashBadRequest
*/
type GetBlockHashBadRequest struct {
}

// NewGetBlockHashBadRequest creates GetBlockHashBadRequest with default headers values
func NewGetBlockHashBadRequest() *GetBlockHashBadRequest {

	return &GetBlockHashBadRequest{}
}

// WriteResponse to the client
func (o *GetBlockHashBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}
