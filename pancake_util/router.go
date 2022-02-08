package pancake_util

import (
	"math/big"

	"go-util/web3_util/contract"

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

func (r *RouterRepo) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Tx {
	methodName := "swapExactTokensForTokens"

	return contract.NewTx().
		SetMethod(methodName).
		AddArgs(amountIn, amountOutMin, path, to, deadline).
		//SetGas(1800 * 1000).
		//SetGasPrice(unit.NewGWei(big.NewFloat(20)).Wei().Uint64()).
		SetContract(r.Contract)
}

// 对于有税的合约，用 SwapExactTokensForTokens 貌似检查不出税率来
func (r *RouterRepo) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Tx {
	methodName := "swapExactTokensForTokensSupportingFeeOnTransferTokens"

	return contract.NewTx().
		SetMethod(methodName).
		AddArgs(amountIn, amountOutMin, path, to, deadline).
		//SetGas(1800 * 1000).
		//SetGasPrice(unit.NewGWei(big.NewFloat(20)).Wei().Uint64()).
		SetContract(r.Contract)
}
