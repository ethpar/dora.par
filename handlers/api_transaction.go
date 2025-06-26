package handlers

import (
	"encoding/json"
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/types/models"
	"github.com/gorilla/mux"
	"math/big"
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
	dbTransaction := services.GlobalBeaconService.GetTransactionByHash(hash)

	contracts := services.GlobalBeaconService.GetContracts()
	contract := contracts[dbTransaction.To]
	coin := ""
	contractName := ""
	if contract != nil {
		coin = contract.Symbol
		contractName = contract.Name
	}

	var erc20Value big.Int
	var erc20Address string
	var erc20Method string
	isErc20 := false
	if dbTransaction.Erc20Value != nil {
		isErc20 = true
		erc20Value1, _ := new(big.Int).SetString(*dbTransaction.Erc20Value, 10)
		erc20Value = *erc20Value1
	}
	if dbTransaction.Erc20Method != nil {
		erc20Method = *dbTransaction.Erc20Method
	}
	if dbTransaction.Erc20Address != nil {
		erc20Address = *dbTransaction.Erc20Address
	}
	txFee := new(big.Int).SetUint64(dbTransaction.GasPrice * dbTransaction.GasUsed)

	value, _ := new(big.Int).SetString(dbTransaction.Value, 10)

	transactionData := &models.ApiTransactionData{
		Hash:             dbTransaction.Hash,
		BlockNumber:      dbTransaction.BlockNumber,
		BlockRank:        dbTransaction.BlockRank,
		TimeStamp:        dbTransaction.TimeStamp,
		Nonce:            dbTransaction.Nonce,
		TransactionIndex: dbTransaction.TransactionIndex,
		From:             dbTransaction.From,
		To:               dbTransaction.To,
		Value:            *value,
		Type:             dbTransaction.Type,
		GasPrice:         dbTransaction.GasPrice,
		GasUsed:          dbTransaction.GasUsed,
		//	Gas:               dbTransaction.Gas,
		CumulativeGasUsed: dbTransaction.CumulativeGasUsed,
		Input:             dbTransaction.Input,
		TxFee:             txFee,
		Erc20Method:       erc20Method,
		Erc20Address:      erc20Address,
		Erc20Value:        erc20Value,
		Contract:          contractName,
		Coin:              coin,
		IsErc20:           isErc20,
	}

	json.NewEncoder(w).Encode(transactionData)

}
