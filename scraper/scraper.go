package scraper

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/4406arthur/grape/repository"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type Scraper struct {
	ctx         context.Context
	ethclient   *ethclient.Client
	blockRepo   repository.BlockRepo
	jobChannel  chan int64
	concurrency int
}

func NewScraper(ctx context.Context, ethclient *ethclient.Client, repo repository.BlockRepo, jobChannel chan int64, concurrency int) *Scraper {
	return &Scraper{
		ctx:         ctx,
		ethclient:   ethclient,
		blockRepo:   repo,
		jobChannel:  jobChannel,
		concurrency: concurrency,
	}
}

func (s *Scraper) Start() {
	for workerID := 1; workerID <= s.concurrency; workerID++ {
		go s.executor(workerID)
	}
}

func (s *Scraper) executor(id int) {
	for blkID := range s.jobChannel {
		log.Debug("executor id: ", id, " current block_num: ", blkID)
		//leverage ETH RPC to fetch Block info
		blk, tx, _ := s.getNBlock(blkID)
		err := s.blockRepo.InsertFullBlockInfo(s.ctx, blk, tx)
		if err != nil {
			//TODO: should push to a persist storage for retry
			log.Error(err.Error())
		}
	}
	// for {
	// 	select {
	// 	case <-s.jobChannel:
	// 		blkID := <-s.jobChannel
	// 		log.Debug("executor id: ", id, " current block_num: ", blkID)
	// 		//leverage ETH RPC to fetch Block info
	// 		// blk, tx, _ := s.getNBlock(blkID)
	// 		// log.Debugf("executor id: %d, print block info %+v %+v", id, blk, tx)
	// 		// err := s.blockRepo.InsertFullBlockInfo(s.ctx, blk, tx)
	// 		// if err != nil {
	// 		// 	//TODO: should push to a persist storage for retry
	// 		// 	log.Error(err.Error())
	// 		// }
	// 	case <-s.ctx.Done():
	// 		log.Warnf("catch signal scraper worker: %d stop now", id)
	// 		return
	// 	}
	// }
}

func (s *Scraper) getNBlock(n int64) (repository.Block, []repository.Tx, error) {

	block, err := s.ethclient.BlockByNumber(s.ctx, big.NewInt(n))
	if err != nil {
		//TODO: throw into retry queue or persist into a failed block hostory table
		log.Warnf("error getting block: %d with error: %s", n, err.Error())
		return repository.Block{}, nil, err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	tx := make([]repository.Tx, 0)
	for _, v := range block.Transactions() {
		// log.Debugf("tx id: %s", v.Hash().String())
		wg.Add(1)
		tmpTx := v
		go func() {
			defer wg.Done()

			chainID, _ := s.ethclient.NetworkID(s.ctx)
			msg, _ := tmpTx.AsMessage(types.NewEIP155Signer(chainID), nil)
			//For contract-creation transactions, To returns empty string, avoid segment fault
			var to string
			if msg.To() == nil {
				to = ""
			} else {
				to = msg.To().String()
			}

			tmp := repository.Tx{
				TxHash:   tmpTx.Hash().String(),
				BlockNum: n,
				TxFrom:   msg.From().String(),
				TxTo:     to,
				TxNonce:  tmpTx.Nonce(),
				TxData:   hex.EncodeToString(msg.Data()),
				TxValue:  tmpTx.Value().String(),
			}

			mu.Lock()
			tx = append(tx, tmp)
			mu.Unlock()
		}()
	}
	wg.Wait()

	blk := repository.Block{
		BlockNum:   block.Header().Number.Int64(),
		BlockHash:  block.Hash().String(),
		ParentHash: block.Header().ParentHash.String(),
		BlockTime:  block.Header().Time,
	}
	return blk, tx, nil
}
