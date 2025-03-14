package beacon

import (
	"encoding/json"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/dora/clients/execution"
	"github.com/ethpandaops/dora/utils"
	"github.com/senseyeio/roger"
	"os"
	"strconv"
	"time"
)

type PendingTransaction struct {
	BlockHash            string        `json:"blockHash,omitempty"`
	BlockNumber          string        `json:"blockNumber,omitempty"`
	From                 string        `json:"from,omitempty"`
	Gas                  string        `json:"gas,omitempty"`
	GasPrice             string        `json:"gasPrice,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	YParity              string        `json:"yParity"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	AccessList           []interface{} `json:"accessList"`
	ChainId              string        `json:"chainId"`
	PublicKey            string        `json:"publicKey"`
	Raw                  string        `json:"raw"`
	Type                 string        `json:"type"`
}

func processPendingTransactions(c *Client, slot phase0.Slot) (err error) {

	if !utils.Config.Graph.Enabled {
		return
	}
	var executionClient = c.indexer.executionPool.GetReadyEndpoint(execution.AnyClient)
	if executionClient == nil {
		return fmt.Errorf("processExecutionBlocks: could not get execution client")
	}

	var rawTransactions, er1 = executionClient.GetRPCClient().GetPendingTransactions(c.getContext())

	if er1 != nil {
		c.logger.Infof("Error on get GetPendingTransactions")
		return er1
	}

	//var j, er = json.Marshal(&res)
	var transactions []PendingTransaction
	err = json.Unmarshal(rawTransactions, &transactions)
	if err != nil {
		c.logger.Infof("Error on Unmarshal transactions")
		return err
	}
	var fileName = strconv.FormatUint(uint64(slot), 10)

	var workingDir = utils.Config.Graph.FilesPath
	/*err = os.WriteFile(workingDir+"/"+fileName, rawTransactions, 0644)
	if err != nil {
		c.logger.Infof("write file error: %v", err.Error())
		return err
	}*/
	var file, err1 = os.Create(workingDir + "/" + fileName + ".csv")
	if err1 != nil {
		c.logger.Infof("create file error: %v", err1.Error())
		return err1
	}
	defer file.Close()
	file.WriteString("hash,rank,from_address,to_address,gas_limit\n")
	for _, transaction := range transactions {
		file.WriteString(transaction.Hash + ",0," + transaction.From + "," + transaction.To + ",1000\n")
	}
	file.Sync()

	c.logger.Infof("pending transactions count: %v", len(transactions))

	if len(transactions) == 0 {
		removeOldFiles(c, workingDir)
		return
	}

	rClient, err := roger.NewRClient("127.0.0.1", 6311)
	if err != nil {
		c.logger.Infof("Failed to connect RServe")
		return
	}
	var session, err11 = rClient.GetSession()
	if err11 != nil {
		c.logger.Infof("Command failed: %v", err.Error())
		return err1
	}
	err11 = session.Assign("tx_df1", workingDir)
	if err11 != nil {
		c.logger.Infof("Command failed: %v", err.Error())
		return err1
	}
	err11 = session.Assign("tx_df2", fileName)
	if err11 != nil {
		c.logger.Infof("Command failed: %v", err.Error())
		return err1
	}
	c.logger.Infof("start eval")
	value, err := session.Eval("edge_list_plot_file(tx_df1, tx_df2)")
	if err != nil {
		c.logger.Infof("Command failed: %v", err.Error())
	} else {
		c.logger.Infof("result: %v", value)
	}
	c.logger.Infof("end eval")

	removeOldFiles(c, workingDir)
	return nil
}

func removeOldFiles(c *Client, folder string) {
	currentTime := time.Now()
	entries, err := os.ReadDir(folder)
	if err != nil {
		c.logger.Warnf("removeOldFiles failed: %v", err.Error())
	}

	for _, e := range entries {
		var fileInfo, err = e.Info()
		if err != nil {
			c.logger.Warnf("removeOldFiles failed: %v", err.Error())
		}
		var fileDate = fileInfo.ModTime()
		if currentTime.Sub(fileDate).Hours() > 24 {
			c.logger.Infof("remove result: %v %v", fileDate, e.Name())
			err = os.Remove(folder + "/" + e.Name())
			if err != nil {
				c.logger.Warnf("removeOldFiles failed: %v", err.Error())
			}
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
