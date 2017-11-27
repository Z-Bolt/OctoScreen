package octoprint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// A Client manages communication with the OctoPrint API.
type Client struct {
	apiKey  string
	baseURL string
	c       *http.Client
}

// NewClient returns a new OctoPrint API client with provided base URL and API
// Key. If baseURL does not have a trailing slash, one is added automatically. If
// `Access Control` is enabled at OctoPrint configuration an apiKey should be
// provided (http://docs.octoprint.org/en/master/api/general.html#authorization).
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		c:       &http.Client{},
	}
}

func (c *Client) doRequest(method, target string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, joinURL(c.baseURL, target), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", c.apiKey)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	js, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	err = cacheRequest(target, js)
	return js, err
}

func (c *Client) handleResponse(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
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

func cacheRequest(uri string, js []byte) error {
	u, _ := url.Parse(uri)
	path := filepath.Join("cache", strings.Replace(u.Path, "/", "-", -1), time.Now().String()+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		panic(err)
	}

	return ioutil.WriteFile(path, js, 0777)
}
