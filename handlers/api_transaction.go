package handlers

import (
	"encoding/json"
	"github.com/ethpandaops/dora/services"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func ApiTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	hash := vars["hash"]

	if strings.Index(hash, "0x") != 0 {
		hash = "0x" + hash
	}
	res := services.GlobalBeaconService.GetTransactionByHash(hash)
	json.NewEncoder(w).Encode(res)

}
