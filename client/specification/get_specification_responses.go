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

package specification

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/haproxytech/models"
)

// GetSpecificationReader is a Reader for the GetSpecification structure.
type GetSpecificationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSpecificationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSpecificationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetSpecificationDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetSpecificationOK creates a GetSpecificationOK with default headers values
func NewGetSpecificationOK() *GetSpecificationOK {
	return &GetSpecificationOK{}
}

/*GetSpecificationOK handles this case with default header values.

Success
*/
type GetSpecificationOK struct {
	Payload interface{}
}

func (o *GetSpecificationOK) Error() string {
	return fmt.Sprintf("[GET /specification][%d] getSpecificationOK  %+v", 200, o.Payload)
}

func (o *GetSpecificationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSpecificationDefault creates a GetSpecificationDefault with default headers values
func NewGetSpecificationDefault(code int) *GetSpecificationDefault {
	return &GetSpecificationDefault{
		_statusCode: code,
	}
}

/*GetSpecificationDefault handles this case with default header values.

General Error
*/
type GetSpecificationDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get specification default response
func (o *GetSpecificationDefault) Code() int {
	return o._statusCode
}

func (o *GetSpecificationDefault) Error() string {
	return fmt.Sprintf("[GET /specification][%d] getSpecification default  %+v", o._statusCode, o.Payload)
}

func (o *GetSpecificationDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}