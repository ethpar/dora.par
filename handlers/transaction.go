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

	v := new(big.Int)
	v.SetString(dbTransaction.Value, 10)
	txValue := weiToEther(v)

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
	}

	data.Data = transactionData

	w.Header().Set("Content-Type", "text/html")

	if handleTemplateError(w, r, "transaction.go", "Transaction", "", transactionTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}
