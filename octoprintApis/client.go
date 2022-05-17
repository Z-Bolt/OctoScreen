package octoprintApis

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Z-Bolt/OctoScreen/logger"
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
// 'Access Control' is enabled at OctoPrint configuration an apiKey should be
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
	isRequired bool,
) ([]byte, error) {
	logger.TraceEnter("Client.doJsonRequest()")

	bytes, err := this.doRequest(method, target, "application/json", body, statusMapping, isRequired)
	if err != nil {
		logOptionalError("Client.doJsonRequest()", "this.doRequest()", err, isRequired)
		logger.TraceLeave("Client.doJsonRequest()")
		return nil, err
	}

	// Use the following only for debugging.
	if logger.LogLevel() == "debug" {
		logger.Debug("Client.doJsonRequest() - converting bytes to JSON")
		json := string(bytes)
		logger.Debugf("JSON response: %s", json)
	}

	logger.TraceLeave("Client.doJsonRequest()")
	return bytes, err
}

func (this *Client) doRequest(
	method string,
	target string,
	contentType string,
	body io.Reader,
	statusMapping StatusMapping,
	isRequired bool,
) ([]byte, error) {
	logger.TraceEnter("Client.doRequest()")
	logger.Debugf("method: %s", method)
	logger.Debugf("target: %s", target)
	logger.Debugf("contentType: %s", contentType)

	url := joinUrl(this.Endpoint, target)
	logger.Debugf("url: %s", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.LogError("Client.doRequest()", "http.NewRequest()", err)
		logger.TraceLeave("Client.doRequest()")
		return nil, err
	}

	req.Header.Add("Host", "localhost:5000")
	req.Header.Add("Accept", "*/*")

	userAgent := fmt.Sprintf("go-octoprint/%s", Version)
	logger.Debugf("userAgent: %s", userAgent)
	req.Header.Add("User-Agent", userAgent)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	// Don't log APIKey due to privacy & security.
	// logger.Debugf("API key: %s", this.APIKey)
	req.Header.Add("X-Api-Key", this.APIKey)

	response, err := this.httpClient.Do(req)
	if err != nil {
		logger.LogError("Client.doRequest()", "this.httpClient.Do()", err)
		logger.TraceLeave("Client.doRequest()")
		return nil, err
	} else {
		logger.Debug("Client.doRequest() - httpClient.Do() passed")
	}

	bytes, err := this.handleResponse(response, statusMapping, isRequired)
	if err != nil {
		logOptionalError("Client.doRequest()", "this.handleResponse()", err, isRequired)
		bytes = nil
	} else {
		logger.Debug("Client.doRequest() - handleResponse() passed")
	}

	logger.TraceLeave("Client.doRequest()")
	return bytes, err
}

func (this *Client) handleResponse(
	httpResponse *http.Response,
	statusMapping StatusMapping,
	isRequired bool,
) ([]byte, error) {
	logger.TraceEnter("Client.handleResponse()")

	defer func() {
		io.Copy(ioutil.Discard, httpResponse.Body)
		httpResponse.Body.Close()
	}()

	if statusMapping != nil {
		if err := statusMapping.Error(httpResponse.StatusCode); err != nil {
			logger.LogError("Client.handleResponse()", "statusMapping.Error()", err)
			logger.TraceLeave("Client.handleResponse()")
			return nil, err
		}
	}

	if httpResponse.StatusCode == 401 {
		logger.Error("Client.handleResponse() - StatusCode is 401")
		logger.TraceLeave("Client.handleResponse()")
		return nil, ErrUnauthorized
	}

	if httpResponse.StatusCode == 204 {
		logger.Error("Client.handleResponse() - StatusCode is 204")
		logger.TraceLeave("Client.handleResponse()")
		return nil, nil
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		logger.LogError("Client.handleResponse()", "ioutil.ReadAll()", err)
		logger.TraceLeave("Client.handleResponse()")
		return nil, err
	}

	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode <= 209 {
		logger.Debugf("Client.handleResponse() - status code %d was within range", httpResponse.StatusCode)
	} else {
		errMsg := fmt.Sprintf("Unexpected status code: %d", httpResponse.StatusCode)
		if httpResponse.StatusCode == 404 {
			logOptionalMessage(errMsg, isRequired)
		} else {
			logger.Error(errMsg)
		}

		err = fmt.Errorf(errMsg)
		body = nil
	}

	logger.TraceLeave("Client.handleResponse()")
	return body, err
}

func logOptionalError(
	currentFunctionName string,
	functionCalledName string,
	err error,
	isRequired bool,
) {
	if isRequired {
		// Some APIs return an error and the error should be logged.
		logger.LogError(currentFunctionName, functionCalledName, err)
	} else {
		// On the other hand, calls to some APIs are optional, and the result should be logged
		// as info and leave it up to the caller to determine whether it's an error or not.
		msg := fmt.Sprintf("%s - %s returned %q", currentFunctionName, functionCalledName, err)
		logger.Info(msg)
	}
}

func logOptionalMessage(
	msg string,
	isRequired bool,
) {
	if isRequired {
		logger.Error(msg)
	} else {
		logger.Info(msg)
	}
}

func joinUrl(base, uri string) string {
	u, _ := url.Parse(uri)
	b, _ := url.Parse(base)
	return b.ResolveReference(u).String()
}
