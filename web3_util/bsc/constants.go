package bsc

import (
	"go-util/web3_util/contract"

	"github.com/ethereum/go-ethereum/common"
)

var (
	Usdt2           = &contract.Contract{Name: "USDT", Symbol: "USDT", Address: common.HexToAddress("0xc2132d05d31c914a87c6611c10748aeb04b58e8f"), Decimals: 18, ChainName: "bsc"}
	Usdc            = &contract.Contract{Name: "USDC", Symbol: "USDC", Address: common.HexToAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174"), Decimals: 18, ChainName: "bsc"}
	PancakeFactory  = &contract.Contract{Name: "PancakeFactory", Symbol: "PancakeFactory", Address: common.HexToAddress("0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"), Decimals: 18, ChainName: "bsc"}
	PancakeRouter   = &contract.Contract{Name: "PancakeRouter", Symbol: "PancakeRouter", Address: common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"), Decimals: 18, ChainName: "bsc"}
	FistSwapFactory = &contract.Contract{Name: "FstSwapFactory", Symbol: "FstSwapFactory", Address: common.HexToAddress("0x9a272d734c5a0d7d84e0a892e891a553e8066dce"), Decimals: 18, ChainName: "bsc"}
	FistSwapRouter  = &contract.Contract{Name: "FstSwapRouter", Symbol: "FstSwapRouter", Address: common.HexToAddress("0x1b6c9c20693afde803b27f8782156c0f892abc2d"), Decimals: 18, ChainName: "bsc"}
	MultiCall       = &contract.Contract{Name: "MultiCall", Symbol: "MultiCall", Address: common.HexToAddress("0x5dc53ed77bbc84f39c76fb4c84ac9f28384a4b55"), Decimals: 18, ChainName: "bsc"}
	Wbnb            = &contract.Contract{Name: "WBNB", Symbol: "WBNB", Address: common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"), Decimals: 18, ChainName: "bsc"}
	Busd            = &contract.Contract{Name: "BUSD", Symbol: "BUSD", Address: common.HexToAddress("0xe9e7cea3dedca5984780bafc599bd69add087d56"), Decimals: 18, ChainName: "bsc"}
	Usdt            = &contract.Contract{Name: "USDT", Symbol: "USDT", Address: common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"), Decimals: 18, ChainName: "bsc"}
	WbnbBusdPair    = &contract.Contract{Name: "WbnbBusdPair", Symbol: "WBPair", Address: common.HexToAddress("0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16"), Decimals: 18, ChainName: "bsc"}
)
