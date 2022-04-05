package indexer

import (
	"context"
	"time"

	"github.com/4406arthur/grape/jobqueue"
	"github.com/4406arthur/grape/repository"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type indexer struct {
	ctx          context.Context
	ethclient    *ethclient.Client
	jobQueue     jobqueue.JobQueue
	summaryRepo  repository.SummaryRepo
	timeInterval time.Duration
	//blockRepo
	//txRepo
}

func NewIndexer(ctx context.Context, ethclient *ethclient.Client, jobQueue jobqueue.JobQueue, summaryRepo repository.SummaryRepo, timeInterval time.Duration) *indexer {

	return &indexer{
		ctx:          ctx,
		ethclient:    ethclient,
		jobQueue:     jobQueue,
		timeInterval: timeInterval,
		summaryRepo:  summaryRepo,
	}
}

func (i *indexer) Start() {
	log.Info("start indexer")
	ticker := time.NewTicker(i.timeInterval)

	for {
		select {
		case <-ticker.C:
			log.Debug("discover new block...")
			// calculate how many block we need scrap
			latestBLockNum := i.summaryRepo.GetMaxBlockNum(i.ctx)
			// latestBLockNum := int64(18000000)
			currentLatest := i.getLatestBLockNum()
			log.Debug("rpc tell the latest block num is: ", currentLatest)
			for block := latestBLockNum + 1; block <= latestBLockNum+1000; block++ {
				//TODO: should be batch multi blocks

				if !i.jobQueue.Enqueue(block) {
					log.Warnf("job queue channel is full when deal with block: %d ", block)
					//something taffic jam happend
					//TODO: notify to scale worker or slow down producer
					time.Sleep(3 * time.Second)
					i.jobQueue.Enqueue(block)
				}
			}
		case <-i.ctx.Done():
			log.Warn("catch signal stop ticker now")
			ticker.Stop()
			return
		}
	}
}

func (i *indexer) getLatestBLockNum() int64 {
	header, err := i.ethclient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return header.Number.Int64()
}
