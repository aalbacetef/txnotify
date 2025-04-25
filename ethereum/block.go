package ethereum

// @NOTE: we only care about a few fields.
type Block struct {
	Hash         string   `json:"hash"`
	Transactions []string `json:"transactions"`
}
