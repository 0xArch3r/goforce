package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/0xArch3r/goforce/types"
)

func newGetFunc(b Transport) Get {
	return func(object string, id string, o ...GetOption) (*types.SObject, error) {
		r := GetRequest{Object: object, ID: id}
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

		obj := &types.SObject{}
		err = json.Unmarshal(data, obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
}

type Get func(object string, id string, o ...GetOption) (*types.SObject, error)

type GetOption func(*GetRequest) error

// GetRequest configures the Get API request.
type GetRequest struct {
	Object string
	ID     string

	ctx context.Context
}

// Do executes the request and returns response or error.
func (r GetRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   string
		//params map[string]string
	)

	method = http.MethodGet

	path = fmt.Sprintf("/sobjects/%v/%v", r.Object, r.ID)

	//params = make(map[string]string)

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
func (f Get) WithContext(v context.Context) GetOption {
	return func(r *GetRequest) error {
		r.ctx = v
		return nil
	}
}
