package db

import (
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/jmoiron/sqlx"
)

func GetTransactions(address string) []*dbtypes.Transaction {
	transactions := []*dbtypes.Transaction{}
	err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations
	FROM transactions
	WHERE "to" = $1 or "from"=$2
	ORDER BY block_number desc, block_rank desc
	`, address, address)
	if err != nil {
		logger.Errorf("Error while fetching Transactions: %v", err)
		return nil
	}
	return transactions
}

func GetTransactionByHash(hash string) *dbtypes.Transaction {
	transactions := []*dbtypes.Transaction{}
	err := ReaderDb.Select(&transactions, `
	SELECT
		"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations
	FROM transactions
	WHERE "hash" = $1
	`, hash)
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
			                          type
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			ON CONFLICT ("hash") DO NOTHING`,
		dbtypes.DBEngineSqlite: `
			INSERT OR IGNORE INTO transactions (
				"hash","block_number",block_rank,created_at,nonce,block_hash,transaction_index,"from","to",value,gas,gas_price,
			                          is_error,receipt_status,input,contract_address,cumulative_gas_used,gas_used,confirmations,
			                          type
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
	}),
		transaction.Hash, transaction.BlockNumber, transaction.BlockRank, transaction.TimeStamp, transaction.Nonce, transaction.BlockHash, transaction.TransactionIndex, transaction.From,
		transaction.To, transaction.Value, transaction.Gas, transaction.GasPrice, transaction.IsError, transaction.TxReceiptStatus, transaction.Input, transaction.ContractAddress, transaction.CumulativeGasUsed,
		transaction.GasUsed, transaction.Confirmations, transaction.Type)
	if err != nil {
		return err
	}
	return nil
}
