package core

import "github.com/ethereum/go-ethereum/ethclient"

func NewConnection(url string) (*ethclient.Client, error) {
	return ethclient.Dial(url)
}
