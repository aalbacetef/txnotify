package rpc

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

// NOTE: We only care about a few fields.
type GetBlockByNumberPayload struct {
	Hash         string   `json:"hash"`
	Transactions []string `json:"transactions"`
}

// GetBlockByNumber will return block information (hash and transaction hashes) given the block's number as a hex-string.
func (client *Client) GetBlockByNumber(blockNum string) (*Response[GetBlockByNumberPayload], error) {
	endpoint := getBlockByNumberEndpoint

	const getFullBlock = false

	params := []any{
		blockNum,
		getFullBlock,
	}

	response, err := Do[GetBlockByNumberPayload](client, endpoint, params)
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
func (client *Client) GetTransactionByBlockNumberAndIndex(blockNum, index string) (*Response[Transaction], error) {
	endpoint := getTransactionByBlockNumberAndIndexEndpoint

	params := []any{blockNum, index}

	response, err := Do[Transaction](client, endpoint, params)

	return response, err
}
