package web3_util

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/wallet"
)

func ToEth(balance float64, decimals int) *big.Int {
	if balance < 0 {
		panic("balance < 0")
	}

	op := big.NewInt(10)
	e := big.NewInt(int64(decimals))
	op.Exp(op, e, nil)

	opFloat := big.NewFloat(0).SetInt(op)

	bal := big.NewFloat(balance)
	bal = big.NewFloat(0).Mul(bal, opFloat)

	balInt, _ := bal.Int(nil)

	return balInt
}
func FromWei(balance *big.Int, decimals int) *big.Float {
	op := big.NewInt(10)
	e := big.NewInt(int64(decimals))
	op.Exp(op, e, nil)

	bal := new(big.Float).SetInt(balance)
	op2 := new(big.Float).SetInt(op)
	return bal.Quo(bal, op2)
}

func NewWalletFromPrivateKeyString(pk string) (key *wallet.Key, err error) {
	pk = pk[2:]
	bs, err := hex.DecodeString(pk)
	if err != nil {
		err = fmt.Errorf("decode: %w", err)
		return
	}

	key, err = wallet.NewWalletFromPrivKey(bs)
	if err != nil {
		err = fmt.Errorf("wallet.NewWalletFromPrivKey: %w", err)
		return
	}
	return
}

func CommonToWeb3(addr common.Address) *web3.Address {
	t := web3.HexToAddress(addr.String())
	return &t
}
