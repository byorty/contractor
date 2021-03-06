// Code generated by go-swagger; DO NOT EDIT.

package saved

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewSearchRelativeParams creates a new SearchRelativeParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchRelativeParams() *SearchRelativeParams {
	return &SearchRelativeParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchRelativeParamsWithTimeout creates a new SearchRelativeParams object
// with the ability to set a timeout on a request.
func NewSearchRelativeParamsWithTimeout(timeout time.Duration) *SearchRelativeParams {
	return &SearchRelativeParams{
		timeout: timeout,
	}
}

// NewSearchRelativeParamsWithContext creates a new SearchRelativeParams object
// with the ability to set a context for a request.
func NewSearchRelativeParamsWithContext(ctx context.Context) *SearchRelativeParams {
	return &SearchRelativeParams{
		Context: ctx,
	}
}

// NewSearchRelativeParamsWithHTTPClient creates a new SearchRelativeParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchRelativeParamsWithHTTPClient(client *http.Client) *SearchRelativeParams {
	return &SearchRelativeParams{
		HTTPClient: client,
	}
}

/* SearchRelativeParams contains all the parameters to send to the API endpoint
   for the search relative operation.

   Typically these are written to a http.Request.
*/
type SearchRelativeParams struct {

	/* Decorate.

	   Run decorators on search result

	   Default: true
	*/
	Decorate *bool

	/* Fields.

	   Comma separated list of fields to return
	*/
	Fields *string

	/* Filter.

	   Filter
	*/
	Filter *string

	/* Limit.

	   Maximum number of messages to return.
	*/
	Limit *int64

	/* Offset.

	   Offset
	*/
	Offset *int64

	/* Query.

	   Query (Lucene syntax)
	*/
	Query string

	/* Range.

	   Relative timeframe to search in. See method description.
	*/
	Range int64

	/* Sort.

	   Sorting (field:asc / field:desc)
	*/
	Sort *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search relative params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchRelativeParams) WithDefaults() *SearchRelativeParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search relative params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchRelativeParams) SetDefaults() {
	var (
		decorateDefault = bool(true)
	)

	val := SearchRelativeParams{
		Decorate: &decorateDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the search relative params
func (o *SearchRelativeParams) WithTimeout(timeout time.Duration) *SearchRelativeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search relative params
func (o *SearchRelativeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search relative params
func (o *SearchRelativeParams) WithContext(ctx context.Context) *SearchRelativeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search relative params
func (o *SearchRelativeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search relative params
func (o *SearchRelativeParams) WithHTTPClient(client *http.Client) *SearchRelativeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search relative params
func (o *SearchRelativeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDecorate adds the decorate to the search relative params
func (o *SearchRelativeParams) WithDecorate(decorate *bool) *SearchRelativeParams {
	o.SetDecorate(decorate)
	return o
}

// SetDecorate adds the decorate to the search relative params
func (o *SearchRelativeParams) SetDecorate(decorate *bool) {
	o.Decorate = decorate
}

// WithFields adds the fields to the search relative params
func (o *SearchRelativeParams) WithFields(fields *string) *SearchRelativeParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the search relative params
func (o *SearchRelativeParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithFilter adds the filter to the search relative params
func (o *SearchRelativeParams) WithFilter(filter *string) *SearchRelativeParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the search relative params
func (o *SearchRelativeParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithLimit adds the limit to the search relative params
func (o *SearchRelativeParams) WithLimit(limit *int64) *SearchRelativeParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the search relative params
func (o *SearchRelativeParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithOffset adds the offset to the search relative params
func (o *SearchRelativeParams) WithOffset(offset *int64) *SearchRelativeParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the search relative params
func (o *SearchRelativeParams) SetOffset(offset *int64) {
	o.Offset = offset
}

// WithQuery adds the query to the search relative params
func (o *SearchRelativeParams) WithQuery(query string) *SearchRelativeParams {
	o.SetQuery(query)
	return o
}

// SetQuery adds the query to the search relative params
func (o *SearchRelativeParams) SetQuery(query string) {
	o.Query = query
}

// WithRange adds the rangeVar to the search relative params
func (o *SearchRelativeParams) WithRange(rangeVar int64) *SearchRelativeParams {
	o.SetRange(rangeVar)
	return o
}

// SetRange adds the range to the search relative params
func (o *SearchRelativeParams) SetRange(rangeVar int64) {
	o.Range = rangeVar
}

// WithSort adds the sort to the search relative params
func (o *SearchRelativeParams) WithSort(sort *string) *SearchRelativeParams {
	o.SetSort(sort)
	return o
}

// SetSort adds the sort to the search relative params
func (o *SearchRelativeParams) SetSort(sort *string) {
	o.Sort = sort
}

// WriteToRequest writes these params to a swagger request
func (o *SearchRelativeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Decorate != nil {

		// query param decorate
		var qrDecorate bool

		if o.Decorate != nil {
			qrDecorate = *o.Decorate
		}
		qDecorate := swag.FormatBool(qrDecorate)
		if qDecorate != "" {

			if err := r.SetQueryParam("decorate", qDecorate); err != nil {
				return err
			}
		}
	}

	if o.Fields != nil {

		// query param fields
		var qrFields string

		if o.Fields != nil {
			qrFields = *o.Fields
		}
		qFields := qrFields
		if qFields != "" {

			if err := r.SetQueryParam("fields", qFields); err != nil {
				return err
			}
		}
	}

	if o.Filter != nil {

		// query param filter
		var qrFilter string

		if o.Filter != nil {
			qrFilter = *o.Filter
		}
		qFilter := qrFilter
		if qFilter != "" {

			if err := r.SetQueryParam("filter", qFilter); err != nil {
				return err
			}
		}
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64

		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {

			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}
	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int64

		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt64(qrOffset)
		if qOffset != "" {

			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}
	}

	// query param query
	qrQuery := o.Query
	qQuery := qrQuery
	if qQuery != "" {

		if err := r.SetQueryParam("query", qQuery); err != nil {
			return err
		}
	}

	// query param range
	qrRange := o.Range
	qRange := swag.FormatInt64(qrRange)
	if qRange != "" {

		if err := r.SetQueryParam("range", qRange); err != nil {
			return err
		}
	}

	if o.Sort != nil {

		// query param sort
		var qrSort string

		if o.Sort != nil {
			qrSort = *o.Sort
		}
		qSort := qrSort
		if qSort != "" {

			if err := r.SetQueryParam("sort", qSort); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
