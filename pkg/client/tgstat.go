// Code generated from jsonrpc schema by rpcgen v2.4.6; DO NOT EDIT.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// Always import time package. Generated models can contain time.Time fields.
	_ time.Time
)

type Client struct {
	httpCLient *httpClient
	Services   *svcServices
}

func NewDefaultClient(endpoint string) *Client {
	return NewClient(endpoint, http.Header{}, &http.Client{})
}

func NewClient(endpoint string, header http.Header, httpClient *http.Client) *Client {
	c := &Client{
		httpCLient: newHTTPClient(endpoint, header, httpClient),
	}

	c.Services = newClientServices(c.httpCLient)
	return c
}

type Item struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type ServiceListResponse struct {
	Services []Item `json:"services"`
}

type svcServices struct {
	clientHttp *httpClient
}

func newClientServices(clienthttp *httpClient) *svcServices {
	return &svcServices{
		clientHttp: clienthttp,
	}
}

type httpClient struct {
	endpoint string
	cl       *http.Client

	requestID uint64
	header    http.Header
}

func newHTTPClient(endpoint string, header http.Header, httpclient *http.Client) *httpClient {
	return &httpClient{
		endpoint: endpoint,
		header:   header,
		cl:       httpclient,
	}
}

func (c *svcServices) GetAllServices(ctx context.Context) (res []Item, err error) {

	req, err := http.NewRequest("GET", c.clientHttp.endpoint, nil)

	if err != nil {
		return nil, fmt.Errorf("failed create new request. err : %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		return nil, fmt.Errorf("failed to do request err : %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result ServiceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Services, nil
}
