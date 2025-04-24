package rpc

import (
	"encoding/json"
	"testing"
)

const testEndpoint = "https://eth.nodeconnect.org"

func TestCurrentBlock(t *testing.T) {
	client := mustMakeClient(t, testEndpoint)

	response, err := client.GetCurrentBlock()
	if err != nil {
		t.Fatalf("could not get block: %v", err)
	}

	t.Logf("block num: %s", response.Result)
}

func TestGetBlockInfo(t *testing.T) {
	client := mustMakeClient(t, testEndpoint)

	v := 0x154d535

	m, err := client.GetBlockByNumber(v)
	if err != nil {
		t.Fatalf("could not fetch block: %v", err)
	}

	d, _ := json.MarshalIndent(m, "", "  ")
	t.Logf("%s", string(d))
}

func mustMakeClient(t *testing.T, endpoint string) *Client {
	t.Helper()

	client, err := NewClient(ClientOptions{Endpoint: endpoint})
	if err != nil {
		t.Fatalf("could not start client: %v", err)
	}

	return client
}
