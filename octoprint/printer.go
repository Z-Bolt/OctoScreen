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

	js, err := p.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	err = cacheRequest(target, js)
	return js, err
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
