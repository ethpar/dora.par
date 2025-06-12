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

func ApiTransactionsErc20(w http.ResponseWriter, r *http.Request) {
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

	var contract string = ""
	if urlArgs.Has("contract") {
		contract = urlArgs.Get("contract")
	}

	transactionsCount, _ := services.GlobalBeaconService.GetTransactionsErc20CountForAddress(address, contract)

	transactions := services.GlobalBeaconService.GetTransactionsErc20ForAddress(address, contract, offset, pageSize)

	var result models.APITransactionsErc20List

	contracts := services.GlobalBeaconService.GetContracts()

	for _, dbTransaction := range transactions {
		contract := contracts[dbTransaction.To]
		coin := ""
		if contract != nil {
			coin = contract.Symbol
		}

		transactionData := models.ApiTransactionErc20DataListItem{
			Hash:      dbTransaction.Hash,
			Method:    *dbTransaction.Erc20Method,
			TimeStamp: dbTransaction.TimeStamp,
			From:      dbTransaction.From,
			To:        *dbTransaction.Erc20Address,
			Amount:    *dbTransaction.Erc20Value,
			IsFrom:    false,
			Contract:  dbTransaction.To,
			Token:     coin,
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
