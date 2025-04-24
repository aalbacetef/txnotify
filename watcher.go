package txnotify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/aalbacetef/txnotify/rpc"
)

func NewWatcher(rpcEndpoint string, pollInterval time.Duration) (*Watcher, error) {
	client, err := rpc.NewClient(rpc.ClientOptions{
		Endpoint: rpcEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("could not initialize client: %w", err)
	}

	watcher := &Watcher{
		pollInterval: pollInterval,
		rpcClient:    client,
		logger:       slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	return watcher, nil
}

type Watcher struct {
	mu            sync.Mutex
	subscriptions []string
	cancel        context.CancelFunc
	pollInterval  time.Duration
	rpcClient     RPCClient
	currentBlock  string
	latestBlock   string
	logger        *slog.Logger
}

type RPCClient interface {
	GetCurrentBlockNumber() (*rpc.Response[string], error)
	GetTransactionByBlockNumberAndIndex(blockNum, index string) (*rpc.Response[ethereum.Transaction], error)
	GetTransactionCountByNumber(blockNum string) (*rpc.Response[string], error)
}

func (watcher *Watcher) Close() error {
	watcher.mu.Lock()
	defer watcher.mu.Unlock()

	if watcher.cancel != nil {
		watcher.cancel()
		watcher.cancel = nil
	}

	return nil
}

func (watcher *Watcher) Subscribe(address string) error {
	watcher.mu.Lock()
	defer watcher.mu.Unlock()

	watcher.subscriptions = append(watcher.subscriptions, address)

	return nil
}

func (watcher *Watcher) Listen(backgroundCtx context.Context) error {
	ctx, cancel := context.WithCancel(backgroundCtx)

	watcher.cancel = cancel

	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return nil

		case <-time.After(watcher.pollInterval):
			watcher.checkNewBlock()
			watcher.processNextBlock()
		}
	}
}

func (watcher *Watcher) checkNewBlock() {
	resp, err := watcher.rpcClient.GetCurrentBlockNumber()
	if err != nil {
		watcher.logger.Error(
			"rpcClient.GetCurrentBlockNumber failed",
			"error", err,
		)
		return
	}

	blockNumber := resp.Result

	watcher.mu.Lock()
	defer watcher.mu.Unlock()

	if watcher.latestBlock == blockNumber {
		watcher.logger.Debug("no new block, skipping...")
		return
	}

	watcher.logger.Info("new block number", "block number", blockNumber)
	watcher.latestBlock = blockNumber
}

// @TODO: handle case where block still has pending tx
// assumes: block number increases incrementally and linearly.
func (watcher *Watcher) processNextBlock() {
	state := watcher.copyState()

	// @TODO: update in case new subs came in
	if state.currentBlock == state.latestBlock {
		watcher.logger.Debug("no new block to process, skipping")
		return
	}

	watcher.logger.Info("processing next block")
	offset := 1
	if state.currentBlock == "" {
		offset = 0
		state.currentBlock = state.latestBlock
	}

	currentBlockNum, err := strToHex(state.currentBlock)
	if err != nil {
		watcher.logger.Error("could not parse current block number", "currentBlockNum", currentBlockNum)
		return
	}

	nextBlockNum := numToStr(currentBlockNum + offset)
	watcher.logger.Info("processing next block", "nextBlockNum", nextBlockNum)

	resp, err := watcher.rpcClient.GetTransactionCountByNumber(nextBlockNum)
	if err != nil {
		watcher.logger.Error("could not fetch transaction count", "error", err)
		return
	}

	watcher.logger.Info(
		"got transaction count",
		"count", resp.Result,
	)

	count, err := strToHex(resp.Result)
	if err != nil {
		watcher.logger.Error("could not parse transaction count string", "countStr", resp.Result)
		return
	}

	watcher.logger.Info("got block transaction count", "blockNum", nextBlockNum, "count", count)

	const batchSize = 25

	transactions := make([]ethereum.Transaction, count)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	for k := 0; k < count; k++ {
		if k%batchSize == 0 {
			time.Sleep(2 * time.Second)
		}
		index := k

		g.Go(func() error {
			resp, err := watcher.rpcClient.GetTransactionByBlockNumberAndIndex(
				nextBlockNum, numToStr(index),
			)
			if err != nil {
				watcher.logger.Error(
					"failed to get transaction",
					"blockNum", nextBlockNum,
					"index", k,
					"error", err,
				)

				return err
			}

			transactions[index] = resp.Result
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		watcher.logger.Info("fetching transactions failed, block will be reprocessed")
		return
	}

	for k, t := range transactions {
		watcher.logger.Info(
			"transaction",
			"k", k,
			"From", t.From,
			"To", t.To,
		)
	}

	watcher.mu.Lock()
	watcher.currentBlock = nextBlockNum
	watcher.mu.Unlock()
}

type State struct {
	subs         []string
	currentBlock string
	latestBlock  string
}

func (watcher *Watcher) copyState() State {
	watcher.mu.Lock()
	defer watcher.mu.Unlock()

	subs := make([]string, len(watcher.subscriptions))
	copy(subs, watcher.subscriptions)

	return State{
		subs:         subs,
		currentBlock: watcher.currentBlock,
		latestBlock:  watcher.latestBlock,
	}
}
