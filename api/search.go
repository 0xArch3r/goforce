package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/0xArch3r/goforce/types"
)

func newSearchFunc(b Transport) Search {
	return func(query string, o ...SearchOption) (*types.SearchResults, error) {
		r := SearchRequest{
			Query:        query,
			In:           "ALL",
			DefaultLimit: 2000,
			OverallLimit: 2000,
		}
		for _, f := range o {
			err := f(&r)
			if err != nil {
				return nil, err
			}
		}

		resp, err := r.Do(r.ctx, b)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.IsError() {
			err := types.ParseSalesforceError(resp.StatusCode, data)
			return nil, err
		}

		res := &types.SearchResults{}
		err = json.Unmarshal(data, res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

type Search func(query string, o ...SearchOption) (*types.SearchResults, error)

type SearchOption func(*SearchRequest) error

type SearchObject struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}

// SearchRequest configures the Search API request.
type SearchRequest struct {
	Query        string         `json:"q"`
	Fields       []string       `json:"fields,omitempty"`
	SObjects     []SearchObject `json:"sobjects,omitempty"`
	In           string         `json:"in"`
	OverallLimit int            `json:"overallLimit,omitempty"`
	DefaultLimit int            `json:"defaultLimit,omitempty"`

	ctx context.Context
}

// Do executes the request and returns response or error.
func (r SearchRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	method := http.MethodPost
	path := "http:///parameterizedSearch"

	payload, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	} else {
		req = req.WithContext(context.Background())
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// WithContext sets the request context.
func (f Search) WithContext(v context.Context) SearchOption {
	return func(r *SearchRequest) error {
		r.ctx = v
		return nil
	}
}

func (f Search) WithFields(fields ...string) SearchOption {
	return func(r *SearchRequest) error {
		r.Fields = fields
		return nil
	}
}

// TODO: Check for valid searchable fields, apparently only some are searchable
func (f Search) In(fields ...string) SearchOption {
	return func(r *SearchRequest) error {
		r.In = strings.Join(fields, ",")
		return nil
	}
}

func (f Search) SObjects(objects ...SearchObject) SearchOption {
	return func(r *SearchRequest) error {
		r.SObjects = objects
		return nil
	}
}

func (f Search) WithOverallLimit(limit int) SearchOption {
	return func(r *SearchRequest) error {
		if limit > 2000 {
			return errors.New("limit exceeds maximum")
		} else if limit <= 0 {
			return errors.New("limit must be greater than 0")
		}
		r.OverallLimit = limit
		return nil
	}
}

func (f Search) WithDefaultlLimit(limit int) SearchOption {
	return func(r *SearchRequest) error {
		if limit > 2000 {
			return errors.New("limit exceeds maximum")
		} else if limit <= 0 {
			return errors.New("limit must be greater than 0")
		}
		r.DefaultLimit = limit
		return nil
	}
}
