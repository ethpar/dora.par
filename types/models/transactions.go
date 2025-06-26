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
	Type              uint8      `json:"type"`
	IsFrom            bool       `json:"is_from"`
	TxFee             *big.Float `json:"tx_fee"`
	Erc20Method       string     `json:"erc20_method"`
	Erc20Address      string     `json:"erc20_address_to"`
	Erc20Value        *big.Float `json:"erc20_value"`
	Coin              string     `json:"coin"`
	Contract          string     `json:"contract"`
	IsErc20           bool       `json:"is_erc20"`
}

type ApiTransactionData struct {
	Hash             string    `json:"hash"`
	BlockNumber      uint64    `json:"block_number"`
	BlockRank        uint64    `json:"block_rank"`
	TimeStamp        time.Time `json:"created_at"`
	Nonce            uint64    `json:"nonce"`
	TransactionIndex uint      `json:"transaction_index"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	Value            big.Int   `json:"value"`
	//Gas               uint64    `json:"gas"`
	GasPrice          uint64   `json:"gas_price"`
	Input             string   `json:"input"`
	CumulativeGasUsed uint64   `json:"cumulative_gas_used"`
	GasUsed           uint64   `json:"gas_used"`
	Type              uint8    `json:"type"`
	TxFee             *big.Int `json:"tx_fee"`
	Erc20Method       string   `json:"erc20_method"`
	Erc20Address      string   `json:"erc20_address_to"`
	Erc20Value        big.Int  `json:"erc20_value"`
	Coin              string   `json:"erc20_coin"`
	Contract          string   `json:"erc20_contract"`
	IsErc20           bool     `json:"is_erc20"`
}

type TransactionDataList struct {
	Hash      string    `json:"hash"`
	TimeStamp time.Time `json:"created_at"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Value     string    `json:"value"`
	IsFrom    bool      `json:"is_from"`
}

type TransactionErc20DataListItem struct {
	Hash        string     `json:"hash"`
	BlockNumber uint64     `json:"block_number"`
	BlockRank   uint64     `json:"block_rank"`
	Method      string     `json:"method"`
	TimeStamp   time.Time  `json:"created_at"`
	From        string     `json:"from"`
	To          string     `json:"to"`
	Amount      *big.Float `json:"amount"`
	IsFrom      bool       `json:"is_from"`
	Coin        string     `json:"coin"`
	Contract    string     `json:"contract"`
}

type ApiTransactionErc20DataListItem struct {
	Hash      string    `json:"hash"`
	Method    string    `json:"method"`
	TimeStamp time.Time `json:"created_at"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    string    `json:"amount"`
	IsFrom    bool      `json:"is_from"`
	Token     string    `json:"token"`
	Contract  string    `json:"contract"`
}

type APITransactionsList struct {
	Data       []TransactionDataList `json:"data"`
	TotalCount uint64                `json:"totalCount"`
	PageSize   uint64                `json:"pageSize"`
	Offset     uint64                `json:"offset"`
}

type APITransactionsErc20List struct {
	Data       []ApiTransactionErc20DataListItem `json:"data"`
	TotalCount uint64                            `json:"totalCount"`
	PageSize   uint64                            `json:"pageSize"`
	Offset     uint64                            `json:"offset"`
}
