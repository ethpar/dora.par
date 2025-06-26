package handlers

import (
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/templates"
	"github.com/ethpandaops/dora/types/models"
	"github.com/gorilla/mux"
	"math/big"
	"net/http"
	"strings"
)

func Transaction(w http.ResponseWriter, r *http.Request) {
	var transactionTemplateFiles = append(layoutTemplateFiles,
		"transaction/transaction.html",
	)
	var transactionTemplate = templates.GetTemplate(transactionTemplateFiles...)
	data := InitPageData(w, r, "transaction", "", "", transactionTemplateFiles)

	vars := mux.Vars(r)
	//address := strings.Replace(vars["address"], "0x", "", -1)
	hash := vars["hash"]

	if strings.Index(hash, "0x") != 0 {
		hash = "0x" + hash
	}

	dbTransaction := services.GlobalBeaconService.GetTransactionByHash(hash)

	if dbTransaction == nil {
		return
	}

	txValue := weiToEtherS(dbTransaction.Value)

	contracts := services.GlobalBeaconService.GetContracts()
	contract := contracts[dbTransaction.To]
	coin := ""
	contractName := ""
	if contract != nil {
		coin = contract.Symbol
		contractName = contract.Name
	}

	var erc20Value *big.Float
	var erc20Address string
	var erc20Method string
	isErc20 := false
	if dbTransaction.Erc20Value != nil {
		isErc20 = true
		erc20Value = weiToEtherS(*dbTransaction.Erc20Value)
	}
	if dbTransaction.Erc20Method != nil {
		erc20Method = *dbTransaction.Erc20Method
	}
	if dbTransaction.Erc20Address != nil {
		erc20Address = *dbTransaction.Erc20Address
	}

	gasPriceGWei := weiToGWei(new(big.Int).SetUint64(dbTransaction.GasPrice))
	txFee := weiToEther(new(big.Int).SetUint64(dbTransaction.GasPrice * dbTransaction.GasUsed))
	transactionData := &models.TransactionData{
		Hash:         dbTransaction.Hash,
		BlockNumber:  dbTransaction.BlockNumber,
		BlockRank:    dbTransaction.BlockRank,
		TimeStamp:    dbTransaction.TimeStamp,
		From:         dbTransaction.From,
		To:           dbTransaction.To,
		Value:        txValue,
		Method:       dbTransaction.Method,
		Type:         dbTransaction.Type,
		IsFrom:       false,
		GasPrice:     dbTransaction.GasPrice,
		GasPriceGWei: gasPriceGWei,
		GasUsed:      dbTransaction.GasUsed,
		TxFee:        txFee,
		Erc20Method:  erc20Method,
		Erc20Address: erc20Address,
		Erc20Value:   erc20Value,
		Contract:     contractName,
		Coin:         coin,
		IsErc20:      isErc20,
	}

	data.Data = transactionData

	w.Header().Set("Content-Type", "text/html")

	if handleTemplateError(w, r, "transaction.go", "Transaction", "", transactionTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}
