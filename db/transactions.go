package db

import (
	"fmt"
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/jmoiron/sqlx"
	"strings"
)

func GetTransactions(address string, start uint64, pageSize uint64) []*dbtypes.Transaction {
	upperAddress := strings.ToUpper(address)
	transactions := []*dbtypes.Transaction{}
	err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations, 
			                          erc20_method, erc20_address_to, erc20_value
	FROM transactions
	WHERE upper("to") = $1 or upper("from")=$2
	ORDER BY block_number desc, block_rank desc LIMIT $3 OFFSET $4
	`, upperAddress, upperAddress, pageSize, start)
	if err != nil {
		logger.Errorf("Error while fetching Transactions: %v", err)
		return nil
	}
	return transactions
}

func GetTransactionsErc20(address string, contract string, start uint64, pageSize uint64) []*dbtypes.Transaction {
	upperAddress := strings.ToUpper(address)
	upperContract := strings.ToUpper(contract)
	transactions := []*dbtypes.Transaction{}
	if contract != "" {
		err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations, 
			                          erc20_method, erc20_address_to, erc20_value
	FROM transactions
	WHERE upper("to") = $1 and (upper("from")=$2 or upper(erc20_address_to) =$3) and erc20_method !=''
	ORDER BY block_number desc, block_rank desc LIMIT $4 OFFSET $5
	`, upperContract, upperAddress, upperAddress, pageSize, start)
		if err != nil {
			logger.Errorf("Error while fetching Transactions: %v", err)
			return nil
		}
	} else {
		err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations, 
			                          erc20_method, erc20_address_to, erc20_value
	FROM transactions
	WHERE (upper("from")=$1 or upper(erc20_address_to) =$2) and erc20_method !=''
	ORDER BY block_number desc, block_rank desc LIMIT $3 OFFSET $4
	`, upperAddress, upperAddress, pageSize, start)
		if err != nil {
			logger.Errorf("Error while fetching Transactions: %v", err)
			return nil
		}
	}
	return transactions
}

func GetTransactionsCount(address string) (uint64, error) {
	upperAddress := strings.ToUpper(address)
	query := fmt.Sprintf("SELECT count(*)  FROM transactions WHERE upper(\"to\") = '%s' or upper(\"from\")='%s'", upperAddress, upperAddress)

	var count uint64
	err := ReaderDb.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count: %v", err)
	}
	return count, nil
}

func GetTransactionsErc20Count(address string, contract string) (uint64, error) {
	upperAddress := strings.ToUpper(address)
	upperContract := strings.ToUpper(contract)
	query := fmt.Sprintf("SELECT count(*)  FROM transactions WHERE (upper(\"from\")='%s' or upper(erc20_address_to)='%s') and erc20_method !=''",
		upperAddress, upperAddress)

	if contract != "" {
		query = fmt.Sprintf("SELECT count(*)  FROM transactions WHERE upper(\"to\") = '%s' and ( upper(\"from\")='%s' or upper(erc20_address_to)='%s') and erc20_method !=''",
			upperContract, upperAddress, upperAddress)
	}

	var count uint64
	err := ReaderDb.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count: %v", err)
	}
	return count, nil
}

func GetTransactionByHash(hash string) *dbtypes.Transaction {
	upperHash := strings.ToUpper(hash)
	transactions := []*dbtypes.Transaction{}
	err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations,
			                          erc20_method, erc20_address_to, erc20_value
	FROM transactions
	WHERE upper("hash") = $1
	`, upperHash)
	if err != nil {
		logger.Errorf("Error while fetching Transactions: %v", err)
		return nil
	}
	if len(transactions) > 0 {
		return transactions[0]
	}
	return nil
}

func InsertTransaction(transaction *dbtypes.Transaction, tx *sqlx.Tx) error {
	_, err := tx.Exec(EngineQuery(map[dbtypes.DBEngineType]string{
		dbtypes.DBEnginePgsql: `
			INSERT INTO transactions (
				"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations,
			                          type,erc20_method, erc20_address_to, erc20_value
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
			ON CONFLICT ("hash") 
DO UPDATE SET
    "block_number" = EXCLUDED."block_number",
    block_rank = EXCLUDED.block_rank,
    created_at = EXCLUDED.created_at,
    nonce = EXCLUDED.nonce,
    block_hash = EXCLUDED.block_hash,
    transaction_index = EXCLUDED.transaction_index,
    "from" = EXCLUDED."from",
    "to" = EXCLUDED."to",
    value = EXCLUDED.value,
    gas = EXCLUDED.gas,
    gas_price = EXCLUDED.gas_price,
    is_error = EXCLUDED.is_error,
    receipt_status = EXCLUDED.receipt_status,
    input = EXCLUDED.input,
    contract_address = EXCLUDED.contract_address,
    cumulative_gas_used = EXCLUDED.cumulative_gas_used,
    gas_used = EXCLUDED.gas_used,
    confirmations = EXCLUDED.confirmations,
    type = EXCLUDED.type,
    erc20_method = EXCLUDED.erc20_method, 
    erc20_address_to = EXCLUDED.erc20_address_to,
    erc20_value = EXCLUDED.erc20_value;`,
		dbtypes.DBEngineSqlite: `
			INSERT OR REPLACE INTO transactions (
    "hash", "block_number", block_rank, created_at, nonce, block_hash, transaction_index, "from", "to", value, gas,
    gas_price, is_error, receipt_status, input, contract_address, cumulative_gas_used, gas_used, confirmations, type
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20);`,
	}),
		transaction.Hash, transaction.BlockNumber, transaction.BlockRank, transaction.TimeStamp, transaction.Nonce, transaction.BlockHash, transaction.TransactionIndex, transaction.From,
		transaction.To, transaction.Value, transaction.Gas, transaction.GasPrice, transaction.IsError, transaction.TxReceiptStatus, transaction.Input, transaction.ContractAddress, transaction.CumulativeGasUsed,
		transaction.GasUsed, transaction.Confirmations, transaction.Type, transaction.Erc20Method, transaction.Erc20Address, transaction.Erc20Value)
	if err != nil {
		return err
	}
	return nil
}
