package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/4406arthur/grape/repository"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type BlockRepository struct {
	DB *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{
		DB: db,
	}
}

func (db *BlockRepository) GetBlockByID(ctx context.Context, number int64) (repository.Block, error) {
	var block repository.Block
	err := db.DB.QueryRowContext(ctx, "SELECT * FROM blk WHERE block_num = $1", number).Scan(&block.BlockNum, &block.BlockHash, &block.ParentHash, &block.BlockTime)
	if err != nil {
		//TODO: should push this lost block to job queue
		// or make a ETH-RPC query then insert
		return block, err
	}
	return block, nil
}

func (db *BlockRepository) GetBlockByHash(ctx context.Context, hash string) (repository.Block, error) {
	var block repository.Block
	err := db.DB.QueryRowContext(ctx, "SELECT * FROM blk WHERE block_hash = $1", hash).Scan(&block.BlockNum, &block.BlockHash, &block.ParentHash, &block.BlockTime)
	if err != nil {
		return block, err
	}
	return block, nil
}

//func (db *BlockRepository) GetLatestNBlocks(int) []repository.Block { return []repository.Block{} }

func (db *BlockRepository) InsertFullBlockInfo(ctx context.Context, blk repository.Block, txList []repository.Tx) error {
	if db.checkBlockExist(blk.BlockNum) {
		return nil
	}
	start := time.Now()
	defer func() {
		log.WithField("time duration: ", time.Since(start)).WithField("count", len(txList)).Debugf("block: %d processing done", blk.BlockNum)
	}()

	t, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("caonnot start a transaction when deal with block: %d", blk.BlockNum)
	}

	_, err = t.ExecContext(ctx, "INSERT INTO blk (block_num, block_hash, parent_hash, block_time) VALUES ($1, $2, $3, $4)", blk.BlockNum, blk.BlockHash, blk.ParentHash, blk.BlockTime)
	if err != nil {
		log.Errorf("FSB12bX: %+v", err.Error())
		return err
	}

	stmt, err := t.Prepare(pq.CopyIn("tx", "tx_hash", "block_num", "tx_from", "tx_to", "tx_nonce", "tx_data", "tx_value"))
	if err != nil {
		log.Errorf("prepare setup failed: %+v", err.Error())
		return err
	}
	defer stmt.Close()

	for _, tx := range txList {
		// log.Debugf("SUAF1231: %s", tx.TxHash)
		_, err = stmt.Exec(tx.TxHash, blk.BlockNum, tx.TxFrom, tx.TxTo, tx.TxNonce, tx.TxData, tx.TxValue)
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	err = t.Commit()
	if err != nil {
		t.Rollback()
		log.Errorf("caonnot commit a transaction when deal with block: %d", blk.BlockNum)
		return err
	}
	return nil
}

func (db *BlockRepository) checkBlockExist(blockNum int64) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT exists (SELECT 1 FROM blk WHERE block_num = $1)", blockNum).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
