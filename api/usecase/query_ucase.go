package usecase

import (
	"context"

	"github.com/4406arthur/grape/repository"

	log "github.com/sirupsen/logrus"
)

type QueryUsecase struct {
	summaryRepo repository.SummaryRepo
	blockRepo   repository.BlockRepo
	txRepo      repository.TxRepo
}

func NewQueryUsecase(summaryRepo repository.SummaryRepo, blockRepo repository.BlockRepo, txRepo repository.TxRepo) *QueryUsecase {
	return &QueryUsecase{
		summaryRepo: summaryRepo,
		blockRepo:   blockRepo,
		txRepo:      txRepo,
	}
}

func (u *QueryUsecase) GetLatestNBlock(c context.Context, n int) ([]repository.Block, error) {
	res := make([]repository.Block, 0)
	latestBlockNum := u.summaryRepo.GetMaxBlockNum(c)
	//log.Infof("current latestBlockNum is: %d", latestBlockNum)
	lowFlag := latestBlockNum - int64(n)

	for i := latestBlockNum; i > lowFlag; i-- {
		b, err := u.blockRepo.GetBlockByID(c, i)
		// log.Debugf("ASDfas: %+v", b)
		if err != nil {
			return res, err
		}
		res = append(res, b)
	}

	return res, nil
}

func (u *QueryUsecase) GetFullBlockByID(c context.Context, id int) (repository.Block, []repository.Tx, error) {
	b, err := u.blockRepo.GetBlockByID(c, int64(id))
	// log.Debugf("ASDfas: %+v", b)
	if err != nil {
		return repository.Block{}, []repository.Tx{}, err
	}
	txs, err := u.txRepo.GetTXsByBlockNum(c, int64(id))
	if err != nil {
		return repository.Block{}, []repository.Tx{}, err
	}

	return b, txs, nil
}

func (u *QueryUsecase) GetTXByHash(c context.Context, hash string) (repository.Tx, error) {
	tx, err := u.txRepo.GetTXByHash(c, hash)
	log.Debugf("XNVXB12: %s %v", hash, tx)
	if err != nil {
		return repository.Tx{}, err
	}

	return tx, nil
}
