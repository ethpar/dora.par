package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethpandaops/dora/services"
)

type EtherscanResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type EtherscanTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	IsError          string `json:"isError"`
	TxReceiptStatus  string `json:"txreceipt_status"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed          string `json:"gasUsed"`
	Confirmations    string `json:"confirmations"`
}

func EtherscanAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	module := r.URL.Query().Get("module")
	action := r.URL.Query().Get("action")
	address := r.URL.Query().Get("address")
	contractaddress := r.URL.Query().Get("contractaddress")

	switch module {
	case "account":
		switch action {
		case "balance":
			handleBalance(w, address)
		case "txlist":
			handleTxList(w, address)
		case "tokentx":
			handleTokenTx(w, address)
		case "tokenbalance":
			handleTokenBalance(w, address, contractaddress)
		default:
			json.NewEncoder(w).Encode(EtherscanResponse{
				Status:  "0",
				Message: "Invalid action",
				Result:  nil,
			})
		}
	case "contract":
		switch action {
		case "getabi":
			handleGetABI(w, address)
		default:
			json.NewEncoder(w).Encode(EtherscanResponse{
				Status:  "0",
				Message: "Invalid action",
				Result:  nil,
			})
		}
	default:
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Invalid module",
			Result:  nil,
		})
	}
}

func handleBalance(w http.ResponseWriter, address string) {
	if !common.IsHexAddress(address) {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Invalid address format",
			Result:  nil,
		})
		return
	}

	clients := services.GlobalBeaconService.GetExecutionClients()
	if len(clients) == 0 {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "No execution clients available",
			Result:  nil,
		})
		return
	}

	// Use the first available client
	client := clients[0].GetRPCClient()
	balance, err := client.GetEthClient().BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error getting balance",
			Result:  nil,
		})
		return
	}

	json.NewEncoder(w).Encode(EtherscanResponse{
		Status:  "1",
		Message: "OK",
		Result:  balance.String(),
	})
}

func handleTxList(w http.ResponseWriter, address string) {
	if !common.IsHexAddress(address) {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Invalid address format",
			Result:  nil,
		})
		return
	}

	clients := services.GlobalBeaconService.GetExecutionClients()
	if len(clients) == 0 {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "No execution clients available",
			Result:  nil,
		})
		return
	}

	client := clients[0].GetRPCClient()
	ethClient := client.GetEthClient()

	// Get the latest block number
	latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error getting latest block",
			Result:  nil,
		})
		return
	}

	// Get transactions for the last 10,000 blocks (Etherscan default)
	startBlock := big.NewInt(0).Sub(latestBlock.Number(), big.NewInt(10000))
	if startBlock.Sign() < 0 {
		startBlock = big.NewInt(0)
	}

	var txs []map[string]interface{}
	for i := startBlock.Int64(); i <= latestBlock.Number().Int64(); i++ {
		block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			continue
		}

		for _, tx := range block.Transactions() {
			if tx.To() != nil && *tx.To() == common.HexToAddress(address) {
				receipt, err := ethClient.TransactionReceipt(context.Background(), tx.Hash())
				if err != nil {
					continue
				}

				txs = append(txs, map[string]interface{}{
					"hash":             tx.Hash().Hex(),
					"nonce":           fmt.Sprintf("%#x", tx.Nonce()),
					"blockHash":       block.Hash().Hex(),
					"blockNumber":     fmt.Sprintf("%#x", block.Number()),
					"transactionIndex": fmt.Sprintf("%#x", receipt.TransactionIndex),
					"from":            func() string {
						sender, err := ethClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
						if err != nil {
							return "0x0000000000000000000000000000000000000000"
						}
						return sender.Hex()
					}(),
					"to":              tx.To().Hex(),
					"value":           fmt.Sprintf("%#x", tx.Value()),
					"gas":             fmt.Sprintf("%#x", tx.Gas()),
					"gasPrice":        fmt.Sprintf("%#x", tx.GasPrice()),
					"isError":         "0",
					"timeStamp":       fmt.Sprintf("%#x", block.Time()),
				})
			}
		}
	}

	json.NewEncoder(w).Encode(EtherscanResponse{
		Status:  "1",
		Message: "OK",
		Result:  txs,
	})
}
