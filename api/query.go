package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/0xArch3r/goforce/types"
)

// Select is (for now) simple queries
type Select func(object string, opts ...SelectOption) (*types.QueryResult, error)

func newSelectFunc(base Transport) Select {
	return func(object string, opts ...SelectOption) (*types.QueryResult, error) {
		r := SelectRequest{
			Object: object,
			Fields: []string{"FIELDS(ALL)"},
		}
		for _, f := range opts {
			err := f(&r)
			if err != nil {
				return nil, err
			}
		}

		resp, err := r.Do(r.ctx, base)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		res := &types.QueryResult{}
		err = json.Unmarshal(data, res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

type SelectOption func(*SelectRequest) error

type SelectRequest struct {
	ctx     context.Context
	Object  string
	Fields  []string
	Limit   int
	OrderBy string
}

func (r SelectRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	method := http.MethodGet

	fields := strings.Join(r.Fields, ",")

	query := fmt.Sprintf("SELECT %v FROM %v", fields, r.Object)

	if r.OrderBy != "" {
		query = fmt.Sprintf("%v ORDER BY %v", query, r.OrderBy)
	}

	if r.Limit > 0 {
		query = fmt.Sprintf("%v LIMIT %v", query, r.Limit)
	}

	path := fmt.Sprintf("/query?q=%s", url.QueryEscape(query))

	req, err := http.NewRequest(method, path, nil)
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
func (f Select) WithContext(v context.Context) SelectOption {
	return func(r *SelectRequest) error {
		r.ctx = v
		return nil
	}
}

func (f Select) Fields(fields ...string) SelectOption {
	return func(r *SelectRequest) error {
		r.Fields = fields
		return nil
	}
}

func (f Select) Limit(limit int) SelectOption {
	return func(r *SelectRequest) error {
		if limit < 1 {
			return errors.New("limit cannot be lower than 1")
		}
		r.Limit = limit
		return nil
	}
}

func (f Select) OrderBy(orderBy string) SelectOption {
	return func(r *SelectRequest) error {
		r.OrderBy = orderBy
		return nil
	}
}
