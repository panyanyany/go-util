package token

import (
	"math/big"

	"tugou-hunter/util/go_util"

	"github.com/panyanyany/pancakeswap-sdk-go/entities"
)

// 不用 balance，因为在交易的时候，输入的只是数量，不是余额
type TokenAmount struct {
	Token *Token

	Amount      *big.Int
	AmountHuman *big.Float
}

func (r *TokenAmount) MustToEntityTokenAmount() *entities.TokenAmount {
	token := r.Token.MustToEntityToken()
	tokenAmount, err := entities.NewTokenAmount(token, r.Amount)
	go_util.Must(err)
	return tokenAmount
}
