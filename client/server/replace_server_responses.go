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

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/haproxytech/models"
)

// ReplaceServerReader is a Reader for the ReplaceServer structure.
type ReplaceServerReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ReplaceServerReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewReplaceServerOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 202:
		result := NewReplaceServerAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewReplaceServerBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewReplaceServerNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewReplaceServerDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewReplaceServerOK creates a ReplaceServerOK with default headers values
func NewReplaceServerOK() *ReplaceServerOK {
	return &ReplaceServerOK{}
}

/*ReplaceServerOK handles this case with default header values.

Server replaced
*/
type ReplaceServerOK struct {
	Payload *models.Server
}

func (o *ReplaceServerOK) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/servers/{name}][%d] replaceServerOK  %+v", 200, o.Payload)
}

func (o *ReplaceServerOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Server)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceServerAccepted creates a ReplaceServerAccepted with default headers values
func NewReplaceServerAccepted() *ReplaceServerAccepted {
	return &ReplaceServerAccepted{}
}

/*ReplaceServerAccepted handles this case with default header values.

Configuration change accepted and reload requested
*/
type ReplaceServerAccepted struct {
	/*ID of the requested reload
	 */
	ReloadID string

	Payload *models.Server
}

func (o *ReplaceServerAccepted) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/servers/{name}][%d] replaceServerAccepted  %+v", 202, o.Payload)
}

func (o *ReplaceServerAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Reload-ID
	o.ReloadID = response.GetHeader("Reload-ID")

	o.Payload = new(models.Server)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceServerBadRequest creates a ReplaceServerBadRequest with default headers values
func NewReplaceServerBadRequest() *ReplaceServerBadRequest {
	return &ReplaceServerBadRequest{}
}

/*ReplaceServerBadRequest handles this case with default header values.

Bad request
*/
type ReplaceServerBadRequest struct {
	Payload *models.Error
}

func (o *ReplaceServerBadRequest) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/servers/{name}][%d] replaceServerBadRequest  %+v", 400, o.Payload)
}

func (o *ReplaceServerBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceServerNotFound creates a ReplaceServerNotFound with default headers values
func NewReplaceServerNotFound() *ReplaceServerNotFound {
	return &ReplaceServerNotFound{}
}

/*ReplaceServerNotFound handles this case with default header values.

The specified resource was not found
*/
type ReplaceServerNotFound struct {
	Payload *models.Error
}

func (o *ReplaceServerNotFound) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/servers/{name}][%d] replaceServerNotFound  %+v", 404, o.Payload)
}

func (o *ReplaceServerNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceServerDefault creates a ReplaceServerDefault with default headers values
func NewReplaceServerDefault(code int) *ReplaceServerDefault {
	return &ReplaceServerDefault{
		_statusCode: code,
	}
}

/*ReplaceServerDefault handles this case with default header values.

General Error
*/
type ReplaceServerDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the replace server default response
func (o *ReplaceServerDefault) Code() int {
	return o._statusCode
}

func (o *ReplaceServerDefault) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/servers/{name}][%d] replaceServer default  %+v", o._statusCode, o.Payload)
}

func (o *ReplaceServerDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}