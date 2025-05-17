package txnotify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/aalbacetef/txnotify/rpc"
)

type RPCClient interface {
	GetBlockByNumber(blockNum string) (*rpc.Response[ethereum.Block], error)
	GetCurrentBlockNumber() (*rpc.Response[string], error)
}

type Notifier interface {
	Notify(address string, txList []ethereum.Transaction)
}

type Config struct {
	PollInterval time.Duration
	BatchSize    int
	BatchDelay   time.Duration
}

// NewWatcher initializes a new Watcher instance with a JSON-RPC client, logger, in-memory cache, and notifier.
func NewWatcher(rpcEndpoint string, cfg Config, notifier Notifier) (*Watcher, error) {
	client, err := rpc.NewClient(rpc.ClientOptions{
		Endpoint: rpcEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("could not initialize client: %w", err)
	}

	pollInterval := cfg.PollInterval
	if pollInterval == 0 {
		pollInterval = defaultPollInterval
	}

	watcher := &Watcher{
		pollInterval: pollInterval,
		rpcClient:    client,
		logger:       slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		cache:        NewInMemoryCache(),
		notifier:     notifier,
	}

	return watcher, nil
}

const (
	defaultPollInterval = 15 * time.Second
)

type Watcher struct {
	mu            sync.Mutex
	subscriptions []string
	cancel        context.CancelFunc
	pollInterval  time.Duration
	rpcClient     RPCClient
	cache         Cache
	currentBlock  string
	latestBlock   string
	logger        *slog.Logger
	notifier      Notifier
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

// Subscribe registers a new address for monitoring by normalizing and adding it to the subscription list.
func (watcher *Watcher) Subscribe(address string) error {
	watcher.mu.Lock()
	defer watcher.mu.Unlock()

	watcher.subscriptions = append(watcher.subscriptions, normalizeAddress(address))

	return nil
}

// Listen starts the polling loop to watch for new Ethereum blocks and process transactions in real-time.
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

// checkNewBlock will fetch the latest block number, skipping if there are no new ones.
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

	watcher.logger.Info(
		"new block number",
		"block number", blockNumber,
	)
	watcher.latestBlock = blockNumber
}

// processNextBlock determines the next block to process, fetches its transactions, updates cache,
// and triggers notifications.
// @TODO: handle case where block still has pending tx
// assumes: block number increases incrementally and linearly.
func (watcher *Watcher) processNextBlock() { //nolint:funlen
	state := watcher.copyState()

	// @TODO: update in case new subs came in
	if state.currentBlock == state.latestBlock {
		watcher.logger.Debug("no new block to process, skipping")
		return
	}

	offset := 1
	if state.currentBlock == "" {
		offset = 0
		state.currentBlock = state.latestBlock

		watcher.mu.Lock()
		watcher.currentBlock = state.latestBlock
		watcher.mu.Unlock()
	}

	currentBlockNum, err := strToHex(state.currentBlock)
	if err != nil {
		watcher.logger.Error("could not parse current block number", "currentBlockNum", currentBlockNum)
		return
	}

	nextBlockNum := state.currentBlock

	processed, err := watcher.cache.GetBlockProcessed(state.currentBlock)
	if err == nil && processed {
		nextBlockNum = numToStr(currentBlockNum + offset)
	}

	watcher.logger.Info(
		"processing next block",
		"latestBlock", state.latestBlock,
		"nextBlockNum", nextBlockNum,
		"currentBlockNum", state.currentBlock,
	)

	block, err := watcher.fetchBlockInfoIfNotExist(nextBlockNum)
	if err != nil {
		watcher.logger.Error("fetchBlockInfo failed", "blockNum", nextBlockNum, "error", err)
		return
	}

	count := len(block.Transactions)

	watcher.logger.Info(
		"got block info",
		"txCount", count,
	)

	if err := watcher.cache.SetBlockProcessed(nextBlockNum); err != nil {
		watcher.logger.Error(
			"could not set block as processed, will be reprocessed",
			"blockNum", nextBlockNum,
			"error", err,
		)
		return
	}

	watcher.logger.Info(
		"processed block",
		"blockNum", nextBlockNum,
	)

	go watcher.notifyForBlock(nextBlockNum, state.subs)

	watcher.mu.Lock()
	watcher.currentBlock = nextBlockNum
	watcher.mu.Unlock()
}

type State struct {
	subs         []string
	currentBlock string
	latestBlock  string
}

// copyState returns a snapshot of the current block state and subscriptions for safe concurrent access.
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

// fetchBlockInfoIfNotExistchecks if a block is already cached, and fetches it from the blockchain if not.
func (watcher *Watcher) fetchBlockInfoIfNotExist(blockNum string) (ethereum.Block, error) {
	if block, err := watcher.cache.GetBlock(blockNum); err == nil {
		return block, nil
	}

	blockInfoResp, err := watcher.rpcClient.GetBlockByNumber(blockNum)
	if err != nil {
		return ethereum.Block{}, fmt.Errorf("could not get block info: %w", err)
	}

	if err := watcher.cache.AddBlock(blockNum, blockInfoResp.Result); err != nil {
		return ethereum.Block{}, fmt.Errorf("could not store block info to cache: %w", err)
	}

	return blockInfoResp.Result, nil
}

// notifyForBlock filters transactions involving subscribed addresses and invokes the notifier for each address.
func (watcher *Watcher) notifyForBlock(blockNum string, subs []string) {
	txxMap := make(map[string][]ethereum.Transaction, len(subs))

	block, err := watcher.cache.GetBlock(blockNum)
	if err != nil {
		watcher.logger.Error("cache.GetBlock failed", "blockNum", blockNum, "error", err)
		return
	}

	for _, tx := range block.Transactions {
		from := normalizeAddress(tx.From)
		txxMap[from] = append(txxMap[from], tx)

		if tx.To == nil {
			continue
		}

		to := normalizeAddress(*tx.To)

		if to == from {
			continue
		}

		txxMap[to] = append(txxMap[to], tx)
	}

	for _, addr := range subs {
		go watcher.notifier.Notify(addr, txxMap[addr])
	}
}
