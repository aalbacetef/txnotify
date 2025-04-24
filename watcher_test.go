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

func TestWatcher(t *testing.T) {
	mock := mustCreateMockClient(t)

	t.Run("it updates the latest block", func(tt *testing.T) {
		watcher := &Watcher{
			pollInterval: 5 * time.Second,
			rpcClient:    mock,
			logger:       slog.New(slog.NewTextHandler(nopWriter{}, nil)),
			cache:        NewInMemoryCache(),
		}

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
		watcher := &Watcher{
			pollInterval: 5 * time.Second,
			rpcClient:    mock,
			logger:       slog.New(slog.NewTextHandler(nopWriter{}, nil)),
			notifier:     mockNotifier{},
			cache:        NewInMemoryCache(),
		}

		watcher.checkNewBlock()
		watcher.processNextBlock()
	})
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
	txCount   *rpc.Response[string]
	txData    *rpc.Response[ethereum.Transaction]
	blockInfo *rpc.Response[ethereum.Block]
}

//go:embed rpc/testdata/tx-count.0x154d535.json
var testTxCountFile []byte

//go:embed rpc/testdata/tx.0x154d535.0x1.json
var testTxDataFile []byte

//go:embed rpc/testdata/block.0x154d535.json
var testBlockInfoFile []byte

func mustCreateMockClient(t *testing.T) *mockRPCClient {
	t.Helper()

	mock := &mockRPCClient{
		blockNum:  "0x154d535",
		index:     "0x1",
		txCount:   &rpc.Response[string]{},
		txData:    &rpc.Response[ethereum.Transaction]{},
		blockInfo: &rpc.Response[ethereum.Block]{},
	}

	if err := json.NewDecoder(bytes.NewReader(testTxCountFile)).Decode(mock.txCount); err != nil {
		t.Fatalf("error decoding tx count file: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(testTxDataFile)).Decode(mock.txData); err != nil {
		t.Fatalf("error decoding tx data file: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(testBlockInfoFile)).Decode(mock.blockInfo); err != nil {
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

func (m *mockRPCClient) GetTransactionByHash(hash string) (*rpc.Response[ethereum.Transaction], error) {
	if hash != m.txData.Result.Hash {
		return nil, fmt.Errorf("mock expects tx hash %s, got %s", m.txData.Result.Hash, hash)
	}

	return m.txData, nil
}
