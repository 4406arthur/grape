package repository

type Block struct {
	BlockNum   int64  `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	ParentHash string `json:"parent_hash"`
	BlockTime  uint64 `json:"block_time"`
}

type Tx struct {
	TxHash   string `json:"tx_hash"`
	BlockNum int64
	TxFrom   string `json:"from"`
	TxTo     string `json:"to"`
	TxNonce  uint64 `json:"nonce"`
	TxData   string `json:"data"`
	TxValue  string `json:"value"`
}
