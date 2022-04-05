package postgres

import (
	"context"
	"database/sql"

	"github.com/4406arthur/grape/repository"
)

type TxRepository struct {
	DB *sql.DB
}

func NewTxRepository(db *sql.DB) *TxRepository {
	return &TxRepository{
		DB: db,
	}
}

func (db *TxRepository) GetTXByHash(ctx context.Context, hash string) (repository.Tx, error) {
	var tx repository.Tx
	err := db.DB.QueryRowContext(ctx, "SELECT tx_hash,block_num,tx_from,tx_to,tx_nonce,tx_data,tx_value FROM tx WHERE tx_hash = $1", hash).Scan(
		&tx.TxHash, &tx.BlockNum, &tx.TxFrom, &tx.TxTo, &tx.TxNonce, &tx.TxData, &tx.TxValue)
	if err != nil {
		//TODO: should push this lost block to job queue
		// or make a ETH-RPC query then insert
		return tx, err
	}
	return tx, nil
}

func (db *TxRepository) GetTXsByBlockNum(ctx context.Context, id int64) ([]repository.Tx, error) {

	stmt, _ := db.DB.PrepareContext(ctx, "SELECT tx_hash,block_num,tx_from,tx_to,tx_nonce,tx_data,tx_value FROM tx WHERE block_num = $1")
	rows, _ := stmt.Query(id)
	txs := make([]repository.Tx, 0)
	for rows.Next() {
		var tx repository.Tx
		if err := rows.Scan(&tx.TxHash, &tx.BlockNum, &tx.TxFrom, &tx.TxTo, &tx.TxNonce, &tx.TxData, &tx.TxValue); err != nil {
			return []repository.Tx{}, err
		}
		txs = append(txs, tx)
	}
	if err := rows.Err(); err != nil {
		return []repository.Tx{}, err
	}
	return txs, nil
}
