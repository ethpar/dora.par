package dbtypes

import "math/big"

type IndexerSyncState struct {
	Epoch uint64 `json:"epoch"`
}

type IndexerPruneState struct {
	Epoch uint64 `json:"epoch"`
}

type IndexerTxState struct {
	BlockNumberStart big.Int `json:"block_number_start"`
	BlockNumberEnd   big.Int `json:"block_number_end"`
}

type IndexerForkState struct {
	ForkId    uint64 `json:"fork_id"`
	Finalized uint64 `json:"finalized"`
}

type DepositIndexerState struct {
	FinalBlock   uint64 `json:"final_block"`
	HeadBlock    uint64 `json:"head_block"`
	DepositIndex uint64 `json:"deposit_index"`
}
