package pancake_util

import (
	"fmt"
	"math/big"
	"sync"

	"tugou-hunter/util/go_util"
	"tugou-hunter/util/web3_util"
	"tugou-hunter/util/web3_util/bsc"
	"tugou-hunter/util/web3_util/cclient"
	"tugou-hunter/util/web3_util/contract"

	"github.com/cihub/seelog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/pancakeswap-sdk-go/constants"
	"github.com/panyanyany/pancakeswap-sdk-go/entities"
)

type Token struct {
	Addr                       string     `json:"addr"`
	Name                       string     `json:"name"`
	Symbol                     string     `json:"symbol"`
	Decimals                   int        `json:"decimals"`
	TotalSupply                *big.Int   `json:"totalSupply"`
	BalanceOfPair              *big.Int   `json:"balanceOfPair"`
	BusdPrice                  *big.Int   `json:"busdPrice"`
	BusdPriceHuman             *big.Float `json:"busdPriceHuman"`
	TotalBusdAmountOfPair      *big.Int   `json:"totalBusdAmountOfPair"`
	TotalBusdAmountOfPairHuman *big.Float `json:"totalBusdAmountOfPairHuman"`
}

type Pair struct {
	PairAddress          string     `json:"pairAddress"`
	Name                 string     `json:"name"`
	Symbol               string     `json:"symbol"`
	Symbol1              string     `json:"symbol1"`
	Symbol0              string     `json:"symbol0"`
	Decimals             int64      `json:"decimals"`
	TotalSupply          *big.Int   `json:"totalSupply"`
	Token0               *Token     `json:"token0"`
	Token1               *Token     `json:"token1"`
	TotalBusdAmount      *big.Int   `json:"totalBusdAmount"`
	TotalBusdAmountHuman *big.Float `json:"totalBusdAmountHuman"`
	FarmTvl              *big.Int   `json:"farmTVL"`
}

func (r *Token) ToHuman() {
	r.BusdPriceHuman = web3_util.FromWei(r.BusdPrice, 12)
	r.TotalBusdAmountOfPairHuman = web3_util.FromWei(r.TotalBusdAmountOfPair, 12)
}

func (r *Token) ToEntityToken() (token *entities.Token) {
	token, err := entities.NewToken(constants.Mainnet, common.HexToAddress(r.Addr), r.Decimals, r.Symbol, r.Name)
	go_util.Must(err)
	return
}

func (r *Pair) ToHuman() {
	r.TotalBusdAmountHuman = web3_util.FromWei(r.TotalBusdAmount, 12)
	r.Token0.ToHuman()
	r.Token1.ToHuman()
}

func (r *Pair) ToNamed(quoteAddr string) (np *NamedPair) {
	np = new(NamedPair)

	np.PairAddress = r.PairAddress
	np.Name = r.Name
	np.Symbol = r.Symbol
	np.Decimals = r.Decimals
	np.TotalSupply = r.TotalSupply
	np.TotalBusdAmount = r.TotalBusdAmount
	np.TotalBusdAmountHuman = r.TotalBusdAmountHuman
	np.FarmTvl = r.FarmTvl

	var tokenQuote, tokenBase *Token
	var symbolQuote, symbolBase string
	if r.Token0.Addr == quoteAddr {
		tokenQuote = r.Token0
		symbolQuote = r.Symbol0

		tokenBase = r.Token1
		symbolBase = r.Symbol1
	} else {
		tokenQuote = r.Token1
		symbolQuote = r.Symbol1

		tokenBase = r.Token0
		symbolBase = r.Symbol0
	}

	np.SymbolBase = symbolBase
	np.TokenBase = tokenBase
	np.SymbolQuote = symbolQuote
	np.TokenQuote = tokenQuote

	return
}

type NamedPair struct {
	PairAddress          string     `json:"pairAddress"`
	Name                 string     `json:"name"`
	Symbol               string     `json:"symbol"`
	SymbolBase           string     `json:"symbolBase"`
	SymbolQuote          string     `json:"symbolQuote"`
	Decimals             int64      `json:"decimals"`
	TotalSupply          *big.Int   `json:"totalSupply"`
	TokenBase            *Token     `json:"tokenBase"`
	TokenQuote           *Token     `json:"tokenQuote"`
	TotalBusdAmount      *big.Int   `json:"totalBusdAmount"`
	TotalBusdAmountHuman *big.Float `json:"totalBusdAmountHuman"`
	FarmTvl              *big.Int   `json:"farmTVL"`
}

