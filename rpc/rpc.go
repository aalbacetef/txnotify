package rpc

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
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

type Response[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	Result  T      `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
	ID      int    `json:"id"`
}
