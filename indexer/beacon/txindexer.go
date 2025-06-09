package beacon

import (
	"context"
	"github.com/ethpandaops/dora/clients/consensus"
	"github.com/ethpandaops/dora/clients/execution"
	"github.com/ethpandaops/dora/db"
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math/big"
	"runtime/debug"
	"time"
)

type TxIndexer struct {
	logger        logrus.FieldLogger
	consensusPool *consensus.Pool
	executionPool *execution.Pool

	running bool
}

func NewTxIndexer(logger logrus.FieldLogger, consensusPool *consensus.Pool, executionPool *execution.Pool) *TxIndexer {

	// Create the indexer instance.
	indexer := &TxIndexer{
		logger:        logger,
		consensusPool: consensusPool,
		executionPool: executionPool,
	}

	return indexer
}

func (indexer *TxIndexer) StartIndexer() {
	if indexer.running {
		return
	}

	indexer.running = true

	go indexer.runIndexerLoop()
}

func (indexer *TxIndexer) runIndexerLoop() {
	defer func() {
		if err := recover(); err != nil {
			indexer.logger.WithError(err.(error)).Errorf("uncaught panic in indexer.beacon.TxIndexer.runIndexerLoop subroutine: %v, stack: %v", err, string(debug.Stack()))
			time.Sleep(10 * time.Second)

			go indexer.runIndexerLoop()
		}
	}()

	//chainState := indexer.consensusPool.GetChainState()
	txState := dbtypes.IndexerTxState{}
	_, err := db.GetExplorerState("indexer.tx", &txState)
	if err != nil {
		indexer.logger.Errorf("fails to read indexer.tx: %v", err)
		return
	}

	indexer.logger.Infof("start with blockNumber:%v", txState.BlockNumberStart)
	blockNumber := txState.BlockNumberStart
	blockNumberEnd := txState.BlockNumberEnd
	b := big.NewInt(1)
	indexer.logger.Infof("blockNumber:%v", blockNumber)
	for {

		clients := indexer.executionPool.GetAllEndpoints()
		if len(clients) == 0 {

			continue
		}

		indexer.logger.Infof("!!!!!!process blockNumber:%v", blockNumber)
		client := clients[0].GetRPCClient()
		executionClient := client.GetEthClient()

		for i := 0; i < 50; i++ {
			//indexer.logger.Infof("!!!process blockNumber:%v %v", blockNumber, i)
			UpdateTransactionsForBlock(context.Background(), executionClient, &blockNumber, indexer.logger)
			blockNumber.Sub(&blockNumber, b)
			if blockNumber.Cmp(&blockNumberEnd) < 0 {
				indexer.logger.Infof("indexed, exit:%v %v", blockNumber, blockNumberEnd)
				indexer.storeState(blockNumber, blockNumberEnd)
				return
			}
		}
		time.Sleep(3 * time.Second)
		indexer.logger.Infof("!!!processed blockNumber:%v", blockNumber)
		indexer.storeState(blockNumber, blockNumberEnd)
	}

}

func (indexer *TxIndexer) storeState(blockNumber big.Int, blockNumberEnd big.Int) {
	err := db.RunDBTransaction(func(tx *sqlx.Tx) error {
		syncState := &dbtypes.IndexerTxState{
			BlockNumberStart: blockNumber,
			BlockNumberEnd:   blockNumberEnd,
		}
		return db.SetExplorerState("indexer.tx", syncState, tx)
	})
	if err != nil {
		indexer.logger.Errorf("save state:%v", err)
	}

}
