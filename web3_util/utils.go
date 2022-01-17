package web3_util

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/panyanyany/go-web3/wallet"
)

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
