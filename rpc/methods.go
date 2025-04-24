package rpc

import "github.com/aalbacetef/txnotify/ethereum"

const (
	getCurrentBlockMethod                       = "eth_blockNumber"
	getBlockByNumberEndpoint                    = "eth_getBlockByNumber"
	getTransactionCountByNumberEndpoint         = "eth_getBlockTransactionCountByNumber"
	getTransactionByBlockNumberAndIndexEndpoint = "eth_getTransactionByBlockNumberAndIndex"
)

func (client *Client) GetCurrentBlockNumber() (*Response[string], error) {
	endpoint := getCurrentBlockMethod

	response, err := Do[string](client, endpoint, []any{})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetBlockByNumber will return block information (hash and transaction hashes) given the block's number as a hex-string.
func (client *Client) GetBlockByNumber(blockNum string) (*Response[ethereum.Block], error) {
	endpoint := getBlockByNumberEndpoint

	const getFullBlock = false

	params := []any{
		blockNum,
		getFullBlock,
	}

	response, err := Do[ethereum.Block](client, endpoint, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetTransactionCountByNumber will fetch the transaction count for a block. Result is a hex-string corresponding to the transaction count. It expects the blockNum to be a hex-string.
func (client *Client) GetTransactionCountByNumber(blockNum string) (*Response[string], error) {
	endpoint := getTransactionCountByNumberEndpoint
	params := []any{blockNum}

	response, err := Do[string](client, endpoint, params)

	return response, err
}

// GetTransactionByBlockNumberAndIndex will fetch the transaction corresponding to the given block and index. It expects both inputs to be hex strings.
func (client *Client) GetTransactionByBlockNumberAndIndex(blockNum, index string) (*Response[ethereum.Transaction], error) {
	endpoint := getTransactionByBlockNumberAndIndexEndpoint

	params := []any{blockNum, index}

	response, err := Do[ethereum.Transaction](client, endpoint, params)

	return response, err
}
