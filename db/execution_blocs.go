package db

import (
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/jmoiron/sqlx"
)

func GetExecutionBlocks(root []byte) []*dbtypes.UnfinalizedExecutionBlock {
	blocks := []*dbtypes.UnfinalizedExecutionBlock{}
	err := ReaderDb.Select(&blocks, `
	SELECT
		slot, rank, root, eth_block_number, eth_block_hash, eth_transaction_count, block
	FROM unfinalized_execution_blocks
	WHERE root = $1
	ORDER BY rank desc
	`, root)
	if err != nil {
		logger.Errorf("Error while fetching ExecutionBlocks: %v", err)
		return nil
	}
	return blocks
}

func InsertUnfinalizedExecutionBlock(block *dbtypes.UnfinalizedExecutionBlock, tx *sqlx.Tx) error {
	_, err := tx.Exec(EngineQuery(map[dbtypes.DBEngineType]string{
		dbtypes.DBEnginePgsql: `
			INSERT INTO unfinalized_execution_blocks (
				root, slot, rank,block,   eth_block_number, eth_block_hash,eth_transaction_count
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (slot,rank) DO NOTHING`,
		dbtypes.DBEngineSqlite: `
			INSERT OR IGNORE INTO unfinalized_execution_blocks (
				root, slot, header_ver, header_ssz, block_ver, block_ssz, status, fork_id, rank
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
	}),
		block.Root, block.Slot, block.Rank, block.Block, block.Eth_block_number, block.Eth_block_hash, block.EthTransactionCount)
	if err != nil {
		return err
	}
	return nil
}
