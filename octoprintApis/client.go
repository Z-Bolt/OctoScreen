package octoprintApis

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	// "log"
	"net/http"
	"net/url"
	"time"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ErrUnauthorized missing or invalid API key
var ErrUnauthorized = errors.New("Missing or invalid API key")


// A Client manages communication with the OctoPrint API.
type Client struct {
	// Endpoint address to the OctoPrint REST API server.
	Endpoint string

	// APIKey used to connect to the OctoPrint REST API server.
	APIKey string

	// HTTP client connection.
	httpClient *http.Client
}

// NewClient returns a new OctoPrint API client with provided base URL and API
// Key. If baseURL does not have a trailing slash, one is added automatically. If
// `Access Control` is enabled at OctoPrint configuration an apiKey should be
// provided (http://docs.octoprint.org/en/master/api/general.html#authorization).
func NewClient(endpoint, apiKey string) *Client {
	return &Client {
		Endpoint: endpoint,
		APIKey:   apiKey,
		httpClient: &http.Client {
			Timeout: time.Second * 3,
			Transport: &http.Transport {
				DisableKeepAlives: true,
			},
		},
	}
}

func (this *Client) doJsonRequest(
	method string,
	target string,
	body io.Reader,
	statusMapping StatusMapping,
) ([]byte, error) {
	LogMessagef("    entering Client.doJsonRequest()")

	bytes, err := this.doRequest(method, target, "application/json", body, statusMapping)
	if err != nil {
		LogError(err, "Client.doJsonRequest(), this.doRequest() failed")
		LogMessagef("    leaving Client.doJsonRequest()")
		return nil, err
	}

	// Use the following for debugging.  Comment out for production.
	json := string(bytes)
	LogMessagef("        JSON response: %s", json)
	LogMessagef("    leaving Client.doJsonRequest()")

	return bytes, err
}

func (this *Client) doRequest(
	method string,
	target string,
	contentType string,
	body io.Reader,
	statusMapping StatusMapping,
) ([]byte, error) {
	LogMessagef("        entering Client.doRequest()")
	LogMessagef("            method: %s", method)
	LogMessagef("            target: %s",target)


	req, err := http.NewRequest(method, joinUrl(this.Endpoint, target), body)
	if err != nil {
		LogError(err, "Client.doRequest(), http.NewRequest() failed")
		LogMessagef("        leaving Client.doRequest()")
		return nil, err
	}

	req.Header.Add("Host", "localhost:5000")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", fmt.Sprintf("go-octoprint/%s", Version))
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	req.Header.Add("X-Api-Key", this.APIKey)
	resp, err := this.httpClient.Do(req)
	if err != nil {
		LogError(err, "Client.doRequest(), this.httpClient.Do() failed")
		LogMessagef("        leaving Client.doRequest()")
		return nil, err
	}

	response, err := this.handleResponse(resp, statusMapping)
	if err != nil {
		LogError(err, "Client.doRequest(), this.handleResponse() failed")
		LogMessagef("        leaving Client.doRequest()")
		return nil, err
	}

	LogMessagef("        leaving Client.doRequest()")
	return response, err
}

func (this *Client) handleResponse(
	httpResponse *http.Response,
	statusMapping StatusMapping,
) ([]byte, error) {
	defer httpResponse.Body.Close()

	if statusMapping != nil {
		if err := statusMapping.Error(httpResponse.StatusCode); err != nil {
			return nil, err
		}
	}

	if httpResponse.StatusCode == 401 {
		return nil, ErrUnauthorized
	}

	if httpResponse.StatusCode == 204 {
		return nil, nil
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode <= 209 {
		return body, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
}


func joinUrl(base, uri string) string {
	u, _ := url.Parse(uri)
	b, _ := url.Parse(base)
	return b.ResolveReference(u).String()
}
