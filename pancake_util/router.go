package pancake_util

import (
	"math/big"

	"tugou-hunter/util/web3_util/contract"

	"github.com/panyanyany/go-web3"
)

type RouterRepo struct {
	Contract *contract.Contract
}

func (r *RouterRepo) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Tx {
	methodName := "swapExactETHForTokensSupportingFeeOnTransferTokens"

	return contract.NewTx().
		SetMethod(methodName).
		AddArgs(amountOutMin, path, to, deadline).
		//SetGas(1800 * 1000).
		//SetGasPrice(unit.NewGWei(big.NewFloat(20)).Wei().Uint64()).
		SetContract(r.Contract)
}
