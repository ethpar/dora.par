package beacon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethpandaops/dora/clients/execution"
	"github.com/ethpandaops/dora/db"
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/jmoiron/sqlx"
	dynssz "github.com/pk910/dynamic-ssz"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ExecutionBlock struct {
	Root     phase0.Root
	Slot     phase0.Slot
	Rank     uint64
	dynSsz   *dynssz.DynSsz
	Block    *types.Block
	blockRaw *json.RawMessage
}

func (block *ExecutionBlock) buildUnfinalizedBlock(compress bool) *dbtypes.UnfinalizedExecutionBlock {

	return &dbtypes.UnfinalizedExecutionBlock{
		Root:                block.Root[:],
		Slot:                uint64(block.Slot),
		Eth_block_hash:      block.Block.Hash().Bytes(),
		Eth_block_number:    block.Block.Number().Uint64(),
		EthTransactionCount: uint64(len(block.Block.Transactions())),
		Block:               *block.blockRaw,
		Status:              0,
		Rank:                block.Rank,
	}
}

func getExecutionHashes(block *Block) (hashes []string) {
	var mapHashes = map[string]string{}

	if block.block == nil {
		return hashes
	}

	var attestaions, err = block.block.Attestations()
	if err != nil {
		return hashes
	}
	for _, attestaion := range attestaions {
		var attestationData, err = attestaion.Data()
		if err != nil {
			continue
		}

		for _, executionHash := range attestationData.ExecutionHashes() {
			var s = executionHash.String()
			if s != "0x0000000000000000000000000000000000000000000000000000000000000000" {
				mapHashes[s] = s
			}
		}
		for _, s2 := range mapHashes {
			hashes = append(hashes, s2)
		}
	}
	return hashes
}

func processExecutionBlocks(c *Client, block *Block) (err error) {
	var executionClient = c.indexer.executionPool.GetReadyEndpoint(execution.AnyClient)
	if executionClient == nil {
		return fmt.Errorf("processExecutionBlocks: could not get execution client")
	}
	if block.block == nil {
		c.logger.Warn("processExecutionBlocks: block.block == nil")
		return
	}
	var blockNumber, err1 = block.block.ExecutionBlockNumber()
	if err1 != nil {
		c.logger.Errorf("processExecutionBlocks:  %v", err1)
		return err1
	}

	c.logger.Infof("check parallel Blocks: %v", blockNumber)
	for rank := 1; rank < 5; rank++ {
		c.logger.Debugf("check GetBlockByNumberAndRank: %v:%v", blockNumber, rank)
		parallelExecutionBlockRaw, err := executionClient.GetRPCClient().GetBlockByNumberAndRankRaw(c.getContext(), blockNumber, uint64(rank))
		if err != nil {
			if err.Error() != "not found" {
				c.logger.Errorf("GetBlockByNumberAndRank: %v:%v %v", blockNumber, rank, err)
			} else {
				c.logger.Infof("GetBlockByNumberAndRank not found: %v:%v", blockNumber, rank)
			}
			continue
		}
		parallelExecutionBlock, err := executionClient.GetRPCClient().DecodeBlockRaw(nil, parallelExecutionBlockRaw)
		if err == nil {
			processExecutionBlock(c, block, parallelExecutionBlock, parallelExecutionBlockRaw, uint64(rank))
		} else {
			c.logger.Errorf("DecodeBlockRaw: %v:%v %v", blockNumber, rank, err)
		}
	}
	return
}

func processExecutionBlocksWithHashes(c *Client, block *Block) (isExists bool, isNew bool, err error) {
	//var _blockNumber uint64
	//var _pBlock *Block
	var hashes = getExecutionHashes(block)
	if hashes != nil && len(hashes) > 0 {
		var rank uint64 = 0
		for _, s := range hashes {
			//c.logger.Info("parallel hash: %v ", s)

			var executionClient = c.indexer.executionPool.GetReadyEndpoint(execution.AnyClient)
			if executionClient == nil {
				return false, false, fmt.Errorf("could not get execution client")
			}

			parallelExecutionBlockRaw, err := executionClient.GetRPCClient().GetBlockByStringHashRaw(c.getContext(), s)
			parallelExecutionBlock, err := executionClient.GetRPCClient().DecodeBlockRaw(nil, parallelExecutionBlockRaw)

			if err == nil {
				processExecutionBlock(c, block, parallelExecutionBlock, parallelExecutionBlockRaw, uint64(rank))
			} else {
				c.logger.Errorf("!paralle error:  %v", err)
			}
		}
	}
	return
}

