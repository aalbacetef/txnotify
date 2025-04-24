package main

import (
	"encoding/json"
	"testing"
)

func TestRPCCurrentBlock(t *testing.T) {
	client, err := NewClient(RPCClientOptions{Endpoint: "https://eth.nodeconnect.org"})
	if err != nil {
		t.Fatalf("could not start client: %v", err)
	}

	v, err := client.GetCurrentBlock()
	if err != nil {
		t.Fatalf("could not get block: %v", err)
	}

	t.Logf("block num: %#0x", v)

	m, err := client.GetBlockByNumber(v)
	if err != nil {
		t.Fatalf("could not fetch block: %v", err)
	}

	d, _ := json.MarshalIndent(m, "", "  ")

	t.Logf("%s", string(d))

}
