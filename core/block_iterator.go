package core

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/atomic"
	"math/big"
	"sync"
)

type BlockIterator struct {
	client    *ethclient.Client
	wg        *sync.WaitGroup
	started   *atomic.Bool
	blockChan chan *types.Block
}

func NewBlockIterator(client *ethclient.Client) *BlockIterator {
	return &BlockIterator{client: client, wg: &sync.WaitGroup{}}
}

func (i *BlockIterator) Start(blockNum int64) <-chan *types.Block {
	i.blockChan = make(chan *types.Block)
	i.started = atomic.NewBool(true)
	i.wg.Add(1)
	go i.iterate(blockNum)
	return i.blockChan
}

func (i *BlockIterator) Stop() {
	i.started.Store(false)
	i.wg.Wait()
}

func (i *BlockIterator) iterate(blockNum int64) {
	blockNumber := blockNum
	retries := 0
	for i.started.Load() && retries < 16 {
		block, err := i.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if err != nil {
			retries++
			continue
		}

		if retries > 0 {
			// reset it to 0 if this was a successful attempt
			retries = 0
		}

		i.blockChan <- block
		blockNumber++
	}
	i.wg.Done()
}
