package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"gitlab.com/key-connect/geth-indexer/db"
	"gitlab.com/key-connect/geth-indexer/utils"
	"testing"
)

func Test_TokenTx(t *testing.T) {
	indexer := TxIndexer{}

	tests := []struct {
		name     string
		txHash   string
		inputHex string
	}{
		{
			name:     "usdc contract interaction",
			txHash:   "0x289d14a203d46a8d097fb5fe2ce505abbda716930b1ac2ced97a658f82c2c825",
			inputHex: "0xa9059cbb000000000000000000000000e78388b4ce79068e89bf8aa7f218ef6b9ab0e9d0000000000000000000000000000000000000000000000000000000174876fde4",
		}, {
			name:     "usdt contract interaction",
			txHash:   "0x70625e5daa306014f149300ea246c28b38048b9c849601e8d58727a15293ac1c",
			inputHex: "0xa9059cbb000000000000000000000000ed4f3d31b5d2d4f9f88617d4bf47d57725613d6e00000000000000000000000000000000000000000000000000000000183c1230",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenTx := db.TokenTx{}
			err := indexer.fillTokenTx(test.inputHex, &tokenTx)
			assert.NoError(t, err)

			to := test.inputHex[8 : 8+64] // 8:72
			toAddress := common.HexToAddress(to)
			value := test.inputHex[72 : 72+64]
			uint256Val, err := uint256.FromHex("0x" + utils.StripZeroPadding(value))
			assert.NoError(t, err)
			assert.Equal(t, toAddress.String(), tokenTx.To)
			assert.Equal(t, uint256Val.Uint64(), tokenTx.Value)
		})
	}
}

/*func Test_TokenTx(t *testing.T) {
	hexData := "a9059cbb000000000000000000000000e78388b4ce79068e89bf8aa7f218ef6b9ab0e9d0000000000000000000000000000000000000000000000000000000174876fde4"
	indexer := TxIndexer{}
	tokenTx := db.TokenTx{}
	err := indexer.fillTokenTx(hexData, &tokenTx)
	assert.NoError(t, err)

	to := hexData[8:8+64]	// 8:72
	toAddress := common.HexToAddress(to)
	value := hexData[72:72+64]
	uint256Val, err := uint256.FromHex("0x"+utils.StripZeroPadding(value))
	assert.NoError(t, err)
	assert.Equal(t, toAddress.String(), tokenTx.To)
	assert.Equal(t, uint256Val.Uint64(), tokenTx.Value)
}*/
