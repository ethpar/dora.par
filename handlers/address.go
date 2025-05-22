package handlers

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/templates"
	"github.com/gorilla/mux"
	"math/big"
	"net/http"
	"strings"
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
}
type Account struct {
	AccountAddress string             `json:"account_address"`
	AccountBalance *big.Float         `json:"account_balance"`
	ERC20Tokens    int                `json:"account_erc20"`
	Transactions   []*TransactionData `json:"transactions"` // Transactions included in this block
}

func Address(w http.ResponseWriter, r *http.Request) {
	var addressTemplateFiles = append(layoutTemplateFiles,
		"address/address.html",
		"address/addressOverview.html",
		"address/addressTransactions.html",
	)
	var addressTemplate = templates.GetTemplate(addressTemplateFiles...)

	vars := mux.Vars(r)
	//address := strings.Replace(vars["address"], "0x", "", -1)
	address := vars["address"]

	if strings.Index(address, "0x") != 0 {
		address = "0x" + address
	}

	data := InitPageData(w, r, "address", "", "", addressTemplateFiles)
	transactions := services.GlobalBeaconService.GetTransactionsForAddress(address)
	var account Account

	for _, dbTransaction := range transactions {

		txValue := weiToEther(new(big.Int).SetUint64(dbTransaction.Value))

		transactionData := &TransactionData{
			Hash:        dbTransaction.Hash,
			BlockNumber: dbTransaction.BlockNumber,
			BlockRank:   dbTransaction.BlockRank,
			TimeStamp:   dbTransaction.TimeStamp,
			From:        dbTransaction.From,
			To:          dbTransaction.To,
			Value:       txValue,
			Method:      dbTransaction.Method,
			Type:        dbTransaction.Type,
			IsFrom:      false,
		}
		if address == transactionData.From {
			transactionData.IsFrom = true
		}
		account.Transactions = append(account.Transactions, transactionData)
	}

	account.AccountAddress = address

	clients := services.GlobalBeaconService.GetExecutionClients()
	if len(clients) > 0 {
		client := clients[0].GetRPCClient()
		balance, err := client.GetEthClient().BalanceAt(context.Background(), common.HexToAddress(address), nil)

		if err == nil {
			account.AccountBalance = weiToEther(new(big.Int).SetUint64(balance.Uint64()))
		}
	}

	data.Data = account
	w.Header().Set("Content-Type", "text/html")
	if handleTemplateError(w, r, "index.go", "Index", "", addressTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}

	// Use the first available client

	/*	w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		//address := strings.Replace(vars["address"], "0x", "", -1)
		address := vars["address"]

		if strings.Index(address, "0x") != 0 {
			address = "0x" + address
		}
		//pageCall := services.FrontendCacheProcessingPage
		//services.GlobalBeaconService.GetTransaction(nil, txHash)
		res := services.GlobalBeaconService.GetTransactionsForAddress(address)
		json.NewEncoder(w).Encode(res)*/

}

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}
