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

package transactions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/haproxytech/models"
)

// GetTransactionReader is a Reader for the GetTransaction structure.
type GetTransactionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetTransactionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetTransactionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewGetTransactionNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewGetTransactionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetTransactionOK creates a GetTransactionOK with default headers values
func NewGetTransactionOK() *GetTransactionOK {
	return &GetTransactionOK{}
}

/*GetTransactionOK handles this case with default header values.

Successful operation
*/
type GetTransactionOK struct {
	Payload *models.Transaction
}

func (o *GetTransactionOK) Error() string {
	return fmt.Sprintf("[GET /services/haproxy/transactions/{id}][%d] getTransactionOK  %+v", 200, o.Payload)
}

func (o *GetTransactionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Transaction)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTransactionNotFound creates a GetTransactionNotFound with default headers values
func NewGetTransactionNotFound() *GetTransactionNotFound {
	return &GetTransactionNotFound{}
}

/*GetTransactionNotFound handles this case with default header values.

The specified resource was not found
*/
type GetTransactionNotFound struct {
	Payload *models.Error
}

func (o *GetTransactionNotFound) Error() string {
	return fmt.Sprintf("[GET /services/haproxy/transactions/{id}][%d] getTransactionNotFound  %+v", 404, o.Payload)
}

func (o *GetTransactionNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTransactionDefault creates a GetTransactionDefault with default headers values
func NewGetTransactionDefault(code int) *GetTransactionDefault {
	return &GetTransactionDefault{
		_statusCode: code,
	}
}

/*GetTransactionDefault handles this case with default header values.

General Error
*/
type GetTransactionDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get transaction default response
func (o *GetTransactionDefault) Code() int {
	return o._statusCode
}

func (o *GetTransactionDefault) Error() string {
	return fmt.Sprintf("[GET /services/haproxy/transactions/{id}][%d] getTransaction default  %+v", o._statusCode, o.Payload)
}

func (o *GetTransactionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}