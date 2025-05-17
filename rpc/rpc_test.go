package rpc

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/aalbacetef/txnotify/ethereum"
)

const (
	testEndpoint    = "https://eth.nodeconnect.org"
	testBlockNumber = "0x154d535"
)

//go:embed testdata/block.0x154d535.json
var testBlockInfo []byte

func TestCurrentBlock(t *testing.T) {
	client := mustMakeClient(t, testEndpoint)

	response, err := client.GetCurrentBlockNumber()
	if err != nil {
		t.Fatalf("could not get block: %v", err)
	}

	if !strings.HasPrefix(response.Result, "0x") {
		t.Fatalf("expected a block number, got %s", response.Result)
	}
}

func TestGetBlockInfo(t *testing.T) {
	client := mustMakeClient(t, testEndpoint)

	response, err := client.GetBlockByNumber(testBlockNumber)
	if err != nil {
		t.Fatalf("could not fetch block: %v", err)
	}

	savedResponse := Response[ethereum.Block]{}

	if err := json.NewDecoder(bytes.NewReader(testBlockInfo)).Decode(&savedResponse); err != nil {
		t.Fatalf("could not decode saved response: %v", err)
	}

	if response.Result.Hash != savedResponse.Result.Hash {
		t.Fatalf("(hash) got %s, want %s", response.Result.Hash, savedResponse.Result.Hash)
	}

	// block will have many transactions, lets trim all but one
	savedTx := savedResponse.Result.Transactions[0]

	for _, got := range response.Result.Transactions {
		if got.Hash == savedTx.Hash {
			return
		}
	}

	t.Fatalf("transaction with hash %s not found", savedTx.Hash)
}

func TestInvalidMethod(t *testing.T) {
	client := mustMakeClient(t, testEndpoint)

	invalidEndpoint := "eth_doesnt_exist"

	wantErr := JSONRPCError{
		Code:    -32601,
		Message: "the method eth_doesnt_exist does not exist/is not available",
		Data:    nil,
	}

	_, err := Do[string](client, invalidEndpoint, []any{})
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("got: '%v', want: '%v'", err, wantErr)
	}
}

func mustMakeClient(t *testing.T, endpoint string) *Client { //nolint:unparam
	t.Helper()

	client, err := NewClient(ClientOptions{Endpoint: endpoint})
	if err != nil {
		t.Fatalf("could not start client: %v", err)
	}

	return client
}

func compareTransaction(t *testing.T, gotTx, wantTx ethereum.Transaction) {
	t.Helper()

	if gotTx.Hash != wantTx.Hash {
		t.Errorf("(hash) got %s, want %s", gotTx.Hash, wantTx.Hash)
	}
	if gotTx.From != wantTx.From {
		t.Errorf("(from) got %s, want %s", gotTx.From, wantTx.From)
	}
	if gotTx.Gas != wantTx.Gas {
		t.Errorf("(gas) got %s, want %s", gotTx.Gas, wantTx.Gas)
	}
	if gotTx.GasPrice != wantTx.GasPrice {
		t.Errorf("(gasPrice) got %s, want %s", gotTx.GasPrice, wantTx.GasPrice)
	}
	if gotTx.Input != wantTx.Input {
		t.Errorf("(input) got %s, want %s", gotTx.Input, wantTx.Input)
	}
	if gotTx.Nonce != wantTx.Nonce {
		t.Errorf("(nonce) got %s, want %s", gotTx.Nonce, wantTx.Nonce)
	}
	if gotTx.R != wantTx.R {
		t.Errorf("(r) got %s, want %s", gotTx.R, wantTx.R)
	}
	if gotTx.S != wantTx.S {
		t.Errorf("(s) got %s, want %s", gotTx.S, wantTx.S)
	}
	if gotTx.Type != wantTx.Type {
		t.Errorf("(type) got %s, want %s", gotTx.Type, wantTx.Type)
	}
	if gotTx.V != wantTx.V {
		t.Errorf("(v) got %s, want %s", gotTx.V, wantTx.V)
	}
	if gotTx.Value != wantTx.Value {
		t.Errorf("(value) got %s, want %s", gotTx.Value, wantTx.Value)
	}

	compareOptionalString(t, "blockHash", gotTx.BlockHash, wantTx.BlockHash)
	compareOptionalString(t, "blockNumber", gotTx.BlockNumber, wantTx.BlockNumber)
	compareOptionalString(t, "chainId", gotTx.ChainID, wantTx.ChainID)
	compareOptionalString(t, "to", gotTx.To, wantTx.To)
	compareOptionalString(t, "transactionIndex", gotTx.TransactionIndex, wantTx.TransactionIndex)
	compareOptionalString(t, "maxPriorityFeePerGas", gotTx.MaxPriorityFeePerGas, wantTx.MaxPriorityFeePerGas)
	compareOptionalString(t, "maxFeePerGas", gotTx.MaxFeePerGas, wantTx.MaxFeePerGas)
	compareOptionalString(t, "yParity", gotTx.YParity, wantTx.YParity)
}

// compareOptionalString compares two optional string pointers and reports an error if they differ.
func compareOptionalString(t *testing.T, fieldName string, got, want *string) {
	t.Helper()
	if got == nil && want == nil {
		return
	}

	if got == nil || want == nil || *got != *want {
		t.Errorf("(%s) got %v, want %v", fieldName, got, want)
	}
}
