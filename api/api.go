package api

import "net/http"

type Api struct {
	Get    Get
	Search Search
	Query  *Query
}

type Query struct {
	Select Select
	Raw    RawQuery
}

func New(base Transport) *Api {
	return &Api{
		Get:    newGetFunc(base),
		Search: newSearchFunc(base),
		Query: &Query{
			Select: newSelectFunc(base),
			Raw:    newRawQueryFunc(base),
		},
	}
}

type Transport interface {
	Perform(*http.Request) (*Response, error)
}
