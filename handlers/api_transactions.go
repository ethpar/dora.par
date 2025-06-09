package handlers

import (
	"encoding/json"
	"github.com/ethpandaops/dora/services"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func ApiTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	address := vars["address"]

	if strings.Index(address, "0x") != 0 {
		address = "0x" + address
	}
	res := services.GlobalBeaconService.GetAllTransactionsForAddress(address)
	json.NewEncoder(w).Encode(res)

}
