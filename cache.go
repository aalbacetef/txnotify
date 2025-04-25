package txnotify

import (
	"fmt"
	"sync"

	"github.com/aalbacetef/txnotify/ethereum"
)

type Cache interface {
	AddBlock(blockNum string, block ethereum.Block) error
	GetBlock(blockNum string) (ethereum.Block, error)
	GetBlockProcessed(blockNum string) (bool, error)
	SetBlockProcessed(blockNum string) error
	AddTx(tx ethereum.Transaction) error
	GetTx(hash string) (ethereum.Transaction, error)
	TxForAddress(address string) ([]ethereum.Transaction, error)
	Subscribe(address string) error
	Unsubscribe(address string) error
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		subscribedAddress: make([]string, 0),
		transactions:      make(map[string]ethereum.Transaction),
		blocks:            make(map[string]ethereum.Block),
		processedBlocks:   make(map[string]bool),
	}
}

type InMemoryCache struct {
	mu sync.Mutex

	subscribedAddress []string
	transactions      map[string]ethereum.Transaction
	blocks            map[string]ethereum.Block
	processedBlocks   map[string]bool
}

func (cache *InMemoryCache) AddBlock(blockNum string, block ethereum.Block) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, exists := cache.blocks[blockNum]; exists {
		return nil
	}

	cache.blocks[blockNum] = block

	return nil
}

func (cache *InMemoryCache) GetBlock(blockNum string) (ethereum.Block, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	block, exists := cache.blocks[blockNum]
	if !exists {
		return ethereum.Block{}, fmt.Errorf("block with number %s not found", blockNum)
	}

	return block, nil
}

func (cache *InMemoryCache) GetBlockProcessed(blockNum string) (bool, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	v, found := cache.processedBlocks[blockNum]
	if !found {
		return false, fmt.Errorf("block with number %s not found", blockNum)
	}

	return v, nil
}

func (cache *InMemoryCache) SetBlockProcessed(blockNum string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, found := cache.blocks[blockNum]; !found {
		return fmt.Errorf("block with number %s not found", blockNum)
	}

	cache.processedBlocks[blockNum] = true

	return nil
}

func (cache *InMemoryCache) AddTx(tx ethereum.Transaction) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, exists := cache.transactions[tx.Hash]; exists {
		return nil
	}

	cache.transactions[tx.Hash] = tx
	return nil
}

type TxNotFoundError struct {
	hash string
}

func (e TxNotFoundError) Error() string {
	return fmt.Sprintf("transaction with hash %s not found", e.hash)
}

func (cache *InMemoryCache) GetTx(hash string) (ethereum.Transaction, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	tx, exists := cache.transactions[hash]
	if !exists {
		return ethereum.Transaction{}, TxNotFoundError{hash}
	}

	return tx, nil
}

func (cache *InMemoryCache) TxForAddress(address string) ([]ethereum.Transaction, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	var result []ethereum.Transaction

	for _, tx := range cache.transactions {
		if tx.From == address || (tx.To != nil && *tx.To == address) {
			result = append(result, tx)
		}
	}

	return result, nil
}

func (cache *InMemoryCache) Subscribe(address string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for _, addr := range cache.subscribedAddress {
		if addr == address {
			return nil
		}
	}

	cache.subscribedAddress = append(cache.subscribedAddress, address)
	return nil
}

func (cache *InMemoryCache) Unsubscribe(address string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for i, addr := range cache.subscribedAddress {
		if addr == address {
			cache.subscribedAddress = append(cache.subscribedAddress[:i], cache.subscribedAddress[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("address %s not subscribed", address)
}
