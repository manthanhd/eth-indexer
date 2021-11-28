package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TokenContractMetadata(t *testing.T) {
	// usdc contract => 0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48	NOT AN ERC20 TOKEN APPARENTLY
	// usdt contract => 0xdac17f958d2ee523a2206206994597c13d831ec7
	// amp contract => 0xff20817765cb7f73d4bde2e66e067e58d11095c2
	url := ""
	client, err := ethclient.Dial(url)
	assert.NoError(t, err)
	contractAddress := common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	contractsManager, err := NewContractsManager(client)
	assert.NoError(t, err)
	contract, err := contractsManager.GetTokenContract(contractAddress)
	assert.NoError(t, err)
	assert.Equal(t, "USDT", contract.Symbol)
	assert.Equal(t, uint8(6), contract.Decimals)
	fmt.Println(contract)
	/*code, err := client.CodeAt(context.Background(), contractAddress, nil)
		assert.NoError(t, err)
		contractCode := hex.EncodeToString(code)
		fmt.Println("contract code", contractCode)
	:wq
		file, err := os.Open("")
		assert.NoError(t, err)
		abi, err := abi2.JSON(file)
		assert.NoError(t, err)
		fmt.Println("decimals function =", abi.Methods["decimals"])
		erc20, err := erc20.NewErc20(contractAddress, client)
		assert.NoError(t, err)
		decimals, err := erc20.Decimals(&bind.CallOpts{})
		assert.NoError(t, err)
		fmt.Println("decimals", decimals)*/
	/*
		// THIS IS ONE WAY TO CALL AN ERC20 TOKEN
		jsonrpcClient, err := jsonrpc.NewClient("https://mainnet.infura.io/v3/")
		assert.NoError(t, err)
		erc20 := erc20.NewERC20(web3.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), jsonrpcClient)
		decimals, err := erc20.Decimals()
		assert.NoError(t, err)
		fmt.Println("decimals", decimals)*/
}
