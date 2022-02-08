package token

import (
	"tugou-hunter/util/go_util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/panyanyany/pancakeswap-sdk-go/constants"
	"github.com/panyanyany/pancakeswap-sdk-go/entities"
)

type Token struct {
	Address  common.Address
	Symbol   string
	Decimals int
	Name     string
}

func (r *Token) MustToEntityToken() *entities.Token {
	t, err := entities.NewToken(constants.Mainnet, r.Address, r.Decimals, r.Symbol, r.Name)
	go_util.Must(err)
	return t
}
