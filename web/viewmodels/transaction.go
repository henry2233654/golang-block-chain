package viewmodels

type GetTransaction struct {
	TxHash    string           `json:"tx_hash"`
	BlockHash string           `json:"block_hash"`
	From      string           `json:"from"`
	To        string           `json:"to"`
	Nonce     uint64           `json:"nonce"`
	Data      string           `json:"data"`
	Value     string           `json:"value"`
	Logs      []TransactionLog `json:"logs"`
}

type TransactionLog struct {
	TxHash string `json:"tx_hash"`
	Index  uint64 `json:"index"`
	Data   string `json:"data"`
}
