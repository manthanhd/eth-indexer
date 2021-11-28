package core

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
	"gitlab.com/key-connect/geth-indexer/erc20"
	"sync"
)

type TokenContract struct {
	Symbol   string
	Decimals uint8
}

type ContractsManager struct {
	client *ethclient.Client
	lock   sync.Mutex
	// contractsCache caches contract address => contract code
	contractsCache *lru.Cache
	// tokenContractCache caches contract address => TokenContract
	tokenContractCache *lru.Cache
}

func NewContractsManager(client *ethclient.Client) (*ContractsManager, error) {
	contractsCache, err := lru.New(256)
	if err != nil {
		return nil, err
	}
	tokenContractCache, err := lru.New(256)
	if err != nil {
		return nil, err
	}
	return &ContractsManager{
		client:             client,
		contractsCache:     contractsCache,
		tokenContractCache: tokenContractCache,
	}, nil
}

func (c *ContractsManager) GetCode(contractAddress common.Address) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	contractCode, ok := c.contractsCache.Get(contractAddress.String())
	if !ok {
		code, err := c.client.CodeAt(context.Background(), contractAddress, nil)
		if err != nil {
			return "", err
		}
		contractCode = hex.EncodeToString(code)
		c.contractsCache.Add(contractAddress.String(), contractCode)
	}
	return contractCode.(string), nil
}

func (c *ContractsManager) GetTokenContract(contractAddress common.Address) (*TokenContract, error) {
	var tokenContract TokenContract
	tokenContractIntf, ok := c.tokenContractCache.Get(contractAddress.String())
	if ok {
		tokenContract = tokenContractIntf.(TokenContract)
	} else {
		// extract additional metadata if its a token contract
		erc20, err := erc20.NewErc20(contractAddress, c.client)
		if err == nil && erc20 != nil {
			symbol, err := erc20.Symbol(&bind.CallOpts{})
			if err != nil {
				return nil, errors.New("contract not erc20")
			}

			decimal, err := erc20.Decimals(&bind.CallOpts{})
			if err != nil {
				return nil, errors.New("contract not erc20")
			}

			tokenContract = TokenContract{
				Symbol:   symbol,
				Decimals: decimal,
			}
			c.tokenContractCache.Add(contractAddress.String(), tokenContract)
		}
	}
	return &tokenContract, nil
}
