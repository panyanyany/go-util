package pancake_util

import (
	"fmt"

	"tugou-hunter/util/web3_util/bsc"
	"tugou-hunter/util/web3_util/contract"

	"github.com/cihub/seelog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/jsonrpc"
)

type Kit struct {
	RouterContract *contract.Contract
	Factory        *contract.Contract
	Router         *RouterRepo
	MultiCall      *contract.Contract
	Provider       jsonrpc.IEth
}

func NewKit(provider jsonrpc.IEth, routerContract *contract.Contract, factoryContract *contract.Contract) (r *Kit) {
	r = new(Kit)
	r.Provider = provider
	r.Factory = factoryContract
	r.RouterContract = routerContract
	return
}

func (r *Kit) MustInit() *Kit {
	err := r.Init()
	if err != nil {
		panic(err)
	}

	return r
}

func (r *Kit) Init() (err error) {
	provider := r.Provider
	//r.Factory = bsc.PancakeFactory
	r.Factory.Provider = provider
	err = r.Factory.LoadAbi()
	if err != nil {
		err = fmt.Errorf("r.Factory.LoadAbi(): %w", err)
		return
	}

	r.Router = &RouterRepo{Contract: r.RouterContract}
	r.Router.Contract.Provider = provider
	err = r.Router.Contract.LoadAbi()
	if err != nil {
		err = fmt.Errorf("r.Router.LoadAbi(): %w", err)
		return
	}

	r.MultiCall = bsc.MultiCall
	r.MultiCall.Provider = provider
	err = r.MultiCall.LoadAbi()
	if err != nil {
		err = fmt.Errorf("r.MultiCall.LoadAbi(): %w", err)
		return
	}

	return
}
func (r *Kit) GetPair(baseAddress, quoteAddress common.Address) (pair *Pair, err error) {
	var result map[string]interface{}
	result, err = r.Factory.Call("getPair", web3.Latest, baseAddress, quoteAddress)
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
	multiRepo := NewMultiCallRepo(r.MultiCall)
	pair, err = multiRepo.GetPairInfoWithPrice(pairAddress.String())
	if err != nil {
		err = fmt.Errorf("multiRepo.GetPairInfoWithPrice: %w", err)
		return
	}

	return
}
