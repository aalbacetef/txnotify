package txnotify

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/aalbacetef/txnotify/rpc"
)

type nopWriter struct{}

func (nopWriter) Write(v []byte) (int, error) {
	return len(v), nil
}

func TestWatcher(t *testing.T) { //nolint:gocognit
	mock := mustCreateMockClient(t)

	t.Run("it updates the latest block", func(tt *testing.T) {
		watcher := mustMakeWatcher(t, mock)

		if watcher.latestBlock != "" {
			tt.Fatalf("latestBlock should not be set")
		}

		watcher.checkNewBlock()

		if watcher.latestBlock != mock.blockNum {
			tt.Fatalf("got %s, want %s", watcher.latestBlock, mock.blockNum)
		}

		if watcher.currentBlock != "" {
			tt.Fatalf("watcher.currentBlock should not be updated")
		}
	})

	t.Run("it processes the next block", func(tt *testing.T) {
		watcher := mustMakeWatcher(t, mock)

		watcher.checkNewBlock()
		watcher.processNextBlock()

		state := watcher.copyState()

		if state.currentBlock != state.latestBlock {
			tt.Fatalf(
				"currentBlock != latestBlock, cb=%s, lb=%s",
				state.currentBlock, state.latestBlock,
			)
		}

		block, err := watcher.cache.GetBlock(watcher.latestBlock)
		if err != nil {
			tt.Fatalf("error: %v", err)
		}

		n := len(block.Transactions)
		if n != 1 {
			tt.Fatalf("should have only one tx, got %d", n)
		}

		wantTx := block.Transactions[0]
		tx, err := watcher.cache.GetTx(wantTx.Hash)
		if err != nil {
			tt.Fatalf("error: %v", err)
		}

		if tx.Hash != wantTx.Hash {
			tt.Fatalf("tx hash mismatch: got %s, want %s", tx.Hash, wantTx.Hash)
		}
	})
}

func mustMakeWatcher(t *testing.T, mock RPCClient) *Watcher {
	t.Helper()

	watcher := &Watcher{
		pollInterval: 5 * time.Second,
		rpcClient:    mock,
		logger:       slog.New(slog.NewTextHandler(nopWriter{}, nil)),
		notifier:     mockNotifier{},
		cache:        NewInMemoryCache(),
	}

	return watcher
}

type mockNotifier struct{}

func (mockNotifier) Notify(address string, txList []ethereum.Transaction) {
	for _, tx := range txList {
		fmt.Printf("%s) got tx: %s\n", address, tx.Hash)
	}
}

type mockRPCClient struct {
	blockNum  string
	index     string
	blockInfo *rpc.Response[ethereum.Block]
}

//go:embed rpc/testdata/block.0x154d535.json
var testBlockInfoFile []byte

func mustCreateMockClient(t *testing.T) *mockRPCClient {
	t.Helper()

	mock := &mockRPCClient{
		blockNum:  "0x154d535",
		index:     "0x1",
		blockInfo: &rpc.Response[ethereum.Block]{},
	}

	if err := json.NewDecoder(bytes.NewReader(testBlockInfoFile)).Decode(&mock.blockInfo); err != nil {
		t.Fatalf("error decoding block info file: %v", err)
	}

	return mock
}

func (m *mockRPCClient) GetBlockByNumber(blockNum string) (*rpc.Response[ethereum.Block], error) {
	if blockNum != m.blockNum {
		return nil, fmt.Errorf("mock expects block number %s, got %s", m.blockNum, blockNum)
	}

	return m.blockInfo, nil
}

func (m *mockRPCClient) GetCurrentBlockNumber() (*rpc.Response[string], error) {
	return &rpc.Response[string]{
		ID:      1,
		JSONRPC: "2.0",
		Result:  m.blockNum,
	}, nil
}