func processExecutionBlock(c *Client, block *Block, parallelExecutionBlock *types.Block,
	parallelExecutionBlockRaw *json.RawMessage, rank uint64) (isExists bool, isNew bool, err error) {
	if parallelExecutionBlock != nil {
		isExists = true
		//c.logger.Infof("parallel block %v %v %v", parallelExecutionBlock.Number(), block.Slot, rank) //json.
		isNew = true
		if isNew {
			var executionBlock = ExecutionBlock{
				Root:     block.Root,
				Slot:     block.Slot,
				Block:    parallelExecutionBlock,
				blockRaw: parallelExecutionBlockRaw,
				Rank:     rank,
			}

			block.ExecutionBlocks[rank] = executionBlock
			parallelDbBlock := executionBlock.buildUnfinalizedBlock(false)
			err = db.RunDBTransaction(func(tx *sqlx.Tx) error {
				err := db.InsertUnfinalizedExecutionBlock(parallelDbBlock, tx)
				if err != nil {
					c.logger.Errorf("!parallel block save error:  %v", err)
					return err
				}
				//c.logger.Debugf("saved execution block: slot: %v  %v:%v", block.Slot, parallelExecutionBlock.Number(), rank)
				c.logger.Infof("saved execution block: slot: %v  %v:%v", block.Slot, parallelExecutionBlock.Number(), rank)
				return nil
			})
		}
	}
	return
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}
type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}
type rpcBlock struct {
	Hash         common.Hash         `json:"hash"`
	Transactions []rpcTransaction    `json:"transactions"`
	UncleHashes  []common.Hash       `json:"uncles"`
	Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
	Requests     []*types.Request    `json:"requests,omitempty"`
}

func restoreExecutionBlocksFromDB(indexer *Indexer, block *Block) {
	for _, dbBlock := range db.GetExecutionBlocks(block.Root[:]) {
		var raw = json.RawMessage(dbBlock.Block)
		var parallelExecutionBlock, err = DecodeBlockRaw(raw)
		if err != nil {
			indexer.logger.Errorf("restoreExecutionBlocksFromDB %v", err)
		}
		var executionBlock = ExecutionBlock{
			Root:     phase0.Root(dbBlock.Root),
			Slot:     phase0.Slot(dbBlock.Slot),
			Block:    parallelExecutionBlock,
			blockRaw: &raw,
			Rank:     dbBlock.Rank,
		}
		block.ExecutionBlocks[dbBlock.Rank] = executionBlock
	}
}

func DecodeBlockRaw(raw json.RawMessage /*, ctx context.Context*/) (*types.Block, error) {
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, errors.New("not found")
	}

	var body rpcBlock
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	/*if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == types.EmptyTxsHash && len(body.Transactions) > 0 {
		return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != types.EmptyTxsHash && len(body.Transactions) == 0 {
		return nil, errors.New("server returned empty transaction list but block header indicates transactions")
	}*/

	var uncles []*types.Header
	/*if ctx != nil {
		if len(body.UncleHashes) > 0 {
			uncles = make([]*types.Header, len(body.UncleHashes))
			reqs := make([]rpc.BatchElem, len(body.UncleHashes))
			for i := range reqs {
				reqs[i] = rpc.BatchElem{
					Method: "eth_getUncleByBlockHashAndIndex",
					Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
					Result: &uncles[i],
				}
			}
			if err := ec.c.BatchCallContext(ctx, reqs); err != nil {
				return nil, err
			}
			for i := range reqs {
				if reqs[i].Error != nil {
					return nil, reqs[i].Error
				}
				if uncles[i] == nil {
					return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
				}
			}
		}
	}*/
	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if tx.From != nil {
			//setSenderFromServer(tx.tx, *tx.From, body.Hash)
		}
		txs[i] = tx.tx
	}
	return types.NewBlockWithHeader(head).WithBody(
		types.Body{
			Transactions: txs,
			Uncles:       uncles,
			Withdrawals:  body.Withdrawals,
			Requests:     body.Requests,
		}), nil
}

type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

/*func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}*/
