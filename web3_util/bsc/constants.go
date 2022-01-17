package bsc

import (
	"tugou-hunter/util/web3_util/contract"

	"github.com/panyanyany/go-web3"
)

var (
	Usdt = &contract.Contract{Name: "USDT", Symbol: "USDT", Address: web3.HexToAddress("0xc2132d05d31c914a87c6611c10748aeb04b58e8f"), Decimals: 18, ChainName: "bsc"}
	Usdc           = &contract.Contract{Name: "USDC", Symbol: "USDC", Address: web3.HexToAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174"), Decimals: 18, ChainName: "bsc"}
	PancakeFactory = &contract.Contract{Name: "PancakeFactory", Symbol: "PancakeFactory", Address: web3.HexToAddress("0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"), Decimals: 18, ChainName: "bsc"}
	MultiCall      = &contract.Contract{Name: "MultiCall", Symbol: "MultiCall", Address: web3.HexToAddress("0x5dc53ed77bbc84f39c76fb4c84ac9f28384a4b55"), Decimals: 18, ChainName: "bsc"}
	PancakeRouter  = &contract.Contract{Name: "PancakeRouter", Symbol: "PancakeRouter", Address: web3.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"), Decimals: 18, ChainName: "bsc"}
	Wbnb           = &contract.Contract{Name: "WBNB", Symbol: "WBNB", Address: web3.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"), Decimals: 18, ChainName: "bsc"}
	Busd           = &contract.Contract{Name: "BUSD", Symbol: "BUSD", Address: web3.HexToAddress("0xe9e7cea3dedca5984780bafc599bd69add087d56"), Decimals: 18, ChainName: "bsc"}
	BscUsd         = &contract.Contract{Name: "BUSD", Symbol: "BUSD", Address: web3.HexToAddress("0x55d398326f99059fF775485246999027B3197955"), Decimals: 18, ChainName: "bsc"}
	WbnbBusdPair   = &contract.Contract{Name: "WbnbBusdPair", Symbol: "WBPair", Address: web3.HexToAddress("0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16"), Decimals: 18, ChainName: "bsc"}

	//Usdc           = (&web3_util.Asset{Name: "USDC", Symbol: "USDC", Address: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174", Decimals: 18, ChainName: "bsc"}).Init()
	//PancakeFactory = (&web3_util.Asset{Name: "PancakeFactory", Symbol: "PancakeFactory", Address: "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73", Decimals: 18, ChainName: "bsc"}).Init()
	//MultiCall      = (&web3_util.Asset{Name: "MultiCall", Symbol: "MultiCall", Address: "0x5dc53ed77bbc84f39c76fb4c84ac9f28384a4b55", Decimals: 18, ChainName: "bsc"}).Init()
	//PancakeRouter  = (&web3_util.Asset{Name: "PancakeRouter", Symbol: "PancakeRouter", Address: "0x10ED43C718714eb63d5aA57B78B54704E256024E", Decimals: 18, ChainName: "bsc"}).Init()
	//Wbnb           = (&web3_util.Asset{Name: "WBNB", Symbol: "WBNB", Address: "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c", Decimals: 18, ChainName: "bsc"}).Init()
	//Busd           = (&web3_util.Asset{Name: "BUSD", Symbol: "BUSD", Address: "0xe9e7cea3dedca5984780bafc599bd69add087d56", Decimals: 18, ChainName: "bsc"}).Init()
	//BscUsd         = (&web3_util.Asset{Name: "BUSD", Symbol: "BUSD", Address: "0x55d398326f99059fF775485246999027B3197955", Decimals: 18, ChainName: "bsc"}).Init()
	//WbnbBusdPair   = (&web3_util.Asset{Name: "WbnbBusdPair", Symbol: "WBPair", Address: "0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16", Decimals: 18, ChainName: "bsc"}).Init()
)
