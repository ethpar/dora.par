package handlers

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethpandaops/dora/services"
	"github.com/ethpandaops/dora/templates"
	"github.com/ethpandaops/dora/types/models"
	"github.com/gorilla/mux"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type Account struct {
	AccountType       string                                 `json:"account_type"`
	AccountAddress    string                                 `json:"account_address"`
	AccountBalance    *big.Float                             `json:"account_balance"`
	ERC20Tokens       int                                    `json:"account_erc20"`
	Transactions      []*models.TransactionData              `json:"transactions"`       // Transactions included in this block
	TransactionsErc20 []*models.TransactionErc20DataListItem `json:"transactions_erc20"` // Transactions included in this block
	IsDefaultPage     bool                                   `json:"default_page"`
	TotalPages        uint64                                 `json:"total_pages"`
	PageSize          uint64                                 `json:"page_size"`
	CurrentPageIndex  uint64                                 `json:"page_index"`
	PrevPage          uint64                                 `json:"prev_index"`
	NextPage          uint64                                 `json:"next_index"`
	LastPage          uint64                                 `json:"last_index"`
}

func Address(w http.ResponseWriter, r *http.Request) {
	var addressTemplateFiles = append(layoutTemplateFiles,
		"address/address.html",
		"address/addressOverview.html",
		"address/addressTransactions.html",
		"address/addressTransactionsErc20.html",
	)
	var addressTemplate = templates.GetTemplate(addressTemplateFiles...)

	vars := mux.Vars(r)
	//address := strings.Replace(vars["address"], "0x", "", -1)
	address := vars["address"]
	urlArgs := r.URL.Query()
	var start uint64 = 0
	var pageSize uint64 = 50
	if urlArgs.Has("c") {
		pageSize, _ = strconv.ParseUint(urlArgs.Get("c"), 10, 64)
	}
	var currentPage uint64 = 1
	if urlArgs.Has("s") {
		currentPage, _ = strconv.ParseUint(urlArgs.Get("s"), 10, 64)
	}

	if strings.Index(address, "0x") != 0 {
		address = "0x" + address
	}

	data := InitPageData(w, r, "address", "", "", addressTemplateFiles)

	var account Account
	account.AccountAddress = address

	initTransactions(&account, start, pageSize, currentPage)

	initTransactionsErc20(&account, start, pageSize, currentPage)

	clients := services.GlobalBeaconService.GetExecutionClients()
	if len(clients) > 0 {
		client := clients[0].GetRPCClient()
		balance, err := client.GetEthClient().BalanceAt(context.Background(), common.HexToAddress(address), nil)
		code, err := client.GetEthClient().CodeAt(context.Background(), common.HexToAddress(address), nil)
		if len(code) == 0 {
			account.AccountType = "Address"
		} else {
			account.AccountType = "Contract"
		}
		if err == nil {
			account.AccountBalance = weiToEther(balance)
		}
	}

	if currentPage == 1 {
		account.PrevPage = 1
	} else {
		account.PrevPage = currentPage - 1
	}
	if currentPage == account.TotalPages {
		account.NextPage = account.TotalPages
	} else {
		account.NextPage = currentPage + 1
	}
	account.LastPage = account.TotalPages

	data.Data = account
	w.Header().Set("Content-Type", "text/html")
	if handleTemplateError(w, r, "index.go", "Index", "", addressTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}

	// Use the first available client

	/*	w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		//address := strings.Replace(vars["address"], "0x", "", -1)
		address := vars["address"]

		if strings.Index(address, "0x") != 0 {
			address = "0x" + address
		}
		//pageCall := services.FrontendCacheProcessingPage
		//services.GlobalBeaconService.GetTransaction(nil, txHash)
		res := services.GlobalBeaconService.GetTransactionsForAddress(address)
		json.NewEncoder(w).Encode(res)*/

}

func initTransactions(account *Account, start uint64, pageSize uint64, currentPage uint64) {
	transactionsCount, _ := services.GlobalBeaconService.GetTransactionsCountForAddress(account.AccountAddress)
	totalPages := transactionsCount / pageSize
	start = pageSize * (currentPage - 1)
	transactions := services.GlobalBeaconService.GetTransactionsForAddress(account.AccountAddress, start, pageSize)

	for _, dbTransaction := range transactions {

		v := new(big.Int)
		v.SetString(dbTransaction.Value, 10)
		txValue := weiToEther(v)

		method := dbTransaction.Method
		if dbTransaction.Erc20Method != nil {
			method = *dbTransaction.Erc20Method
		}
		transactionData := &models.TransactionData{
			Hash:        dbTransaction.Hash,
			BlockNumber: dbTransaction.BlockNumber,
			BlockRank:   dbTransaction.BlockRank,
			TimeStamp:   dbTransaction.TimeStamp,
			From:        dbTransaction.From,
			To:          dbTransaction.To,
			Value:       txValue,
			Method:      method,
			Type:        dbTransaction.Type,
			IsFrom:      false,
		}
		if account.AccountAddress == transactionData.From {
			transactionData.IsFrom = true
		}
		account.Transactions = append(account.Transactions, transactionData)
	}

	account.PageSize = pageSize
	account.TotalPages = totalPages
	account.CurrentPageIndex = currentPage
	account.IsDefaultPage = true
}

func initTransactionsErc20(account *Account, start uint64, pageSize uint64, currentPage uint64) {
	transactionsCount, _ := services.GlobalBeaconService.GetTransactionsErc20CountForAddress(account.AccountAddress)
	totalPages := transactionsCount / pageSize
	start = pageSize * (currentPage - 1)
	transactions := services.GlobalBeaconService.GetTransactionsErc20ForAddress(account.AccountAddress, start, pageSize)

	contracts := services.GlobalBeaconService.GetContracts()

	for _, dbTransaction := range transactions {
		contract := contracts[dbTransaction.To]
		coin := ""
		if contract != nil {
			coin = contract.Symbol
		}

		v := new(big.Int)
		s := dbTransaction.Erc20Value
		v.SetString(*s, 10)
		txValue := weiToEther(v)

		transactionData := &models.TransactionErc20DataListItem{
			Hash:        dbTransaction.Hash,
			BlockNumber: dbTransaction.BlockNumber,
			BlockRank:   dbTransaction.BlockRank,
			Method:      *dbTransaction.Erc20Method,
			TimeStamp:   dbTransaction.TimeStamp,
			From:        dbTransaction.From,
			To:          *dbTransaction.Erc20Address,
			Amount:      txValue,
			IsFrom:      false,
			Contract:    dbTransaction.To,
			Coin:        coin,
		}

		if account.AccountAddress == transactionData.From {
			transactionData.IsFrom = true
		}
		account.TransactionsErc20 = append(account.TransactionsErc20, transactionData)
	}

	account.PageSize = pageSize
	account.TotalPages = totalPages
	account.CurrentPageIndex = currentPage
	account.IsDefaultPage = true
}
func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func weiToGWei(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
}
