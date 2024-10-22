package goforce

import (
	"net/http"
	"strings"
)

type Option func(client *Client) error

func WithUrl(url string) Option {
	return func(client *Client) error {
		// Remove trailing "/" from base url to prevent "//" when paths are appended
		client.BaseURL = strings.TrimSuffix(url, "/")
		return nil
	}
}

func WithApiVersion(version string) Option {
	return func(client *Client) error {
		client.ApiVersion = version
		return nil
	}
}

func WithClientId(client_id string) Option {
	return func(client *Client) error {
		client.ClientID = client_id
		return nil
	}
}

func WithHttpClient(hc *http.Client) Option {
	return func(client *Client) error {
		client.HttpClient = hc
		return nil
	}
}

func WithAuthRetry() Option {
	return func(client *Client) error {
		client.AuthRetry = true
		return nil
	}
}
