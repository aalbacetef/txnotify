package rpc

// Transaction represents a JSON-RPC transaction object as per the Ethereum JSON-RPC specification.
type Transaction struct {
	AccessList *[]AccessListEntry `json:"accessList,omitempty"`

	// BlockHash is the 32-byte hash of the block including this transaction.
	// Null when the transaction is pending.
	BlockHash *string `json:"blockHash,omitempty"`

	// BlockNumber is the number of the block including this transaction.
	// Null when the transaction is pending.
	BlockNumber *string `json:"blockNumber,omitempty"`

	// ChainID is the optional chain ID specifying the network (e.g., "0x1" for Ethereum mainnet).
	// Returned only for EIP-1559 transactions.
	ChainID *string `json:"chainId,omitempty"`

	// From is the 20-byte address of the sender.
	From     string `json:"from"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`

	Hash string `json:"hash"`

	Input string `json:"input"`

	MaxPriorityFeePerGas *string `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         *string `json:"maxFeePerGas,omitempty"`

	// Nonce is the number of transactions made by the sender prior to this one, encoded as a hexadecimal string.
	Nonce string `json:"nonce"`

	R string `json:"r"`
	S string `json:"s"`

	// To is the 20-byte address of the receiver.
	// Null for contract creation transactions.
	To *string `json:"to,omitempty"`

	// TransactionIndex is the transaction's index position in the block, encoded as a hexadecimal string.
	// Null when the transaction is pending.
	TransactionIndex *string `json:"transactionIndex,omitempty"`

	// Type is the transaction type (e.g., "0x0" for legacy, "0x2" for EIP-1559), encoded as a hexadecimal string.
	Type string `json:"type"`

	V string `json:"v"`

	Value string `json:"value"`

	YParity *string `json:"yParity,omitempty"`
}

// AccessListEntry represents an entry in the access list for access list transactions (EIP-2930).
type AccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}
