package rpc

import "github.com/aalbacetef/txnotify/ethereum"

const (
	getCurrentBlockMethod    = "eth_blockNumber"
	getBlockByNumberEndpoint = "eth_getBlockByNumber"
)

func (client *Client) GetCurrentBlockNumber() (*Response[string], error) {
	endpoint := getCurrentBlockMethod

	return Do[string](client, endpoint, []any{})
}

// GetBlockByNumber will return block information (hash and transactions) given
// the block's number as a hex-string.
func (client *Client) GetBlockByNumber(blockNum string) (*Response[ethereum.Block], error) {
	endpoint := getBlockByNumberEndpoint

	const getFullBlock = true

	params := []any{
		blockNum,
		getFullBlock,
	}

	return Do[ethereum.Block](client, endpoint, params)
}
