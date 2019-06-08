package octoprint

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ErrUnauthorized missing or invalid API key
var ErrUnauthorized = errors.New("Missing or invalid API key")

// A Client manages communication with the OctoPrint API.
type Client struct {
	// Endpoint address to the OctoPrint REST API server.
	Endpoint string
	// APIKey used to connect to the OctoPrint REST API server.
	APIKey string

	c *http.Client
}

// NewClient returns a new OctoPrint API client with provided base URL and API
// Key. If baseURL does not have a trailing slash, one is added automatically. If
// `Access Control` is enabled at OctoPrint configuration an apiKey should be
// provided (http://docs.octoprint.org/en/master/api/general.html#authorization).
func NewClient(endpoint, apiKey string) *Client {
	return &Client{
		Endpoint: endpoint,
		APIKey:   apiKey,
		c: &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		},
	}
}

func (c *Client) doJSONRequest(
	method, target string, body io.Reader, m statusMapping,
) ([]byte, error) {
	return c.doRequest(method, target, "application/json", body, m)
}

func (c *Client) doRequest(
	method, target, contentType string, body io.Reader, m statusMapping,
) ([]byte, error) {
	req, err := http.NewRequest(method, joinURL(c.Endpoint, target), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "localhost:5000")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", fmt.Sprintf("go-octoprint/%s", Version))
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	req.Header.Add("X-Api-Key", c.APIKey)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	return c.handleResponse(resp, m)
}

func (c *Client) handleResponse(r *http.Response, m statusMapping) ([]byte, error) {
	defer r.Body.Close()

	if m != nil {
		if err := m.Error(r.StatusCode); err != nil {
			return nil, err
		}
	}

	if r.StatusCode == 401 {
		return nil, ErrUnauthorized
	}

	if r.StatusCode == 204 {
		return nil, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 200 && r.StatusCode <= 209 {
		return body, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
}

func joinURL(base, uri string) string {
	u, _ := url.Parse(uri)
	b, _ := url.Parse(base)
	return b.ResolveReference(u).String()
}

type statusMapping map[int]string

func (m *statusMapping) Error(code int) error {
	err, ok := (*m)[code]
	if ok {
		return fmt.Errorf(err)
	}

	return nil
}
