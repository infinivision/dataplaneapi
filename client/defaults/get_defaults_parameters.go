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

package defaults

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetDefaultsParams creates a new GetDefaultsParams object
// with the default values initialized.
func NewGetDefaultsParams() *GetDefaultsParams {
	var ()
	return &GetDefaultsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetDefaultsParamsWithTimeout creates a new GetDefaultsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetDefaultsParamsWithTimeout(timeout time.Duration) *GetDefaultsParams {
	var ()
	return &GetDefaultsParams{

		timeout: timeout,
	}
}

// NewGetDefaultsParamsWithContext creates a new GetDefaultsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetDefaultsParamsWithContext(ctx context.Context) *GetDefaultsParams {
	var ()
	return &GetDefaultsParams{

		Context: ctx,
	}
}

// NewGetDefaultsParamsWithHTTPClient creates a new GetDefaultsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetDefaultsParamsWithHTTPClient(client *http.Client) *GetDefaultsParams {
	var ()
	return &GetDefaultsParams{
		HTTPClient: client,
	}
}

/*GetDefaultsParams contains all the parameters to send to the API endpoint
for the get defaults operation typically these are written to a http.Request
*/
type GetDefaultsParams struct {

	/*TransactionID
	  ID of the transaction where we want to add the operation. Cannot be used when version is specified.

	*/
	TransactionID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get defaults params
func (o *GetDefaultsParams) WithTimeout(timeout time.Duration) *GetDefaultsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get defaults params
func (o *GetDefaultsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get defaults params
func (o *GetDefaultsParams) WithContext(ctx context.Context) *GetDefaultsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get defaults params
func (o *GetDefaultsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get defaults params
func (o *GetDefaultsParams) WithHTTPClient(client *http.Client) *GetDefaultsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get defaults params
func (o *GetDefaultsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTransactionID adds the transactionID to the get defaults params
func (o *GetDefaultsParams) WithTransactionID(transactionID *string) *GetDefaultsParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the get defaults params
func (o *GetDefaultsParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetDefaultsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.TransactionID != nil {

		// query param transaction_id
		var qrTransactionID string
		if o.TransactionID != nil {
			qrTransactionID = *o.TransactionID
		}
		qTransactionID := qrTransactionID
		if qTransactionID != "" {
			if err := r.SetQueryParam("transaction_id", qTransactionID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}