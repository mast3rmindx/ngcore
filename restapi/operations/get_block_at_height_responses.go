// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ngchain/ngcore/models"
)

// GetBlockAtHeightOKCode is the HTTP code returned for type GetBlockAtHeightOK
const GetBlockAtHeightOKCode int = 200

/*GetBlockAtHeightOK OK

swagger:response getBlockAtHeightOK
*/
type GetBlockAtHeightOK struct {

	/*
	  In: Body
	*/
	Payload *models.Block `json:"body,omitempty"`
}

// NewGetBlockAtHeightOK creates GetBlockAtHeightOK with default headers values
func NewGetBlockAtHeightOK() *GetBlockAtHeightOK {

	return &GetBlockAtHeightOK{}
}

// WithPayload adds the payload to the get block at height o k response
func (o *GetBlockAtHeightOK) WithPayload(payload *models.Block) *GetBlockAtHeightOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get block at height o k response
func (o *GetBlockAtHeightOK) SetPayload(payload *models.Block) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBlockAtHeightOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetBlockAtHeightBadRequestCode is the HTTP code returned for type GetBlockAtHeightBadRequest
const GetBlockAtHeightBadRequestCode int = 400

/*GetBlockAtHeightBadRequest Error

swagger:response getBlockAtHeightBadRequest
*/
type GetBlockAtHeightBadRequest struct {

	/*error text
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetBlockAtHeightBadRequest creates GetBlockAtHeightBadRequest with default headers values
func NewGetBlockAtHeightBadRequest() *GetBlockAtHeightBadRequest {

	return &GetBlockAtHeightBadRequest{}
}

// WithPayload adds the payload to the get block at height bad request response
func (o *GetBlockAtHeightBadRequest) WithPayload(payload string) *GetBlockAtHeightBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get block at height bad request response
func (o *GetBlockAtHeightBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBlockAtHeightBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}