package handlers

import (
	"encoding/json"
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/types/models"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func ApiTransactionsAfter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	address := vars["address"]

	if strings.Index(address, "0x") != 0 {
		address = "0x" + address
	}

	urlArgs := r.URL.Query()
	var result models.APITransactionsAfterList

	if urlArgs.Has("after") {
		layout := "2006-01-02T15:04:05Z"
		value := urlArgs.Get("after")
		after, _ := time.Parse(layout, value)
		transactions := services.GlobalBeaconService.GetTransactionsForAddressAfter(address, after)
		for _, transaction := range transactions {
			result.Transactions = append(result.Transactions, *transaction)
		}
		json.NewEncoder(w).Encode(result)
	}

}
