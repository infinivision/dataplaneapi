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

package stats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/models"
)

// GetStatsOKCode is the HTTP code returned for type GetStatsOK
const GetStatsOKCode int = 200

/*GetStatsOK Success

swagger:response getStatsOK
*/
type GetStatsOK struct {

	/*
	  In: Body
	*/
	Payload models.NativeStats `json:"body,omitempty"`
}

// NewGetStatsOK creates GetStatsOK with default headers values
func NewGetStatsOK() *GetStatsOK {

	return &GetStatsOK{}
}

// WithPayload adds the payload to the get stats o k response
func (o *GetStatsOK) WithPayload(payload models.NativeStats) *GetStatsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stats o k response
func (o *GetStatsOK) SetPayload(payload models.NativeStats) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStatsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.NativeStats{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GetStatsDefault General Error

swagger:response getStatsDefault
*/
type GetStatsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetStatsDefault creates GetStatsDefault with default headers values
func NewGetStatsDefault(code int) *GetStatsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetStatsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get stats default response
func (o *GetStatsDefault) WithStatusCode(code int) *GetStatsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get stats default response
func (o *GetStatsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get stats default response
func (o *GetStatsDefault) WithPayload(payload *models.Error) *GetStatsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stats default response
func (o *GetStatsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStatsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
