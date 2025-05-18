package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethpandaops/dora/services"
)

// ERC20 standard transfer event
var transferEventSig = []byte("Transfer(address,address,uint256)")

// Standard ERC20 ABI methods we care about
const erc20ABI = `[
	{
		"constant": true,
		"inputs": [{"name": "_owner", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "balance", "type": "uint256"}],
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "from", "type": "address"},
			{"indexed": true, "name": "to", "type": "address"},
			{"indexed": false, "name": "value", "type": "uint256"}
		],
		"name": "Transfer",
		"type": "event"
	}
]`

type ERC20Transfer struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	ContractAddress  string `json:"contractAddress"`
	TokenName        string `json:"tokenName"`
	TokenSymbol      string `json:"tokenSymbol"`
	TokenDecimal     string `json:"tokenDecimal"`
	TransactionIndex string `json:"transactionIndex"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	GasUsed          string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Confirmations    string `json:"confirmations"`
}

func GetTokenTx(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	handleTokenTx(w, address)
}

func handleTokenTx(w http.ResponseWriter, address string) {
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

	// Get the latest block number
	latestBlock, err := client.GetEthClient().BlockByNumber(context.Background(), nil)
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error getting latest block",
			Result:  nil,
		})
		return
	}

	// Get token transfers for the last 10,000 blocks (Etherscan default)
	startBlock := big.NewInt(0).Sub(latestBlock.Number(), big.NewInt(10000))
	if startBlock.Sign() < 0 {
		startBlock = big.NewInt(0)
	}

	// Create Transfer event signature
	transferSigHash := common.BytesToHash(transferEventSig)
	targetAddr := common.HexToAddress(address)

	var transfers []ERC20Transfer
	for i := startBlock.Int64(); i <= latestBlock.Number().Int64(); i++ {
		block, err := client.GetEthClient().BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			continue
		}

		for _, tx := range block.Transactions() {
			receipt, err := client.GetEthClient().TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				continue
			}

			for _, log := range receipt.Logs {
				// Check if this log is a Transfer event
				if len(log.Topics) == 3 && log.Topics[0] == transferSigHash {
					from := common.BytesToAddress(log.Topics[1].Bytes())
					to := common.BytesToAddress(log.Topics[2].Bytes())

					// Check if the address is involved in the transfer
					if from == targetAddr || to == targetAddr {
						transfers = append(transfers, ERC20Transfer{
							BlockNumber:       fmt.Sprintf("%d", block.Number().Uint64()),
							TimeStamp:        fmt.Sprintf("%d", block.Time()),
							Hash:             tx.Hash().Hex(),
							From:             from.Hex(),
							To:              to.Hex(),
							ContractAddress: log.Address.Hex(),
							Value:           fmt.Sprintf("%#x", new(big.Int).SetBytes(log.Data)),
							TokenName:       "", // Would need contract call to get these
							TokenSymbol:     "",
							TokenDecimal:    "18", // Most common, but would need contract call to verify
							TransactionIndex: fmt.Sprintf("%d", receipt.TransactionIndex),
							Gas:             fmt.Sprintf("%d", tx.Gas()),
							GasPrice:        fmt.Sprintf("%d", tx.GasPrice().Uint64()),
							GasUsed:         fmt.Sprintf("%d", receipt.GasUsed),
							CumulativeGasUsed: fmt.Sprintf("%d", receipt.CumulativeGasUsed),
							Confirmations:    fmt.Sprintf("%d", latestBlock.Number().Uint64()-block.Number().Uint64()),
						})
					}
				}
			}
		}
	}

	json.NewEncoder(w).Encode(EtherscanResponse{
		Status:  "1",
		Message: "OK",
		Result:  transfers,
	})
}

func GetTokenBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	contractAddress := r.URL.Query().Get("contractaddress")
	handleTokenBalance(w, address, contractAddress)
}

func handleTokenBalance(w http.ResponseWriter, address string, contractAddress string) {
	if !common.IsHexAddress(address) || !common.IsHexAddress(contractAddress) {
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

	// Parse ERC20 ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error parsing ABI",
			Result:  nil,
		})
		return
	}

	// Pack the balanceOf function call
	data, err := parsedABI.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error packing function call",
			Result:  nil,
		})
		return
	}

	// Make the call using RPC
	var result string
	err = client.GetEthClient().Client().CallContext(context.Background(), &result, "eth_call", map[string]interface{}{
		"to":   contractAddress,
		"data": common.Bytes2Hex(data),
	}, "latest")

	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error calling contract",
			Result:  nil,
		})
		return
	}

	// Unpack the result
	balance := new(big.Int)
	resultBytes := common.FromHex(result)
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", resultBytes)
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error unpacking result",
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

func GetContractABI(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	handleGetABI(w, address)
}

func handleGetABI(w http.ResponseWriter, address string) {
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

	// Get contract code to check if it exists
	code, err := client.GetEthClient().CodeAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Error getting contract code",
			Result:  nil,
		})
		return
	}

	if len(code) == 0 {
		json.NewEncoder(w).Encode(EtherscanResponse{
			Status:  "0",
			Message: "Contract source code not verified",
			Result:  "",
		})
		return
	}

	// Note: In a real implementation, you would need to maintain a database of verified contract ABIs
	// For now, we'll just return a message that verification is required
	json.NewEncoder(w).Encode(EtherscanResponse{
		Status:  "0",
		Message: "Contract source code not verified. Please verify the contract to see the ABI.",
		Result:  "",
	})
}
