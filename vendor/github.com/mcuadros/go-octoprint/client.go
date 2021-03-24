package octoprint

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	method, target string, body io.Reader, m StatusMapping,
) ([]byte, error) {
	return c.doRequest(method, target, "application/json", body, m)
}

func (c *Client) doJSONRequestWithLogging(
	method, target string, body io.Reader, m StatusMapping,
) ([]byte, error) {
	return c.doRequestWithLogging(method, target, "application/json", body, m)
}



func (c *Client) doRequest(
	method, target, contentType string, body io.Reader, m StatusMapping,
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




func (c *Client) doRequestWithLogging(
	method, target, contentType string, body io.Reader, m StatusMapping,
) ([]byte, error) {


	log.Println("Now in Client.doRequest()")



	req, err := http.NewRequest(method, joinURL(c.Endpoint, target), body)
	if err != nil {
		log.Println("Client.doRequest() - NewRequest() failed")
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
	// resp, err := c.c.DoWithLogging(req)
	if err != nil {
		log.Println("Client.doRequest() - c.c.Do() failed")
		return nil, err
	}





	log.Println("!!! Now in Client.doRequest() - finished calling DoWithLogging() !!!")
	if resp != nil {
		log.Printf("!!! Now in Client.doRequest() - resp.Status: %s", resp.Status)
		log.Printf("!!! Now in Client.doRequest() - resp.StatusCode: %d", resp.StatusCode)
	} else {
		log.Printf("!!! Now in Client.doRequest() - resp was nil")
	}

	if err != nil {
		log.Printf("!!! Now in Client.doRequest() - err: %s", err.Error())
	} else {
		log.Printf("!!! Now in Client.doRequest() - err was nil")
	}
	return nil, err



	// return c.handleResponseWithLogging(resp, m)
}



func (c *Client) handleResponse(r *http.Response, m StatusMapping) ([]byte, error) {
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



func (c *Client) handleResponseWithLogging(r *http.Response, m StatusMapping) ([]byte, error) {

	log.Println("Now in Client.handleResponse()")

	defer r.Body.Close()

	if m != nil {
		log.Println("Client.handleResponse() - m is nil")

		if err := m.Error(r.StatusCode); err != nil {
			log.Println("Client.handleResponse() - m.Error is not nil and is: ", err)
			log.Println("Client.handleResponse() - r.StatusCode: ", r.StatusCode)
			return nil, err
		}
	}

	if r.StatusCode == 401 {
		log.Println("Client.handleResponse() - status code is 401, returning")
		return nil, ErrUnauthorized
	}

	if r.StatusCode == 204 {
		log.Println("Client.handleResponse() - status code is 204, returning")
		return nil, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Client.handleResponse() - ReadAll() failed, err is: ", err)
		return nil, err
	}

	if r.StatusCode >= 200 && r.StatusCode <= 209 {
		log.Println("Client.handleResponse() - status code appears to be good, returning")
		return body, nil
	}
	
	log.Println("Client.handleResponse() - looks like it failed, status code was ", r.StatusCode)


	return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
}





func joinURL(base, uri string) string {
	u, _ := url.Parse(uri)
	b, _ := url.Parse(base)
	return b.ResolveReference(u).String()
}

type StatusMapping map[int]string

func (m *StatusMapping) Error(code int) error {
	err, ok := (*m)[code]
	if ok {
		return fmt.Errorf(err)
	}

	return nil
}
