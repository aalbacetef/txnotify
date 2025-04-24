package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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

type GetCurrentBlockRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type GetCurrentBlockResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

func (client *Client) GetCurrentBlock() (int, error) {
	id := rand.Intn(100) // AA: to improve

	payload := GetCurrentBlockRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []any{},
		ID:      id,
	}

	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return 0, fmt.Errorf("could not encode body: %w", err)
	}

	resp, err := client.httpClient.Post(
		client.endpoint,
		"application/json",
		body,
	)
	if err != nil {
		return 0, fmt.Errorf("post failed: %w", err)
	}

	defer resp.Body.Close()

	respBody := GetCurrentBlockResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return 0, fmt.Errorf("could not decode response body: %w", err)
	}

	if respBody.ID != id {
		return 0, fmt.Errorf("id mismatch: got %d, want %d", respBody.ID, id)
	}

	v, err := strconv.ParseInt(respBody.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse hex value: %w", err)
	}

	return int(v), nil
}

type GetBlockByNumRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

// We only care about a few fields
type GetBlockByNumResponse struct {
	JSONRPC string `json:"jsonrpc"`

	Result struct {
		Hash         string   `json:"hash"`
		Transactions []string `json:"transactions"`
	} `json:"result"`

	ID int `json:"id"`
}

func (client *Client) GetBlockByNumber(blockNum int) (*GetBlockByNumResponse, error) {
	id := rand.Intn(100) // AA: to improve

	payload := GetBlockByNumRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []any{
			fmt.Sprintf("%#0x", blockNum),
			false,
		},
		ID: id,
	}

	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, fmt.Errorf("could not encode body: %w", err)
	}

	resp, err := client.httpClient.Post(
		client.endpoint,
		"application/json",
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("post failed: %w", err)
	}

	defer resp.Body.Close()

	respBody := GetBlockByNumResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	return &respBody, nil
}
