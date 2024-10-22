package api

import "net/http"

type Api struct {
	Get    Get
	Search Search
}

func New(base Transport) *Api {
	return &Api{
		Get:    newGetFunc(base),
		Search: newSearchFunc(base),
	}
}

type Transport interface {
	Perform(*http.Request) (*Response, error)
}
