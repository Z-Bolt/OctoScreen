package octoprint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Axis string

const (
	XAxis Axis = "x"
	YAxis Axis = "y"
	ZAxis Axis = "z"
)

type Printer struct {
	key      string
	endpoint string
	c        *http.Client
}

func NewPrinter(endpoint, key string) *Printer {
	return &Printer{
		endpoint: endpoint,
		key:      key,
		c:        &http.Client{},
	}
}

func (p *Printer) doRequest(method, target string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, joinURL(p.endpoint, target), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", p.key)

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}

	return p.handleResponse(resp)
}

func (c *Printer) handleResponse(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 200 && r.StatusCode <= 209 {
		return body, nil
	}

	fmt.Println(body)
	return nil, nil
}

func joinURL(base, uri string) string {
	u, _ := url.Parse(base)
	u.Path = path.Join(u.Path, uri)
	return u.String()
}
