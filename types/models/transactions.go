package models

import (
	"math/big"
	"time"
)

type TransactionData struct {
	Hash              string     `json:"hash"`
	BlockNumber       uint64     `json:"block_number"`
	BlockRank         uint64     `json:"block_rank"`
	TimeStamp         time.Time  `json:"created_at"`
	Nonce             uint64     `json:"nonce"`
	BlockHash         string     `json:"block_hash"`
	TransactionIndex  uint       `json:"transaction_index"`
	From              string     `json:"from"`
	To                string     `json:"to"`
	Value             *big.Float `json:"value"`
	Gas               uint64     `json:"gas"`
	GasPrice          uint64     `json:"gas_price"`
	GasPriceGWei      *big.Float `json:"gas_price_gwei"`
	IsError           bool       `json:"is_error"`
	TxReceiptStatus   string     `json:"receipt_status"`
	Input             string     `json:"input"`
	ContractAddress   string     `json:"contract_address"`
	CumulativeGasUsed uint64     `json:"cumulative_gas_used"`
	GasUsed           uint64     `json:"gas_used"`
	Confirmations     int        `json:"confirmations"`
	Method            string     `json:"method"`
	Type              string     `json:"type"`
	IsFrom            bool       `json:"is_from"`
	TxFee             *big.Float `json:"tx_fee"`
}
