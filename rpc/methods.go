package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
)

type GetCurrentBlockRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type GetCurrentBlockResponse Response[string]

func (client *Client) GetCurrentBlock() (*GetCurrentBlockResponse, error) {
	id := rand.Intn(100) // AA: to improve

	payload := GetCurrentBlockRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []any{},
		ID:      id,
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

	respBody := GetCurrentBlockResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	if respBody.ID != id {
		return nil, fmt.Errorf("id mismatch: got %d, want %d", respBody.ID, id)
	}

	return &respBody, nil
}

type GetBlockByNumberRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

// We only care about a few fields
type GetBlockByNumberResponseInner struct {
	Hash         string   `json:"hash"`
	Transactions []string `json:"transactions"`
}

type GetBlockByNumberResponse Response[GetBlockByNumberResponseInner]

func (client *Client) GetBlockByNumber(blockNum int) (*GetBlockByNumberResponse, error) {
	id := rand.Intn(100) // AA: to improve

	payload := GetBlockByNumberRequest{
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

	respBody := GetBlockByNumberResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	return &respBody, nil
}
