package db

import (
	"github.com/ethpandaops/dora/dbtypes"
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
