package viewmodels

type GetTransaction struct {
	TxHash string           `json:"tx_hash"`
	From   string           `json:"from"`
	To     string           `json:"to"`
	Nonce  uint64           `json:"nonce"`
	Data   []byte           `json:"data"`
	Value  string           `json:"value"`
	Logs   []TransactionLog `json:"logs"`
}

type TransactionLog struct {
	TxHash string `json:"tx_hash"`
	Index  uint   `json:"index"`
	Data   []byte `json:"data"`
}
