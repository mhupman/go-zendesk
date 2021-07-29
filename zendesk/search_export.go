package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
)

// SearchExportOptions are the options that can be provided to the search export API
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/search/#parameters-2
type SearchExportOptions struct {
    CursorOption
    PageSize int `url:"page[size]"`
    FilterType string `url:"filter[type]"`
	Query     string `url:"query"`
}

// searchExportOptions is a shadow of SearchExportOptions used for encoding cursor
// option as a different query paramater while allowing users to use the common
// CursorOption struct
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/search/#parameters-2
type searchExportOptions struct {
    StartTime int64  `url:"start_time,omitempty"`
    Cursor    string `url:"page[after],omitempty"`
    PageSize int `url:"page[size]"`
    FilterType string `url:"filter[type]"`
    Query     string `url:"query"`
}

type SearchExportResultsMeta struct {
    HasMore bool `json:"has_more"`
    Cursor
}

type SearchExportAPI interface {
	SearchExport(ctx context.Context, opts *SearchExportOptions) (SearchExportResults, SearchExportResultsMeta, error)
}

type SearchExportResults struct {
	results []interface{}
}

func (r *SearchExportResults) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.results)
}

func (r *SearchExportResults) UnmarshalJSON(b []byte) error {
	var (
		results []interface{}
		tmp     []json.RawMessage
	)

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	for _, v := range tmp {
		value, err := r.getObject(v)
		if err != nil {
			return err
		}

		results = append(results, value)
	}

	r.results = results

	return nil
}

func (r *SearchExportResults) getObject(blob json.RawMessage) (interface{}, error) {
	m := make(map[string]interface{})

	err := json.Unmarshal(blob, &m)
	if err != nil {
		return nil, err
	}

	t, ok := m["result_type"].(string)
	if !ok {
		return nil, fmt.Errorf("could not assert result type to string. json was: %v", blob)
	}

	var value interface{}

	switch t {
	case "group":
		var g Group
		err = json.Unmarshal(blob, &g)
		value = g
	case "ticket":
		var t Ticket
		err = json.Unmarshal(blob, &t)
		value = t
	case "user":
		var u User
		err = json.Unmarshal(blob, &u)
		value = u
	case "organization":
		var o Organization
		err = json.Unmarshal(blob, &o)
		value = o
	case "topic":
		var t Topic
		err = json.Unmarshal(blob, &t)
		value = t
	default:
		err = fmt.Errorf("value of result was an unsupported type %s", t)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

// String return string formatted for Search Export results
func (r *SearchExportResults) String() string {
	return fmt.Sprintf("%v", r.results)
}

// List return internal array in Search Export Results
func (r *SearchExportResults) List() []interface{} {
	return r.results
}

// SearchExport allows users to query zendesk's search export endpoint.
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/search/#export-search-results
func (z *Client) SearchExport(ctx context.Context, opts *SearchExportOptions) (SearchExportResults, SearchExportResultsMeta, error) {
    var data struct {
        Results SearchExportResults `json:"results"`
        Meta SearchExportResultsMeta `json:"meta"`
    }

    if opts == nil {
        return SearchExportResults{}, SearchExportResultsMeta{}, &OptionsError{opts}
    }

    customOpts := searchExportOptions{
        StartTime: opts.StartTime,
        Cursor:    opts.Cursor,
        PageSize:     opts.PageSize,
        FilterType:   opts.FilterType,
        Query:        opts.Query,
    }
    u, err := addOptions("/search/export", customOpts)
    if err != nil {
        return SearchExportResults{}, SearchExportResultsMeta{}, err
    }

    body, err := z.get(ctx, u)
    if err != nil {
        return SearchExportResults{}, SearchExportResultsMeta{}, err
    }

    err = json.Unmarshal(body, &data)
    if err != nil {
        return SearchExportResults{}, SearchExportResultsMeta{}, err
    }

    return data.Results, data.Meta, nil
}
