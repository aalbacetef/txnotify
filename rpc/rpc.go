package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func (client *Client) generateID() int {
	id := rand.Intn(100) // AA: to improve
	return id
}

type ClientOptions struct {
	Endpoint string
	Timeout  time.Duration
}

const defaultTimeout = 30 * time.Second

func NewClient(options ClientOptions) (*Client, error) {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	if options.Endpoint == "" {
		return nil, MissingFieldError{"Endpoint"}
	}

	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		endpoint:   options.Endpoint,
	}, nil
}

type MissingFieldError struct {
	name string
}

func (e MissingFieldError) Error() string {
	return fmt.Sprintf("missing field: %s", e.name)
}

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e JSONRPCError) Error() string {
	return fmt.Sprintf(
		"json-rpc error: code=%d, message='%s', data=%v",
		e.Code, e.Message, e.Data,
	)
}

type Response[T any] struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  T            `json:"result,omitempty"`
	Error   JSONRPCError `json:"error,omitempty"`
	ID      int          `json:"id"`
}

func Do[T any](client *Client, method string, params []any) (*Response[T], error) {
	id := client.generateID()

	req := Request{
		ID:      id,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	reqBody := &bytes.Buffer{}

	if err := json.NewEncoder(reqBody).Encode(req); err != nil {
		return nil, fmt.Errorf("could not encode body: %w", err)
	}

	resp, err := client.httpClient.Post(
		client.endpoint,
		"application/json",
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("post failed: %w", err)
	}

	defer resp.Body.Close()

	respBody := Response[T]{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	if respBody.Error.Code != 0 {
		return nil, respBody.Error
	}

	if respBody.ID != id {
		return nil, fmt.Errorf("id mismatch: got %d, want %d", respBody.ID, id)
	}

	return &respBody, nil
}
