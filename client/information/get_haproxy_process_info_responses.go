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

package information

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/haproxytech/models"
)

// GetHaproxyProcessInfoReader is a Reader for the GetHaproxyProcessInfo structure.
type GetHaproxyProcessInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetHaproxyProcessInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetHaproxyProcessInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetHaproxyProcessInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetHaproxyProcessInfoOK creates a GetHaproxyProcessInfoOK with default headers values
func NewGetHaproxyProcessInfoOK() *GetHaproxyProcessInfoOK {
	return &GetHaproxyProcessInfoOK{}
}

/*GetHaproxyProcessInfoOK handles this case with default header values.

Success
*/
type GetHaproxyProcessInfoOK struct {
	Payload *models.ProcessInfo
}

func (o *GetHaproxyProcessInfoOK) Error() string {
	return fmt.Sprintf("[GET /services/haproxy/info][%d] getHaproxyProcessInfoOK  %+v", 200, o.Payload)
}

func (o *GetHaproxyProcessInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ProcessInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetHaproxyProcessInfoDefault creates a GetHaproxyProcessInfoDefault with default headers values
func NewGetHaproxyProcessInfoDefault(code int) *GetHaproxyProcessInfoDefault {
	return &GetHaproxyProcessInfoDefault{
		_statusCode: code,
	}
}

/*GetHaproxyProcessInfoDefault handles this case with default header values.

General Error
*/
type GetHaproxyProcessInfoDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get haproxy process info default response
func (o *GetHaproxyProcessInfoDefault) Code() int {
	return o._statusCode
}

func (o *GetHaproxyProcessInfoDefault) Error() string {
	return fmt.Sprintf("[GET /services/haproxy/info][%d] getHaproxyProcessInfo default  %+v", o._statusCode, o.Payload)
}

func (o *GetHaproxyProcessInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}