type PancakeRepoKit struct {
	FactoryRepo   *FactoryRepo
	RouterRepo    *RouterRepo
	MultiCallRepo *MultiCallRepo
}

func InitAllContracts(client *cclient.CClient) (repokit *PancakeRepoKit, err error) {
	repokit = new(PancakeRepoKit)

	factoryRepo := &FactoryRepo{bsc.PancakeFactory}
	factoryRepo.Contract.Provider = client.Endpoints.EthClient
	err = factoryRepo.Contract.LoadAbi()
	if err != nil {
		err = fmt.Errorf("factoryRepo.LoadAbi(): %w", err)
		return
	}

	routerRepo := &RouterRepo{Contract: bsc.PancakeRouter}
	routerRepo.Contract.Provider = client.Endpoints.EthClient
	err = routerRepo.Contract.LoadAbi()
	if err != nil {
		err = fmt.Errorf("routerRepo.LoadAbi(): %w", err)
		return
	}

	multiCallRepo := &MultiCallRepo{bsc.MultiCall}
	multiCallRepo.Contract.Provider = client.Endpoints.EthClient
	err = multiCallRepo.Contract.LoadAbi()
	if err != nil {
		err = fmt.Errorf("multiCallRepo.LoadAbi(): %w", err)
		return
	}

	repokit.FactoryRepo = factoryRepo
	repokit.RouterRepo = routerRepo
	repokit.MultiCallRepo = multiCallRepo
	return
}

type GetPairOutput struct {
	Pair  *Pair
	Error error
}

func GetBnbOrBusdPair(PancakeContract, MultiCall *contract.Contract, quoteAddress common.Address) (*Pair, error) {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	outs := make([]*GetPairOutput, 0, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		pair, err := GetPair(PancakeContract, MultiCall, bsc.Wbnb.Address, quoteAddress)
		pairOut := &GetPairOutput{Pair: pair, Error: nil}
		if err != nil {
			err = fmt.Errorf("pancake_util.GetPair(BNB): %w", err)
			pairOut.Error = err
		}
		lock.Lock()
		outs = append(outs, pairOut)
		lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pair, err := GetPair(PancakeContract, MultiCall, bsc.Busd.Address, quoteAddress)
		pairOut := &GetPairOutput{Pair: pair, Error: nil}
		if err != nil {
			err = fmt.Errorf("pancake_util.GetPair(BUSD): %w", err)
			pairOut.Error = err
		}
		lock.Lock()
		outs = append(outs, pairOut)
		lock.Unlock()
	}()
	wg.Wait()

	var maxPair *GetPairOutput
	for _, pair := range outs {
		//if pair.Error != nil {
		//	maxPair = pair
		//	return maxPair.Pair, maxPair.Error
		//}
		if pair.Pair == nil {
			continue
		}
		if maxPair == nil {
			maxPair = pair
		} else if maxPair.Pair.TotalBusdAmount.Cmp(pair.Pair.TotalBusdAmount) == -1 {
			maxPair = pair
		}
	}

	if maxPair == nil {
		return nil, nil
	}

	return maxPair.Pair, maxPair.Error
}

func GetPair(PancakeContract, MultiCall *contract.Contract, baseAddress, quoteAddress common.Address) (pair *Pair, err error) {
	var result map[string]interface{}
	result, err = PancakeContract.Call("getPair", web3.Latest, baseAddress, quoteAddress)
	if err != nil {
		err = fmt.Errorf("Factory.Call: %w", err)
		//seelog.Error(err)
		return
	}
	pairAddress := result["0"].(web3.Address)
	if pairAddress.String() == web3.ZeroAddress.String() {
		// no pair
		seelog.Debugf("no pair")
		return
	}
	multiRepo := NewMultiCallRepo(MultiCall)
	pair, err = multiRepo.GetPairInfoWithPrice(pairAddress.String())
	if err != nil {
		err = fmt.Errorf("multiRepo.GetPairInfoWithPrice: %w", err)
		return
	}

	return
}
