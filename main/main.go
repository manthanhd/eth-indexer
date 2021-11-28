package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/viper"
	"gitlab.com/key-connect/geth-indexer/core"
	"gitlab.com/key-connect/geth-indexer/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	database, err := gorm.Open(sqlite.Open("index.database"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = database.AutoMigrate(&db.Block{}); err != nil {
		panic(err)
	}

	if err = database.AutoMigrate(&db.Tx{}); err != nil {
		panic(err)
	}
	if err = database.AutoMigrate(&db.TokenTx{}); err != nil {
		panic(err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}

	gethUrl := viper.GetString("eth.url")
	if len(gethUrl) == 0 {
		log.Error("geth_url must be defined")
		os.Exit(1)
	}

	var lastBlock db.Block
	result := database.Order("Number desc").Limit(1).Find(&lastBlock)
	if result.Error != nil && result.Error.Error() != "record not found" {
		panic(result.Error)
	}

	startingBlockNumber := int64(0)
	if lastBlock != (db.Block{}) {
		startingBlockNumber = int64(lastBlock.Number + 1)
	}

	startingBlockNumber = 13696939 - 50

	fmt.Println("starting block", startingBlockNumber)

	client, err := ethclient.Dial(gethUrl)
	if err != nil {
		panic(err)
	}

	blockIterator := core.NewBlockIterator(client)
	blockChan := blockIterator.Start(startingBlockNumber)

	contractsManager, err := core.NewContractsManager(client)
	if err != nil {
		os.Exit(1)
	}
	txIndexer := core.NewTxIndexer(blockChan, database, contractsManager)
	txIndexer.Start()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-c
		fmt.Println("stopping...")
		blockIterator.Stop()
		if err := txIndexer.Stop(); err != nil {
			log.Error("error stopping tx indexer", "err", err)
		}
		wg.Done()
		os.Exit(1)
	}()

	wg.Wait()
}
