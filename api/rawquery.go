package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/0xArch3r/goforce/types"
)

// RawQuery is (for now) simple queries
type RawQuery func(object string, opts ...RawQueryOption) (*types.QueryResult, error)

func newRawQueryFunc(base Transport) RawQuery {
	return func(query string, opts ...RawQueryOption) (*types.QueryResult, error) {
		r := RawQueryRequest{
			Query: query,
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

type RawQueryOption func(*RawQueryRequest) error

type RawQueryRequest struct {
	ctx   context.Context
	Query string
}

func (r RawQueryRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	method := http.MethodGet

	path := fmt.Sprintf("/query?q=%s", url.QueryEscape(r.Query))

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
func (f RawQuery) WithContext(v context.Context) RawQueryOption {
	return func(r *RawQueryRequest) error {
		r.ctx = v
		return nil
	}
}
