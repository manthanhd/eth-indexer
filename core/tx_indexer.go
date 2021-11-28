package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"gitlab.com/key-connect/geth-indexer/db"
	"gitlab.com/key-connect/geth-indexer/utils"
	"go.uber.org/atomic"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"time"
)

type Indexer interface {
	Start() error
	Stop() error
}

type TxIndexer struct {
	blockChan           <-chan *types.Block
	db                  *gorm.DB
	txSignerMap         map[uint64]types.Signer
	started             *atomic.Bool
	wg                  *sync.WaitGroup
	blocksIndexed       uint64
	transactionsIndexed uint64
	contractsManager    *ContractsManager
}

func NewTxIndexer(blockChan <-chan *types.Block, db *gorm.DB, contractsManager *ContractsManager) *TxIndexer {
	return &TxIndexer{
		blockChan:        blockChan,
		db:               db,
		txSignerMap:      make(map[uint64]types.Signer, 3),
		wg:               &sync.WaitGroup{},
		contractsManager: contractsManager,
	}
}

func (t *TxIndexer) Start() {
	go t.startAsync()
}

func (t *TxIndexer) StartBlocking() {
	t.startAsync()
}

func (t *TxIndexer) startAsync() {
	dbSession := t.db.Session(&gorm.Session{CreateBatchSize: 1024})
	defer dbSession.Commit()
	lastUpdate := time.Now()
	t.started = atomic.NewBool(true)
	t.wg.Add(1)
	for t.started.Load() {
		ethBlock, open := <-t.blockChan
		if !open {
			fmt.Println("Indexer stopping, channel was closed", "blocksIndexed", t.blocksIndexed, "transactionsIndexed", t.transactionsIndexed)
			break
		}

		block := db.Block{
			Number:     ethBlock.NumberU64(),
			Hash:       ethBlock.Hash().String(),
			ParentHash: ethBlock.ParentHash().String(),
			Timestamp:  time.Unix(int64(ethBlock.Time()), 0),
			GasUsed:    ethBlock.GasUsed(),
			GasLimit:   ethBlock.GasLimit(),
		}

		lenTransactions := len(ethBlock.Transactions())
		if lenTransactions == 0 {
			fmt.Println("block had no transactions", "blockNumber", block.Number, "blockHash", block.Hash)
			dbSession.Create(&block)
			t.blocksIndexed++
			continue
		}

		transactions := make([]db.Tx, 0, lenTransactions)
		tokenTransactions := make([]db.TokenTx, 0, lenTransactions)

		for _, ethTx := range ethBlock.Transactions() {
			from, err := t.GetSenderAddress(block.Number, ethTx)
			if err != nil {
				// shut down the node so that it can resume
				fmt.Println("error reading sender address of transaction", "txHash", ethTx.Hash(), "err", err)
				return
			}

			inputData := ethTx.Data()
			hexData := hex.EncodeToString(inputData)
			txType := t.getTransactionType(hexData)
			fmt.Println("processing tx", ethTx.Hash())
			fmt.Println("from", from, "nonce", ethTx.Nonce())
			to := ethTx.To()
			if to == nil {
				// this is a contract creation transaction
				/*b, err := rlp.EncodeToBytes([]string{from.String(), fmt.Sprintf("%d", ethTx.Nonce())})
				if err != nil {
					fmt.Println("err", err)
					return
				}
				hash := crypto.Keccak256(b)
				toAddress := common.BytesToAddress(hash)
				fmt.Println("toAddress", toAddress)
				to = &toAddress
				return*/

				// skip contract creation transaction#
				fmt.Println("skipping contract creation transaction", "hash", ethTx.Hash(), "blockNumber", block.Number)
				continue
			}
			// todo fix broken for contract creation transaction 0x8100e5624f5d12367af7245c37091d8091d474d2b1a283703b4be61e41791c22
			tx := db.Tx{
				Hash:         ethTx.Hash().String(),
				BlockNumber:  block.Number,
				Timestamp:    block.Timestamp,
				From:         from.String(),
				To:           to.String(),
				Value:        ethTx.Value().Uint64(),
				EthTxType:    ethTx.Type(),
				TxType:       txType,
				GasUsed:      ethTx.Gas(),
				GasLimit:     ethTx.GasFeeCap().Uint64(),
				InputDataHex: hexData,
			}

			transactions = append(transactions, tx)

			if len(hexData) >= 8 {
				// this tx has a contract interaction
				fmt.Println("contract interaction", "tx", ethTx.Hash(), "data", hexData)
				switch txType {
				case db.TokenTxType:
					tokenTx := db.TokenTx{
						TxHash:          ethTx.Hash().String(),
						BlockNumber:     block.Number,
						Timestamp:       block.Timestamp,
						From:            from.String(),
						ContractAddress: to.String(),
						GasUsed:         ethTx.Gas(),
						GasLimit:        ethTx.GasFeeCap().Uint64(),
					}

					if err = t.fillTokenTx(hexData, &tokenTx); err != nil {
						fmt.Println("error parsing token transaction", "hash", ethTx.Hash(), "err", err)
						return
					}

					tokenContract, err := t.contractsManager.GetTokenContract(*ethTx.To())
					if err != nil {
						fmt.Println("error getting token contract", "hash", ethTx.Hash(), "err", err)
						return
					}

					tokenTx.Symbol = tokenContract.Symbol
					tokenTx.Decimals = tokenContract.Decimals

					tokenTransactions = append(tokenTransactions, tokenTx)
				}
			}
		}
		dbSession.Create(&block)
		dbSession.Create(transactions)
		dbSession.Create(tokenTransactions)

		t.blocksIndexed++
		t.transactionsIndexed += uint64(lenTransactions)

		if time.Since(lastUpdate) > 30*time.Second {
			fmt.Println("Transaction indexing progress", "blocksIndexed", t.blocksIndexed, "transactionsIndexed", t.transactionsIndexed, "lastBlockNum", block.Number)
			lastUpdate = time.Now()
		}
	}
	t.wg.Done()
}

func (t *TxIndexer) GetSenderAddress(blockNumber uint64, tx *types.Transaction) (common.Address, error) {
	signer := types.MakeSigner(params.MainnetChainConfig, big.NewInt(int64(blockNumber)))
	return signer.Sender(tx)
}

func (t *TxIndexer) Stop() error {
	t.started.Store(false)
	t.wg.Wait()
	return nil
}

// assumes len is greater than 8
func (t *TxIndexer) getTransactionType(hexData string) db.TxType {
	if len(hexData) >= 10 {
		fnName := hexData[2:10]
		switch fnName {
		case "a9059cbb":
			if len(hexData) == 138 {
				return db.TokenTxType
			}
		}
	}
	return db.UnknownTxType
}

// assumes this is a token tx
// fills to and value fields in given db.TokenTx object
// returns error if something went wrong
func (t *TxIndexer) fillTokenTx(data string, tokenTx *db.TokenTx) error {
	if len(data) != 138 {
		return errors.New("input data not expected length")
	}

	to := data[10 : 10+64] // 8:72
	toAddress := common.HexToAddress(to)
	value := data[74 : 74+64]
	value = utils.StripZeroPadding(value)
	val, err := uint256.FromHex("0x" + value) // remember to remove padding
	if err != nil {
		return err
	}
	tokenTx.To = toAddress.String()
	tokenTx.Value = val.Uint64()
	return nil
}
