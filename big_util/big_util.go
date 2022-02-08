package big_util

import "math/big"

func IntPercent(in *big.Int, a, b int64) (out *big.Int) {
	out = big.NewInt(0).Mul(in, big.NewInt(a))
	out = big.NewInt(0).Div(out, big.NewInt(b))
	return
}
