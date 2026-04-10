package handlers

import (
	"encoding/json"
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/types/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func ApiTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	address := vars["address"]

	if strings.Index(address, "0x") != 0 {
		address = "0x" + address
	}

	urlArgs := r.URL.Query()
	var pageSize uint64 = 50
	if urlArgs.Has("pageSize") {
		pageSize, _ = strconv.ParseUint(urlArgs.Get("pageSize"), 10, 64)
	}
	var offset uint64 = 0
	if urlArgs.Has("offset") {
		offset, _ = strconv.ParseUint(urlArgs.Get("offset"), 10, 64)
	}

	transactionsCount, _ := services.GlobalBeaconService.GetTransactionsCountForAddress(address)

	transactions := services.GlobalBeaconService.GetTransactionsForAddress(address, offset, pageSize, transactionsCount)

	var result models.APITransactionsList

	for _, dbTransaction := range transactions {

		transactionData := models.TransactionDataList{
			Hash:      dbTransaction.Hash,
			TimeStamp: dbTransaction.TimeStamp,
			From:      dbTransaction.From,
			To:        dbTransaction.To,
			Value:     dbTransaction.Value,
			IsFrom:    false,
		}
		if address == transactionData.From {
			transactionData.IsFrom = true
		}
		result.Data = append(result.Data, transactionData)
	}

	result.PageSize = pageSize
	result.Offset = offset
	result.TotalCount = transactionsCount

	json.NewEncoder(w).Encode(result)

}
