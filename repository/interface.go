package repository

import "context"

type BlockRepo interface {
	GetBlockByID(context.Context, int64) (Block, error)
	GetBlockByHash(context.Context, string) (Block, error)
	//GetLatestNBlocks(int) ([]Block, error)
	InsertFullBlockInfo(context.Context, Block, []Tx) error
}

type TxRepo interface {
	GetTXByHash(context.Context, string) (Tx, error)
	GetTXsByBlockNum(context.Context, int64) ([]Tx, error)
}

type SummaryRepo interface {
	GetMaxBlockNum(context.Context) int64
}
