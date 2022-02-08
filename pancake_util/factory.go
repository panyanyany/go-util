package pancake_util

import "go-util/web3_util/contract"

// @see https://docs.pancakeswap.finance/code/smart-contracts/pancakeswap-exchange/factory-v2#getpair

type FactoryRepo struct {
	Contract *contract.Contract
}
