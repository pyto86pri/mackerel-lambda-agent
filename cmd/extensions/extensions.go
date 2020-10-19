package extensions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Client ...
type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	ID         string
}

// NewClient ...
func NewClient() (*Client, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s", os.Getenv("AWS_LAMBDA_RUNTIME_API")))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return &Client{
		BaseURL:    u,
		HTTPClient: client,
	}, nil
}

func (c *Client) urlFor(path string) *url.URL {
	u, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("Invalid url")
	}
	u.Path = path
	return u
}

func (c *Client) request(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		return nil, fmt.Errorf("Internal server error")
	}
	if resp.StatusCode == 400 {
		return nil, fmt.Errorf("Badrequest")
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("Forbidden")
	}
	return resp, nil
}

// RegisterRequest ...
type RegisterRequest struct {
	Events []string `json:"events"`
}

// RegisterResponse ...
type RegisterResponse struct {
	FunctionName    string `json:"functionName"`
	FunctionVersion string `json:"functionVersion"`
	Handler         string `json:"handler"`
}

// Register ...
func (c *Client) Register() (*RegisterResponse, error) {
	payload := &RegisterRequest{
		Events: []string{"INVOKE", "SHUTDOWN"},
	}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.urlFor("2020-01-01/extension/register").String(), &body)
	if err != nil {
		return nil, err
	}
	filename := filepath.Base(os.Args[0])
	req.Header.Set("Lambda-Extension-Name", filename)

	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	id := resp.Header.Get("Lambda-Extension-Identifier")
	if len(id) == 0 {
		return nil, fmt.Errorf("Did not get identifier")
	}
	c.ID = id

	var rr RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return &rr, nil
}

func (c *Client) withIDHeader(req *http.Request) *http.Request {
	if c.ID == "" {
		panic("Cannot find identifier")
	}
	req.Header.Set("Lambda-Extension-Identifier", c.ID)
	return req
}

func (c *Client) requestWithIDHeader(req *http.Request) (*http.Response, error) {
	return c.request(c.withIDHeader(req))
}

// Event ...
type Event struct {
	EventType  string `json:"eventType"`
	DeadlineMs int    `json:"deadlineMs"`
}

type xRayTracingInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// EventInvoke ...
type EventInvoke struct {
	Event
	RequestID          string          `json:"requestId"`
	InvokedFunctionArn string          `json:"invokedFunctionArn"`
	Tracing            xRayTracingInfo `json:"tracing"`
}

// EventShutdown ...
type EventShutdown struct {
	Event
	ShutdownReason string `json:"shutdownReason"`
}

// Next ...
func (c *Client) Next() (*Event, error) {
	req, err := http.NewRequest("GET", c.urlFor("2020-01-01/extension/event/next").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.requestWithIDHeader(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// TODO: parse event based on type
	var event Event
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// ErrorRequest ...
type ErrorRequest struct {
	ErrorType    string   `json:"errorType"`
	ErrorMessage string   `json:"errorMessage"`
	StackTrace   []string `json:"stackTrace"`
}

// ErrorResponse ...
type ErrorResponse struct {
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
}

// InitError ...
func (c *Client) InitError(payload *ErrorRequest) (*ErrorResponse, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.urlFor("2020-01-01/extension/init/error").String(), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Lambda-Extension-Function-Error-Type", payload.ErrorType)

	resp, err := c.requestWithIDHeader(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var er ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&er)
	if err != nil {
		return nil, err
	}

	return &er, nil
}

// ExitError ...
func (c *Client) ExitError(payload *ErrorRequest) (*ErrorResponse, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.urlFor("2020-01-01/extension/exit/error").String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Lambda-Extension-Function-Error-Type", payload.ErrorType)

	resp, err := c.requestWithIDHeader(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var er ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&er)
	if err != nil {
		return nil, err
	}

	return &er, nil
}
