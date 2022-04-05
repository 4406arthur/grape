package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/4406arthur/grape/api/http"
	"github.com/4406arthur/grape/api/usecase"
	"github.com/4406arthur/grape/indexer"
	jobqueue "github.com/4406arthur/grape/jobqueue/memory"
	repo "github.com/4406arthur/grape/repository/postgres"
	"github.com/4406arthur/grape/scraper"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.SetLevel(log.DebugLevel)
		log.Info("Service RUN on DEBUG mode")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	//DB config setup
	host := viper.GetString(`database.host`)
	port := viper.GetString(`database.port`)
	user := viper.GetString(`database.user`)
	pass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
	// open database
	dbConn, err := sql.Open(`postgres`, conn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	ctx := context.Background()
	jq := jobqueue.NewMemoryQueue(viper.GetInt(`jobqueue.size`))
	client, err := ethclient.Dial(viper.GetString(`ethrpc.endpoint`))
	if err != nil {
		log.Fatal(err)
	}
	pollingSecInterval := time.Second * time.Duration(viper.GetInt(`ethrpc.pollingInterval`))
	summaryRepo := repo.NewSummaryRepository(dbConn)
	indexer := indexer.NewIndexer(ctx, client, jq, summaryRepo, pollingSecInterval)
	blkRepo := repo.NewBlockRepository(dbConn)
	txRepo := repo.NewTxRepository(dbConn)
	go indexer.Start()
	s := scraper.NewScraper(ctx, client, blkRepo, jq.Subscribe(), viper.GetInt(`ethrpc.scraperWorker`))
	s.Start()

	e := echo.New()
	q := usecase.NewQueryUsecase(summaryRepo, blkRepo, txRepo)
	http.NewQueryHandler(e, q)

	// go waitSignal(cancel)
	log.Fatal(e.Start(viper.GetString("server.address")))
}

// func waitSignal(cancel context.CancelFunc) {
// 	ch := make(chan os.Signal, 1)
// 	signal.Notify(
// 		ch,
// 		syscall.SIGINT,
// 		syscall.SIGQUIT,
// 		syscall.SIGTERM,
// 	)
// 	for {
// 		sig := <-ch
// 		switch sig {
// 		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
// 			cancel()
// 			log.Warn("GET: ", sig)
// 			time.Sleep(2 * time.Second)
// 			return
// 		}
// 	}
// }
