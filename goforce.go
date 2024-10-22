package goforce

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/0xArch3r/goforce/api"
	"github.com/0xArch3r/goforce/types"
)

const (
	DefaultAPIVersion = "54.0"
	DefaultClientID   = "goforce"
	DefaultURL        = "https://login.salesforce.com"
)

// Client is the main instance to access salesforce.
type BaseClient struct {
	SessionID string
	User      struct {
		Id       string
		Name     string
		FullName string
		Email    string
	}
	Password      string
	ClientID      string
	ApiVersion    string
	BaseURL       string
	InstanceURL   string
	UseToolingAPI bool
	HttpClient    *http.Client
	AuthRetry     bool
}

type Client struct {
	BaseClient
	*api.Api
}

// NewClient creates a new instance of the client.
func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		BaseClient: BaseClient{
			ApiVersion: DefaultAPIVersion,
			BaseURL:    DefaultURL,
			ClientID:   DefaultClientID,
			HttpClient: http.DefaultClient,
		},
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	client.Api = api.New(client)

	return client, nil
}

// Perform delegates to Transport to execute a request and return a response.
func (c *BaseClient) Perform(req *http.Request) (*api.Response, error) {
	original_path := req.URL.Path
	query := req.URL.RawQuery
	u := fmt.Sprintf("%v/services/data/v%v%v?%v", c.InstanceURL, c.ApiVersion, original_path, query)

	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	req.URL = url
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.SessionID))
	req.Header.Add("Content-Type", "application/json")

	// Retrieve the original request.
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &api.Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}, nil

}

// LoginPassword signs into salesforce using password. token is optional if trusted IP is configured.
// Ref: https://developer.salesforce.com/docs/atlas.en-us.214.0.api_rest.meta/api_rest/intro_understanding_username_password_oauth_flow.htm
// Ref: https://developer.salesforce.com/docs/atlas.en-us.214.0.api.meta/api/sforce_api_calls_login.htm
func (client *BaseClient) LoginPassword(username, password, token string) error {
	// Use the SOAP interface to acquire session ID with username, password, and token.
	// Do not use REST interface here as REST interface seems to have strong checking against client_id, while the SOAP
	// interface allows a non-exist placeholder client_id to be used.
	soapBody := `<?xml version="1.0" encoding="utf-8" ?>
        <env:Envelope
                xmlns:xsd="http://www.w3.org/2001/XMLSchema"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                xmlns:env="http://schemas.xmlsoap.org/soap/envelope/"
                xmlns:urn="urn:partner.soap.sforce.com">
            <env:Header>
                <urn:CallOptions>
                    <urn:client>%s</urn:client>
                    <urn:defaultNamespace>sf</urn:defaultNamespace>
                </urn:CallOptions>
            </env:Header>
            <env:Body>
                <n1:login xmlns:n1="urn:partner.soap.sforce.com">
                    <n1:username>%s</n1:username>
                    <n1:password>%s%s</n1:password>
                </n1:login>
            </env:Body>
        </env:Envelope>`
	soapBody = fmt.Sprintf(soapBody, client.ClientID, username, html.EscapeString(password), token)

	url := fmt.Sprintf("%s/services/Soap/u/%s", client.BaseURL, client.ApiVersion)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(soapBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("charset", "UTF-8")
	req.Header.Add("SOAPAction", "login")

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		theError := types.ParseSalesforceError(resp.StatusCode, buf.Bytes())
		return theError
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResponse struct {
		XMLName      xml.Name `xml:"Envelope"`
		ServerURL    string   `xml:"Body>loginResponse>result>serverUrl"`
		SessionID    string   `xml:"Body>loginResponse>result>sessionId"`
		UserID       string   `xml:"Body>loginResponse>result>userId"`
		UserEmail    string   `xml:"Body>loginResponse>result>userInfo>userEmail"`
		UserFullName string   `xml:"Body>loginResponse>result>userInfo>userFullName"`
		UserName     string   `xml:"Body>loginResponse>result>userInfo>userName"`
	}

	err = xml.Unmarshal(respData, &loginResponse)
	if err != nil {
		return err
	}

	// Now we should all be good and the sessionID can be used to talk to salesforce further.
	client.SessionID = loginResponse.SessionID
	client.InstanceURL = parseHost(loginResponse.ServerURL)
	client.User.Id = loginResponse.UserID
	client.User.Name = loginResponse.UserName
	client.User.Email = loginResponse.UserEmail
	client.User.FullName = loginResponse.UserFullName
	client.Password = password

	return nil
}

func parseHost(input string) string {
	parsed, err := url.Parse(input)
	if err == nil {
		return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	}
	return "Failed to parse URL input"
}
